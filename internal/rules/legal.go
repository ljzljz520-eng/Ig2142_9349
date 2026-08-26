package rules

import "example.com/othello-records/internal/domain"

func LegalMoves(board domain.Board, player domain.Disc) []domain.Move {
	moves := make([]domain.Move, 0)
	if player != domain.Black && player != domain.White {
		return moves
	}
	for row := 0; row < domain.BoardSize; row++ {
		for column := 0; column < domain.BoardSize; column++ {
			move := domain.Move{Row: row, Column: column, Player: player}
			if len(flipsForMove(board, move)) > 0 {
				moves = append(moves, move)
			}
		}
	}
	return moves
}

func IsLegalMove(board domain.Board, move domain.Move) bool {
	if move.Validate() != nil {
		return false
	}
	return len(flipsForMove(board, move)) > 0
}

func flipsForMove(board domain.Board, move domain.Move) [][2]int {
	if !inside(move.Row, move.Column) || board.Cells[move.Row][move.Column] != domain.Empty {
		return nil
	}
	flips := make([][2]int, 0)
	for _, direction := range directions {
		row := move.Row + direction.rowDelta
		column := move.Column + direction.columnDelta
		line := make([][2]int, 0)
		for inside(row, column) && board.Cells[row][column] == move.Player.Opponent() {
			line = append(line, [2]int{row, column})
			row += direction.rowDelta
			column += direction.columnDelta
		}
		if len(line) > 0 && inside(row, column) && board.Cells[row][column] == move.Player {
			flips = append(flips, line...)
		}
	}
	return flips
}
