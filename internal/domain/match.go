package domain

import (
	"errors"
	"fmt"
)

type Score struct {
	Black int
	White int
}

func (s Score) Total() int { return s.Black + s.White }

func (s Score) Winner() Disc {
	if s.Black > s.White {
		return Black
	}
	if s.White > s.Black {
		return White
	}
	return Empty
}

func (s Score) Equal(other Score) bool {
	return s.Black == other.Black && s.White == other.White
}

type MatchStatus string

const (
	StatusActive    MatchStatus = "active"
	StatusCompleted MatchStatus = "completed"
	StatusAbandoned MatchStatus = "abandoned"
)

type MatchRecord struct {
	ID          string
	BlackPlayer string
	WhitePlayer string
	Board       Board
	Score       Score
	Turn        Disc
	Winner      Disc
	Status      MatchStatus
	Sequence    int
	Label       string
}

func (m MatchRecord) Validate() error {
	if m.ID == "" {
		return errors.New("match id is required")
	}
	if m.BlackPlayer == "" || m.WhitePlayer == "" {
		return errors.New("both players are required")
	}
	if m.BlackPlayer == m.WhitePlayer {
		return errors.New("players must be distinct")
	}
	if m.Status != StatusActive && m.Status != StatusCompleted && m.Status != StatusAbandoned {
		return fmt.Errorf("invalid match status %q", m.Status)
	}
	if m.Score.Black < 0 || m.Score.White < 0 || m.Score.Total() > BoardSize*BoardSize {
		return errors.New("score is outside board bounds")
	}
	return nil
}

func (m MatchRecord) IsComplete() bool { return m.Status == StatusCompleted }

func (m MatchRecord) WithBoard(board Board) MatchRecord {
	m.Board = board.Clone()
	m.Score = board.Counts()
	return m
}

func (m MatchRecord) Finish() MatchRecord {
	m.Status = StatusCompleted
	m.Score = m.Board.Counts()
	m.Winner = m.Score.Winner()
	return m
}

type PlayerProfile struct {
	ID          string
	DisplayName string
	Rank        string
	Active      bool
	Games       int
	Wins        int
	Losses      int
	Draws       int
}

func (p PlayerProfile) Validate() error {
	if p.ID == "" {
		return errors.New("player id is required")
	}
	if p.DisplayName == "" {
		return errors.New("player display name is required")
	}
	if p.Games < 0 || p.Wins < 0 || p.Losses < 0 || p.Draws < 0 {
		return errors.New("player statistics cannot be negative")
	}
	if p.Wins+p.Losses+p.Draws != p.Games {
		return errors.New("player statistics do not balance")
	}
	return nil
}

type BoardSnapshot struct {
	ID      string
	MatchID string
	Ply     int
	Board   Board
	Turn    Disc
}

func (s BoardSnapshot) Validate() error {
	if s.ID == "" || s.MatchID == "" {
		return errors.New("snapshot identifiers are required")
	}
	if s.Ply < 0 {
		return errors.New("snapshot ply cannot be negative")
	}
	if s.Turn != Black && s.Turn != White {
		return errors.New("snapshot turn must be black or white")
	}
	return nil
}

type MatchEvent struct {
	ID        string
	MatchID   string
	Ply       int
	PlayerID  string
	Row       int
	Column    int
	Action    string
	CreatedBy string
}

func (e MatchEvent) Validate() error {
	if e.ID == "" || e.MatchID == "" || e.PlayerID == "" {
		return errors.New("event identifiers are required")
	}
	if e.Ply < 1 || e.Row < 0 || e.Row >= BoardSize || e.Column < 0 || e.Column >= BoardSize {
		return errors.New("event coordinates are invalid")
	}
	if e.Action == "" {
		return errors.New("event action is required")
	}
	return nil
}

type MatchFilter struct {
	PlayerID string
	Status   MatchStatus
	Winner   Disc
	Limit    int
}

func (f MatchFilter) Normalized() MatchFilter {
	if f.Limit <= 0 {
		f.Limit = 50
	}
	if f.Limit > 500 {
		f.Limit = 500
	}
	return f
}

type MatchSummary struct {
	MatchID     string
	BlackPlayer string
	WhitePlayer string
	Score       Score
	Winner      Disc
	Status      MatchStatus
}
