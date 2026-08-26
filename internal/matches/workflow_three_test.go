package matches

import (
	"context"
	"example.com/othello-records/internal/domain"
	"example.com/othello-records/internal/fixtures"
	"testing"
)

func TestWorkflowHistoryFilterReturnsCompletedPlayerMatches(t *testing.T) {
	service, database := newTestService(t)
	defer database.Close()
	ctx := context.Background()
	match := fixtures.Match("workflow-three", 3)
	if err := service.Start(ctx, match); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Complete(ctx, match.ID, fixtures.FinalBoard(2, 9)); err != nil {
		t.Fatal(err)
	}
	items, err := service.History(ctx, domain.MatchFilter{PlayerID: "white-player", Status: domain.StatusCompleted})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].ID != match.ID {
		t.Fatalf("items = %#v", items)
	}
}
