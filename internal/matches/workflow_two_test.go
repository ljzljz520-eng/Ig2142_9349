package matches

import (
	"context"
	"example.com/othello-records/internal/domain"
	"example.com/othello-records/internal/fixtures"
	"testing"
)

func TestWorkflowLegalMovePersistsEventAndSnapshot(t *testing.T) {
	service, database := newTestService(t)
	defer database.Close()
	ctx := context.Background()
	match := fixtures.Match("workflow-two", 2)
	if err := service.Start(ctx, match); err != nil {
		t.Fatal(err)
	}
	result, err := service.SubmitMove(ctx, match.ID, domain.Move{Row: 2, Column: 3, Player: domain.Black}, "event-one")
	if err != nil {
		t.Fatal(err)
	}
	if result.Flipped != 1 {
		t.Fatalf("flipped = %d", result.Flipped)
	}
	events, err := service.Events(ctx, match.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("events = %d", len(events))
	}
	snapshots, err := service.Snapshots(ctx, match.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshots) != 2 {
		t.Fatalf("snapshots = %d", len(snapshots))
	}
}
