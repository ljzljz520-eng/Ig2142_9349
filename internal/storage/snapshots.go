package storage

import (
	"context"
	"example.com/othello-records/internal/domain"
	"fmt"
)

type SnapshotRepository struct{ database *Database }

func NewSnapshotRepository(database *Database) *SnapshotRepository {
	return &SnapshotRepository{database: database}
}

func (r *SnapshotRepository) Save(ctx context.Context, snapshot domain.BoardSnapshot) error {
	if err := snapshot.Validate(); err != nil {
		return err
	}
	_, err := r.database.db.ExecContext(ctx, `INSERT INTO BoardSnapshot (id, match_id, ply, board, turn) VALUES (?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET board=excluded.board, turn=excluded.turn`, snapshot.ID, snapshot.MatchID, snapshot.Ply, snapshot.Board.Encode(), snapshot.Turn)
	if err != nil {
		return fmt.Errorf("save snapshot: %w", err)
	}
	return nil
}

func (r *SnapshotRepository) List(ctx context.Context, matchID string) ([]domain.BoardSnapshot, error) {
	rows, err := r.database.db.QueryContext(ctx, `SELECT id, match_id, ply, board, turn FROM BoardSnapshot WHERE match_id = ? ORDER BY ply`, matchID)
	if err != nil {
		return nil, fmt.Errorf("list snapshots: %w", err)
	}
	defer rows.Close()
	snapshots := make([]domain.BoardSnapshot, 0)
	for rows.Next() {
		var snapshot domain.BoardSnapshot
		var encoded string
		var turn int
		if err := rows.Scan(&snapshot.ID, &snapshot.MatchID, &snapshot.Ply, &encoded, &turn); err != nil {
			return nil, fmt.Errorf("scan snapshot: %w", err)
		}
		board, err := domain.DecodeBoard(encoded)
		if err != nil {
			return nil, fmt.Errorf("decode snapshot: %w", err)
		}
		snapshot.Board, snapshot.Turn = board, domain.Disc(turn)
		snapshots = append(snapshots, snapshot)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate snapshots: %w", err)
	}
	return snapshots, nil
}
