package fixtures

import (
	"example.com/othello-records/internal/domain"
	"fmt"
)

func Board(rows ...string) domain.Board {
	if len(rows) != domain.BoardSize {
		panic("fixture board requires eight rows")
	}
	encoded := ""
	for index, row := range rows {
		if index > 0 {
			encoded += "/"
		}
		encoded += row
	}
	board, err := domain.DecodeBoard(encoded)
	if err != nil {
		panic(err)
	}
	return board
}

func Match(id string, sequence int) domain.MatchRecord {
	return domain.MatchRecord{ID: id, BlackPlayer: "black-player", WhitePlayer: "white-player", Board: domain.NewBoard(), Turn: domain.Black, Status: domain.StatusActive, Sequence: sequence, Label: fmt.Sprintf("fixture-%d", sequence)}
}

func FinalBoard(black, white int) domain.Board {
	if black < 0 || white < 0 || black+white > domain.BoardSize*domain.BoardSize {
		panic("fixture score exceeds board")
	}
	rows := make([]string, domain.BoardSize)
	remainingBlack, remainingWhite := black, white
	for row := range rows {
		bytes := make([]byte, domain.BoardSize)
		for column := range bytes {
			switch {
			case remainingBlack > 0:
				bytes[column] = 'B'
				remainingBlack--
			case remainingWhite > 0:
				bytes[column] = 'W'
				remainingWhite--
			default:
				bytes[column] = '.'
			}
		}
		rows[row] = string(bytes)
	}
	return Board(rows...)
}
