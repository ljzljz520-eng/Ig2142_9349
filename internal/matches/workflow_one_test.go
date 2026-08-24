package matches

import (
	"context"
	"example.com/othello-records/internal/domain"
	"example.com/othello-records/internal/fixtures"
	"testing"
)

func TestWorkflowRegisterStartAndComplete(t *testing.T) {
	service, database := newTestService(t)
	defer database.Close()
	ctx := context.Background()
	match := fixtures.Match("workflow-one", 1)
	if err := service.Start(ctx, match); err != nil {
		t.Fatal(err)
	}
	completed, err := service.Complete(ctx, match.ID, fixtures.FinalBoard(6, 2))
	if err != nil {
		t.Fatal(err)
	}
	if completed.Status != domain.StatusCompleted {
		t.Fatalf("status = %s", completed.Status)
	}
}
