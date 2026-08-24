package app

import (
	"context"
	"example.com/othello-records/internal/domain"
)

func SeedPlayers(ctx context.Context, services *Services) error {
	players := []struct{ id, name, rank string }{{"black-player", "Black Player", "beginner"}, {"white-player", "White Player", "beginner"}}
	for _, item := range players {
		if _, err := services.Players.Register(ctx, item.id, item.name, item.rank); err != nil {
			return err
		}
	}
	return nil
}

func SeedMatch(ctx context.Context, services *Services, id string) error {
	match := domain.MatchRecord{ID: id, BlackPlayer: "black-player", WhitePlayer: "white-player", Status: domain.StatusActive, Sequence: 1, Label: "demo"}
	return services.Matches.Start(ctx, match)
}
