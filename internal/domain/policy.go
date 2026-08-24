package domain

import "errors"

type MatchPolicy struct {
	RequireDistinctPlayers bool
	RequireActivePlayers   bool
	AllowAbandonment       bool
	MaximumSequence        int
}

func DefaultMatchPolicy() MatchPolicy {
	return MatchPolicy{RequireDistinctPlayers: true, RequireActivePlayers: true, AllowAbandonment: true, MaximumSequence: 1000000}
}

func (p MatchPolicy) ValidateMatch(match MatchRecord, black, white PlayerProfile) error {
	if p.RequireDistinctPlayers {
		if err := EnsureDistinctPlayers(match.BlackPlayer, match.WhitePlayer); err != nil {
			return err
		}
	}
	if p.RequireActivePlayers && (!black.Active || !white.Active) {
		return errors.New("both players must be active")
	}
	if p.MaximumSequence > 0 && match.Sequence > p.MaximumSequence {
		return errors.New("match sequence exceeds policy")
	}
	return nil
}

func (p MatchPolicy) CanAbandon(status MatchStatus) bool {
	return p.AllowAbandonment && status == StatusActive
}

func (p MatchPolicy) CanEdit(status MatchStatus) bool {
	return status == StatusActive
}

func (p MatchPolicy) Outcome(score Score) string {
	if score.Black == score.White {
		return "draw"
	}
	if score.Black > score.White {
		return "black-win"
	}
	return "white-win"
}

func (p MatchPolicy) ValidateScore(score Score) error {
	if score.Black < 0 || score.White < 0 {
		return errors.New("score cannot be negative")
	}
	if score.Total() > BoardSize*BoardSize {
		return errors.New("score exceeds board capacity")
	}
	return nil
}
