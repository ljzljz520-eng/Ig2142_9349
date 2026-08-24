package storage

import (
	"context"
	"example.com/othello-records/internal/domain"
	"testing"
)

func TestPlayerRepositoryUpdatesStatistics(t *testing.T) {
	database, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	if err := database.Initialize(context.Background()); err != nil {
		t.Fatal(err)
	}
	repository := NewPlayerRepository(database)
	player := domain.PlayerProfile{ID: "p", DisplayName: "Player One", Rank: "novice", Active: true}
	if err := repository.Save(context.Background(), player); err != nil {
		t.Fatal(err)
	}
	player.Games, player.Wins = 1, 1
	if err := repository.Save(context.Background(), player); err != nil {
		t.Fatal(err)
	}
	loaded, err := repository.Get(context.Background(), "p")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Wins != 1 || loaded.Games != 1 {
		t.Fatalf("loaded = %#v", loaded)
	}
}
