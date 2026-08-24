package storage

import (
	"context"
	"database/sql"
	"example.com/othello-records/internal/domain"
	"fmt"
)

type MatchRepository struct{ database *Database }

func NewMatchRepository(database *Database) *MatchRepository {
	return &MatchRepository{database: database}
}

func (r *MatchRepository) Save(ctx context.Context, match domain.MatchRecord) error {
	if err := match.Validate(); err != nil {
		return err
	}
	_, err := r.database.db.ExecContext(ctx, `INSERT INTO MatchRecord
(id, black_player, white_player, board, black_score, white_score, turn, winner, status, sequence, label)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET black_player=excluded.black_player, white_player=excluded.white_player,
board=excluded.board, black_score=excluded.black_score, white_score=excluded.white_score, turn=excluded.turn,
winner=excluded.winner, status=excluded.status, sequence=excluded.sequence, label=excluded.label`,
		match.ID, match.BlackPlayer, match.WhitePlayer, match.Board.Encode(), match.Score.Black, match.Score.White,
		match.Turn, match.Winner, match.Status, match.Sequence, match.Label)
	if err != nil {
		return fmt.Errorf("save match: %w", err)
	}
	return nil
}

func (r *MatchRepository) Get(ctx context.Context, id string) (domain.MatchRecord, error) {
	row := r.database.db.QueryRowContext(ctx, `SELECT id, black_player, white_player, board, black_score, white_score, turn, winner, status, sequence, label FROM MatchRecord WHERE id = ?`, id)
	return scanMatch(row)
}

func (r *MatchRepository) List(ctx context.Context, filter domain.MatchFilter) ([]domain.MatchRecord, error) {
	filter = filter.Normalized()
	query := `SELECT id, black_player, white_player, board, black_score, white_score, turn, winner, status, sequence, label FROM MatchRecord WHERE 1=1`
	args := make([]any, 0, 4)
	if filter.PlayerID != "" {
		query += ` AND (black_player = ? OR white_player = ?)`
		args = append(args, filter.PlayerID, filter.PlayerID)
	}
	if filter.Status != "" {
		query += ` AND status = ?`
		args = append(args, filter.Status)
	}
	if filter.Winner != domain.Empty {
		query += ` AND winner = ?`
		args = append(args, filter.Winner)
	}
	query += ` ORDER BY sequence DESC, id DESC LIMIT ?`
	args = append(args, filter.Limit)
	rows, err := r.database.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list matches: %w", err)
	}
	defer rows.Close()
	matches := make([]domain.MatchRecord, 0)
	for rows.Next() {
		match, err := scanMatch(rows)
		if err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate matches: %w", err)
	}
	return matches, nil
}

func (r *MatchRepository) Delete(ctx context.Context, id string) error {
	return r.database.WithTx(ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `DELETE FROM MatchEvent WHERE match_id = ?`, id); err != nil {
			return fmt.Errorf("delete match events: %w", err)
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM BoardSnapshot WHERE match_id = ?`, id); err != nil {
			return fmt.Errorf("delete match snapshots: %w", err)
		}
		result, err := tx.ExecContext(ctx, `DELETE FROM MatchRecord WHERE id = ?`, id)
		if err != nil {
			return fmt.Errorf("delete match: %w", err)
		}
		count, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if count == 0 {
			return sql.ErrNoRows
		}
		return nil
	})
}

type scanner interface{ Scan(...any) error }

func scanMatch(row scanner) (domain.MatchRecord, error) {
	var match domain.MatchRecord
	var encoded string
	var turn, winner int
	if err := row.Scan(&match.ID, &match.BlackPlayer, &match.WhitePlayer, &encoded, &match.Score.Black, &match.Score.White, &turn, &winner, &match.Status, &match.Sequence, &match.Label); err != nil {
		return domain.MatchRecord{}, fmt.Errorf("scan match: %w", err)
	}
	board, err := domain.DecodeBoard(encoded)
	if err != nil {
		return domain.MatchRecord{}, fmt.Errorf("decode match board: %w", err)
	}
	match.Board, match.Turn, match.Winner = board, domain.Disc(turn), domain.Disc(winner)
	return match, nil
}
