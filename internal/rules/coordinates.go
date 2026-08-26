package rules

import (
	"errors"
	"example.com/othello-records/internal/domain"
)

func CoordinateFor(move domain.Move) (string, error) {
	if err := move.Validate(); err != nil {
		return "", err
	}
	return move.Key(), nil
}

func MovesByRow(board domain.Board, player domain.Disc) map[int][]domain.Move {
	grouped := make(map[int][]domain.Move)
	for _, move := range LegalMoves(board, player) {
		grouped[move.Row] = append(grouped[move.Row], move)
	}
	return grouped
}

func ValidateSequence(board domain.Board, moves []domain.Move) error {
	current := board.Clone()
	for _, move := range moves {
		if _, err := ApplyMove(current, move); err != nil {
			return err
		}
		result, _ := ApplyMove(current, move)
		current = result.Board
	}
	if len(moves) == 0 {
		return errors.New("move sequence cannot be empty")
	}
	return nil
}
