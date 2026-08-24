package storage

import (
	"context"
	"database/sql"
	"example.com/othello-records/internal/domain"
	"fmt"
)

type MatchTransaction struct {
	Match    domain.MatchRecord
	Snapshot domain.BoardSnapshot
	Event    domain.MatchEvent
}

func (d *Database) SaveMatchTransaction(ctx context.Context, transaction MatchTransaction) error {
	if err := transaction.Match.Validate(); err != nil {
		return err
	}
	if err := transaction.Snapshot.Validate(); err != nil {
		return err
	}
	if err := transaction.Event.Validate(); err != nil {
		return err
	}
	return d.WithTx(ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `INSERT INTO MatchRecord (id, black_player, white_player, board, black_score, white_score, turn, winner, status, sequence, label) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(id) DO UPDATE SET board=excluded.board, black_score=excluded.black_score, white_score=excluded.white_score, turn=excluded.turn, winner=excluded.winner, status=excluded.status`, transaction.Match.ID, transaction.Match.BlackPlayer, transaction.Match.WhitePlayer, transaction.Match.Board.Encode(), transaction.Match.Score.Black, transaction.Match.Score.White, transaction.Match.Turn, transaction.Match.Winner, transaction.Match.Status, transaction.Match.Sequence, transaction.Match.Label); err != nil {
			return fmt.Errorf("save transaction match: %w", err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO BoardSnapshot (id, match_id, ply, board, turn) VALUES (?, ?, ?, ?, ?)`, transaction.Snapshot.ID, transaction.Snapshot.MatchID, transaction.Snapshot.Ply, transaction.Snapshot.Board.Encode(), transaction.Snapshot.Turn); err != nil {
			return fmt.Errorf("save transaction snapshot: %w", err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO MatchEvent (id, match_id, ply, player_id, row_index, column_index, action, created_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, transaction.Event.ID, transaction.Event.MatchID, transaction.Event.Ply, transaction.Event.PlayerID, transaction.Event.Row, transaction.Event.Column, transaction.Event.Action, transaction.Event.CreatedBy); err != nil {
			return fmt.Errorf("save transaction event: %w", err)
		}
		return nil
	})
}

func (d *Database) TouchMatch(ctx context.Context, matchID string, label string) error {
	result, err := d.db.ExecContext(ctx, `UPDATE MatchRecord SET label = ? WHERE id = ?`, label, matchID)
	if err != nil {
		return fmt.Errorf("touch match: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("match %q was not found", matchID)
	}
	return nil
}

func (d *Database) CountCompletedByPlayer(ctx context.Context, playerID string) (int, error) {
	var count int
	if err := d.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM MatchRecord WHERE status = 'completed' AND (black_player = ? OR white_player = ?)`, playerID, playerID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count player matches: %w", err)
	}
	return count, nil
}
