package app

import (
	"context"
	"errors"
	"example.com/othello-records/internal/domain"
	"example.com/othello-records/internal/reporting"
)

type HealthSnapshot struct {
	DatabasePath string
	Players      int
	Matches      int
	Snapshots    int
	Events       int
	SchemaReady  bool
}

func (s *Services) Health(ctx context.Context) (HealthSnapshot, error) {
	if s == nil || s.Database == nil {
		return HealthSnapshot{}, errors.New("services are not initialized")
	}
	counts, err := s.Database.EntityCounts(ctx)
	if err != nil {
		return HealthSnapshot{}, err
	}
	if err := s.Database.VerifySchema(ctx); err != nil {
		return HealthSnapshot{}, err
	}
	return HealthSnapshot{DatabasePath: s.Database.Path(), Players: counts.Players, Matches: counts.Matches, Snapshots: counts.Snapshots, Events: counts.Events, SchemaReady: true}, nil
}

func (s *Services) HistoryExport(ctx context.Context, filter domain.MatchFilter) ([]byte, error) {
	if s == nil || s.Database == nil {
		return nil, errors.New("services are not initialized")
	}
	matches, err := s.Matches.History(ctx, filter)
	if err != nil {
		return nil, err
	}
	players, err := s.Players.List(ctx, false)
	if err != nil {
		return nil, err
	}
	return reporting.ExportJSON(matches, players)
}

func (s *Services) SeedDemo(ctx context.Context) error {
	if err := SeedPlayers(ctx, s); err != nil {
		return err
	}
	if err := SeedMatch(ctx, s, "demo-match"); err != nil {
		return err
	}
	return nil
}

func (s *Services) CloseAndOptimize(ctx context.Context) error {
	if s == nil || s.Database == nil {
		return nil
	}
	if err := s.Database.Optimize(ctx); err != nil {
		return err
	}
	return s.Close()
}

func (s *Services) Validate(ctx context.Context) error {
	if s == nil || s.Database == nil || s.Matches == nil || s.Players == nil {
		return errors.New("services are incomplete")
	}
	return s.Database.VerifySchema(ctx)
}
