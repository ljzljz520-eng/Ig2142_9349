package domain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ParseMoveKey(key string, player Disc) (Move, error) {
	cleaned := strings.ToUpper(strings.TrimSpace(key))
	if len(cleaned) != 2 {
		return Move{}, errors.New("move key must contain a column and row")
	}
	column := int(cleaned[0] - 'A')
	row, err := strconv.Atoi(string(cleaned[1]))
	if err != nil {
		return Move{}, errors.New("move row must be numeric")
	}
	move := Move{Row: row - 1, Column: column, Player: player}
	if err := move.Validate(); err != nil {
		return Move{}, err
	}
	return move, nil
}

func FormatScore(score Score) string {
	return fmt.Sprintf("%d-%d", score.Black, score.White)
}

func FormatWinner(winner Disc) string {
	if winner == Black {
		return "black"
	}
	if winner == White {
		return "white"
	}
	return "draw"
}

func ValidateBoard(board Board) error {
	score := board.Counts()
	if score.Black+score.White > BoardSize*BoardSize {
		return errors.New("board has too many discs")
	}
	for row := 0; row < BoardSize; row++ {
		for column := 0; column < BoardSize; column++ {
			cell := board.Cells[row][column]
			if cell != Empty && cell != Black && cell != White {
				return errors.New("board contains an unknown disc")
			}
		}
	}
	return nil
}

func MatchSummaryLine(match MatchRecord) string {
	return match.ID + " " + match.BlackPlayer + " " + FormatScore(match.Score) + " " + FormatWinner(match.Winner)
}

func BoardDifference(before, after Board) int {
	changed := 0
	for row := 0; row < BoardSize; row++ {
		for column := 0; column < BoardSize; column++ {
			if before.Cells[row][column] != after.Cells[row][column] { changed++ }
		}
	}
	return changed
}

func IsKnownStatus(status MatchStatus) bool {
	return status == StatusActive || status == StatusCompleted || status == StatusAbandoned
}
