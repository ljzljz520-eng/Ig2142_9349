package storage

import (
	"context"
	"example.com/othello-records/internal/domain"
	"path/filepath"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.db")
	ctx := context.Background()
	database, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := database.Initialize(ctx); err != nil {
		t.Fatal(err)
	}
	repository := NewMatchRepository(database)
	match := domain.MatchRecord{ID: "persisted", BlackPlayer: "b", WhitePlayer: "w", Board: domain.NewBoard(), Score: domain.NewBoard().Counts(), Turn: domain.Black, Winner: domain.Empty, Status: domain.StatusCompleted, Sequence: 1, Label: "reopen"}
	if err := repository.Save(ctx, match); err != nil {
		t.Fatal(err)
	}
	if err := database.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	if err := reopened.Initialize(ctx); err != nil {
		t.Fatal(err)
	}
	loaded, err := NewMatchRepository(reopened).Get(ctx, "persisted")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.ID != match.ID || loaded.Score != match.Score {
		t.Fatalf("loaded = %#v", loaded)
	}
}
