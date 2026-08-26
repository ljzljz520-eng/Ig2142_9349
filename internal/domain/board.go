package domain

import (
	"errors"
	"fmt"
	"strings"
)

const BoardSize = 8

type Disc uint8

const (
	Empty Disc = iota
	Black
	White
)

func (d Disc) String() string {
	switch d {
	case Black:
		return "B"
	case White:
		return "W"
	default:
		return "."
	}
}

func ParseDisc(value string) (Disc, error) {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "B", "BLACK":
		return Black, nil
	case "W", "WHITE":
		return White, nil
	case ".", "EMPTY":
		return Empty, nil
	default:
		return Empty, fmt.Errorf("unknown disc %q", value)
	}
}

func (d Disc) Opponent() Disc {
	if d == Black {
		return White
	}
	if d == White {
		return Black
	}
	return Empty
}

type Board struct {
	Cells [BoardSize][BoardSize]Disc
}

func NewBoard() Board {
	var board Board
	board.Cells[3][3] = White
	board.Cells[4][4] = White
	board.Cells[3][4] = Black
	board.Cells[4][3] = Black
	return board
}

func (b Board) Clone() Board {
	var clone Board
	for row := 0; row < BoardSize; row++ {
		for column := 0; column < BoardSize; column++ {
			clone.Cells[row][column] = b.Cells[row][column]
		}
	}
	return clone
}

func (b Board) At(row, column int) (Disc, error) {
	if row < 0 || row >= BoardSize || column < 0 || column >= BoardSize {
		return Empty, errors.New("board coordinate out of range")
	}
	return b.Cells[row][column], nil
}

func (b *Board) Place(row, column int, disc Disc) error {
	if row < 0 || row >= BoardSize || column < 0 || column >= BoardSize {
		return errors.New("board coordinate out of range")
	}
	if disc != Black && disc != White {
		return errors.New("only black or white discs can be placed")
	}
	if b.Cells[row][column] != Empty {
		return errors.New("board coordinate is occupied")
	}
	b.Cells[row][column] = disc
	return nil
}

func (b Board) Counts() Score {
	var score Score
	for row := 0; row < BoardSize; row++ {
		for column := 0; column < BoardSize; column++ {
			switch b.Cells[row][column] {
			case Black:
				score.Black++
			case White:
				score.White++
			}
		}
	}
	return score
}

func (b Board) Encode() string {
	rows := make([]string, 0, BoardSize)
	for row := 0; row < BoardSize; row++ {
		var builder strings.Builder
		for column := 0; column < BoardSize; column++ {
			builder.WriteString(b.Cells[row][column].String())
		}
		rows = append(rows, builder.String())
	}
	return strings.Join(rows, "/")
}

func DecodeBoard(encoded string) (Board, error) {
	parts := strings.Split(encoded, "/")
	if len(parts) != BoardSize {
		return Board{}, errors.New("board must contain eight rows")
	}
	var board Board
	for row, part := range parts {
		if len(part) != BoardSize {
			return Board{}, errors.New("board rows must contain eight cells")
		}
		for column, symbol := range part {
			disc, err := ParseDisc(string(symbol))
			if err != nil {
				return Board{}, err
			}
			board.Cells[row][column] = disc
		}
	}
	return board, nil
}

func (b Board) IsFull() bool {
	for row := 0; row < BoardSize; row++ {
		for column := 0; column < BoardSize; column++ {
			if b.Cells[row][column] == Empty {
				return false
			}
		}
	}
	return true
}
