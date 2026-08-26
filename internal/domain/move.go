package domain

import "errors"

type Move struct {
	Row    int
	Column int
	Player Disc
}

func (m Move) Validate() error {
	if m.Player != Black && m.Player != White {
		return errors.New("move player must be black or white")
	}
	if m.Row < 0 || m.Row >= BoardSize || m.Column < 0 || m.Column >= BoardSize {
		return errors.New("move coordinate is invalid")
	}
	return nil
}

func (m Move) Key() string {
	return string(rune('A'+m.Column)) + string(rune('1'+m.Row))
}

type MoveResult struct {
	Move    Move
	Flipped int
	Score   Score
	Board   Board
}
