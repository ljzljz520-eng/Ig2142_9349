package storage

import (
	"context"
	"fmt"
)

func (d *Database) Initialize(ctx context.Context) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS PlayerProfile (
			id TEXT PRIMARY KEY,
			display_name TEXT NOT NULL,
			rank TEXT NOT NULL,
			active INTEGER NOT NULL,
			games INTEGER NOT NULL,
			wins INTEGER NOT NULL,
			losses INTEGER NOT NULL,
			draws INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS MatchRecord (
			id TEXT PRIMARY KEY,
			black_player TEXT NOT NULL,
			white_player TEXT NOT NULL,
			board TEXT NOT NULL,
			black_score INTEGER NOT NULL,
			white_score INTEGER NOT NULL,
			turn INTEGER NOT NULL,
			winner INTEGER NOT NULL,
			status TEXT NOT NULL,
			sequence INTEGER NOT NULL,
			label TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS BoardSnapshot (
			id TEXT PRIMARY KEY,
			match_id TEXT NOT NULL,
			ply INTEGER NOT NULL,
			board TEXT NOT NULL,
			turn INTEGER NOT NULL,
			FOREIGN KEY(match_id) REFERENCES MatchRecord(id)
		)`,
		`CREATE TABLE IF NOT EXISTS MatchEvent (
			id TEXT PRIMARY KEY,
			match_id TEXT NOT NULL,
			ply INTEGER NOT NULL,
			player_id TEXT NOT NULL,
			row_index INTEGER NOT NULL,
			column_index INTEGER NOT NULL,
			action TEXT NOT NULL,
			created_by TEXT NOT NULL,
			FOREIGN KEY(match_id) REFERENCES MatchRecord(id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_matchrecord_status ON MatchRecord(status)`,
		`CREATE INDEX IF NOT EXISTS idx_matchevent_match ON MatchEvent(match_id, ply)`,
	}
	for _, statement := range statements {
		if _, err := d.db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("initialize schema: %w", err)
		}
	}
	return nil
}
