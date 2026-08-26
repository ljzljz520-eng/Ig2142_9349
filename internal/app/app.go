package app

import (
	"context"
	"example.com/othello-records/internal/httpapi"
	"example.com/othello-records/internal/matches"
	"example.com/othello-records/internal/storage"
	"fmt"
)

type Services struct {
	Database *storage.Database
	Players  *matches.PlayerService
	Matches  *matches.MatchService
	Queries  *matches.QueryService
	HTTP     *httpapi.Server
}

func OpenServices(ctx context.Context, path string) (*Services, error) {
	database, err := storage.Open(path)
	if err != nil {
		return nil, err
	}
	if err := database.Initialize(ctx); err != nil {
		_ = database.Close()
		return nil, err
	}
	playerRepository := storage.NewPlayerRepository(database)
	matchRepository := storage.NewMatchRepository(database)
	eventRepository := storage.NewEventRepository(database)
	snapshotRepository := storage.NewSnapshotRepository(database)
	history := matches.NewHistoryWriter(matchRepository, snapshotRepository, eventRepository)
	matchService := matches.NewService(matchRepository, playerRepository, history)
	playerService := matches.NewPlayerService(playerRepository)
	queryService := matches.NewQueryService(matchService)
	server := httpapi.NewServer(matchService, playerService, queryService)
	return &Services{Database: database, Players: playerService, Matches: matchService, Queries: queryService, HTTP: server}, nil
}

func (s *Services) Close() error {
	if s == nil || s.Database == nil {
		return nil
	}
	return s.Database.Close()
}

func (s *Services) Run(ctx context.Context, address string) error {
	if s == nil || s.HTTP == nil {
		return fmt.Errorf("services are not initialized")
	}
	return s.HTTP.Serve(ctx, address)
}
