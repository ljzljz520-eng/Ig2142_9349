package rules

import (
	"example.com/othello-records/internal/domain"
)

type PositionAnalysis struct {
	Score      domain.Score
	BlackMoves int
	WhiteMoves int
	BlackEdges int
	WhiteEdges int
	IsTerminal bool
	NextPlayer domain.Disc
}

func Analyze(board domain.Board, next domain.Disc) PositionAnalysis {
	score := board.Counts()
	analysis := PositionAnalysis{Score: score, BlackMoves: len(LegalMoves(board, domain.Black)), WhiteMoves: len(LegalMoves(board, domain.White)), NextPlayer: next}
	analysis.BlackEdges, analysis.WhiteEdges = edgeCounts(board)
	analysis.IsTerminal = IsTerminal(board, next)
	if analysis.IsTerminal {
		analysis.NextPlayer = domain.Empty
	}
	return analysis
}

func edgeCounts(board domain.Board) (int, int) {
	black, white := 0, 0
	for index := 0; index < domain.BoardSize; index++ {
		for _, point := range [][2]int{{0, index}, {domain.BoardSize - 1, index}, {index, 0}, {index, domain.BoardSize - 1}} {
			switch board.Cells[point[0]][point[1]] {
			case domain.Black:
				black++
			case domain.White:
				white++
			}
		}
	}
	return black, white
}

func ScoreDifference(score domain.Score, perspective domain.Disc) int {
	if perspective == domain.White {
		return score.White - score.Black
	}
	return score.Black - score.White
}

func LeadingPlayer(score domain.Score) domain.Disc {
	if score.Black == score.White {
		return domain.Empty
	}
	if score.Black > score.White {
		return domain.Black
	}
	return domain.White
}

func StableMoveCount(board domain.Board, player domain.Disc) int {
	count := 0
	for _, move := range LegalMoves(board, player) {
		result, err := ApplyMove(board, move)
		if err == nil && result.Score.Total() > board.Counts().Total() {
			count++
		}
	}
	return count
}
