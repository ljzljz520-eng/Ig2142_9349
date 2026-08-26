package domain

import "testing"

func TestMatchValidationAndFinish(t *testing.T) {
	match := MatchRecord{ID: "m1", BlackPlayer: "b", WhitePlayer: "w", Status: StatusActive, Board: NewBoard()}
	if err := match.Validate(); err != nil {
		t.Fatal(err)
	}
	finished := match.Finish()
	if finished.Status != StatusCompleted {
		t.Fatalf("status = %s", finished.Status)
	}
	if finished.Score != (Score{Black: 2, White: 2}) {
		t.Fatalf("score = %#v", finished.Score)
	}
}
