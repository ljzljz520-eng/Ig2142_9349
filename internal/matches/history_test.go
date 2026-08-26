package matches

import (
	"context"
	"example.com/othello-records/internal/domain"
	"example.com/othello-records/internal/fixtures"
	"example.com/othello-records/internal/storage"
	"path/filepath"
	"testing"
)

func newTestService(t *testing.T) (*MatchService, *storage.Database) {
	t.Helper()
	database, err := storage.Open(filepath.Join(t.TempDir(), "matches.db"))
	if err != nil {
		t.Fatal(err)
	}
	if err := database.Initialize(context.Background()); err != nil {
		t.Fatal(err)
	}
	players := storage.NewPlayerRepository(database)
	if err := players.Save(context.Background(), domain.PlayerProfile{ID: "black-player", DisplayName: "Black Player", Rank: "beginner", Active: true}); err != nil {
		t.Fatal(err)
	}
	if err := players.Save(context.Background(), domain.PlayerProfile{ID: "white-player", DisplayName: "White Player", Rank: "beginner", Active: true}); err != nil {
		t.Fatal(err)
	}
	matches := storage.NewMatchRepository(database)
	history := NewHistoryWriter(matches, storage.NewSnapshotRepository(database), storage.NewEventRepository(database))
	return NewService(matches, players, history), database
}

func TestMatchHistoryStoresFinalScore(t *testing.T) {
	service, database := newTestService(t)
	defer database.Close()
	ctx := context.Background()
	first := fixtures.Match("game-one", 1)
	second := fixtures.Match("game-two", 2)
	if err := service.Start(ctx, first); err != nil {
		t.Fatal(err)
	}
	if err := service.Start(ctx, second); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Complete(ctx, first.ID, fixtures.FinalBoard(10, 4)); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Complete(ctx, second.ID, fixtures.FinalBoard(3, 11)); err != nil {
		t.Fatal(err)
	}
	record, err := service.Get(ctx, second.ID)
	if err != nil {
		t.Fatal(err)
	}
	if record.Score != (domain.Score{Black: 3, White: 11}) {
		t.Fatalf("score = %#v", record.Score)
	}
	if record.Winner != domain.White {
		t.Fatalf("winner = %v", record.Winner)
	}
}
