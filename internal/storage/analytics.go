package storage

import (
	"context"
	"example.com/othello-records/internal/domain"
	"fmt"
)

type MatchTotals struct {
	Completed int
	Active    int
	Abandoned int
	BlackWins int
	WhiteWins int
	Draws     int
}

func (r *MatchRepository) Totals(ctx context.Context) (MatchTotals, error) {
	var totals MatchTotals
	row := r.database.db.QueryRowContext(ctx, `SELECT
	SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END),
	SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END),
	SUM(CASE WHEN status = 'abandoned' THEN 1 ELSE 0 END),
	SUM(CASE WHEN status = 'completed' AND winner = ? THEN 1 ELSE 0 END),
	SUM(CASE WHEN status = 'completed' AND winner = ? THEN 1 ELSE 0 END),
	SUM(CASE WHEN status = 'completed' AND winner = ? THEN 1 ELSE 0 END) FROM MatchRecord`, domain.Black, domain.White, domain.Empty)
	if err := row.Scan(&totals.Completed, &totals.Active, &totals.Abandoned, &totals.BlackWins, &totals.WhiteWins, &totals.Draws); err != nil {
		return MatchTotals{}, fmt.Errorf("read match totals: %w", err)
	}
	return totals, nil
}

func (r *MatchRepository) Latest(ctx context.Context, count int) ([]domain.MatchRecord, error) {
	if count <= 0 {
		count = 10
	}
	if count > 100 {
		count = 100
	}
	return r.List(ctx, domain.MatchFilter{Limit: count})
}
