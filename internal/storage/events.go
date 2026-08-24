package storage

import (
	"context"
	"example.com/othello-records/internal/domain"
	"fmt"
)

type EventRepository struct{ database *Database }

func NewEventRepository(database *Database) *EventRepository {
	return &EventRepository{database: database}
}

func (r *EventRepository) Save(ctx context.Context, event domain.MatchEvent) error {
	if err := event.Validate(); err != nil {
		return err
	}
	_, err := r.database.db.ExecContext(ctx, `INSERT INTO MatchEvent (id, match_id, ply, player_id, row_index, column_index, action, created_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		event.ID, event.MatchID, event.Ply, event.PlayerID, event.Row, event.Column, event.Action, event.CreatedBy)
	if err != nil {
		return fmt.Errorf("save event: %w", err)
	}
	return nil
}

func (r *EventRepository) List(ctx context.Context, matchID string) ([]domain.MatchEvent, error) {
	rows, err := r.database.db.QueryContext(ctx, `SELECT id, match_id, ply, player_id, row_index, column_index, action, created_by FROM MatchEvent WHERE match_id = ? ORDER BY ply`, matchID)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}
	defer rows.Close()
	events := make([]domain.MatchEvent, 0)
	for rows.Next() {
		var event domain.MatchEvent
		if err := rows.Scan(&event.ID, &event.MatchID, &event.Ply, &event.PlayerID, &event.Row, &event.Column, &event.Action, &event.CreatedBy); err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate events: %w", err)
	}
	return events, nil
}
