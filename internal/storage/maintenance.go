package storage

import (
	"context"
	"fmt"
)

type EntityCounts struct {
	Players   int
	Matches   int
	Snapshots int
	Events    int
}

func (d *Database) EntityCounts(ctx context.Context) (EntityCounts, error) {
	queries := []struct {
		name   string
		target *int
	}{{"PlayerProfile", nil}, {"MatchRecord", nil}, {"BoardSnapshot", nil}, {"MatchEvent", nil}}
	counts := EntityCounts{}
	for index := range queries {
		var count int
		if err := d.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+queries[index].name).Scan(&count); err != nil {
			return EntityCounts{}, fmt.Errorf("count %s: %w", queries[index].name, err)
		}
		switch index {
		case 0:
			counts.Players = count
		case 1:
			counts.Matches = count
		case 2:
			counts.Snapshots = count
		case 3:
			counts.Events = count
		}
	}
	return counts, nil
}

func (d *Database) VerifySchema(ctx context.Context) error {
	var count int
	if err := d.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name IN ('PlayerProfile', 'MatchRecord', 'BoardSnapshot', 'MatchEvent')`).Scan(&count); err != nil {
		return err
	}
	if count != 4 {
		return fmt.Errorf("schema has %d domain tables", count)
	}
	return nil
}

func (d *Database) Optimize(ctx context.Context) error {
	if _, err := d.db.ExecContext(ctx, `PRAGMA optimize`); err != nil {
		return fmt.Errorf("optimize database: %w", err)
	}
	return nil
}

func (d *Database) BackupPath() string {
	if len(d.path) > 3 && d.path[len(d.path)-3:] == ".db" {
		return d.path[:len(d.path)-3] + ".backup.db"
	}
	return d.path + ".backup"
}
