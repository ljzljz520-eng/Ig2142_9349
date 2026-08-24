package rules

import (
	"example.com/othello-records/internal/domain"
	"testing"
)

func TestInitialBoardHasFourLegalMoves(t *testing.T) {
	moves := LegalMoves(domain.NewBoard(), domain.Black)
	if len(moves) != 4 {
		t.Fatalf("legal moves = %d", len(moves))
	}
	result, err := ApplyMove(domain.NewBoard(), moves[0])
	if err != nil {
		t.Fatal(err)
	}
	if result.Flipped != 1 || result.Score.Total() != 5 {
		t.Fatalf("result = %#v", result)
	}
}
