package rules

import (
	"errors"
	"example.com/othello-records/internal/domain"
)

func ApplyMove(board domain.Board, move domain.Move) (domain.MoveResult, error) {
	if err := move.Validate(); err != nil {
		return domain.MoveResult{}, err
	}
	flips := flipsForMove(board, move)
	if len(flips) == 0 {
		return domain.MoveResult{}, errors.New("move is not legal")
	}
	updated := board.Clone()
	updated.Cells[move.Row][move.Column] = move.Player
	for _, point := range flips {
		updated.Cells[point[0]][point[1]] = move.Player
	}
	return domain.MoveResult{Move: move, Flipped: len(flips), Score: updated.Counts(), Board: updated}, nil
}

func HasAnyMove(board domain.Board, player domain.Disc) bool {
	for row := 0; row < domain.BoardSize; row++ {
		for column := 0; column < domain.BoardSize; column++ {
			if IsLegalMove(board, domain.Move{Row: row, Column: column, Player: player}) {
				return true
			}
		}
	}
	return false
}

func IsTerminal(board domain.Board, player domain.Disc) bool {
	if board.IsFull() {
		return true
	}
	return !HasAnyMove(board, player) && !HasAnyMove(board, player.Opponent())
}

func NextTurn(board domain.Board, current domain.Disc) domain.Disc {
	other := current.Opponent()
	if HasAnyMove(board, other) {
		return other
	}
	if HasAnyMove(board, current) {
		return current
	}
	return domain.Empty
}
