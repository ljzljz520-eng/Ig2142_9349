package domain

import "testing"

func TestBoardRoundTripPreservesCounts(t *testing.T) {
	board := NewBoard()
	encoded := board.Encode()
	decoded, err := DecodeBoard(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if !decoded.Counts().Equal(board.Counts()) {
		t.Fatalf("counts changed: %#v %#v", decoded.Counts(), board.Counts())
	}
}

func TestBoardRejectsInvalidCoordinates(t *testing.T) {
	board := NewBoard()
	if _, err := board.At(-1, 0); err == nil {
		t.Fatal("expected coordinate error")
	}
	if err := board.Place(0, 0, Empty); err == nil {
		t.Fatal("expected disc error")
	}
}
