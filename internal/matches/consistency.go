package matches

import (
	"context"
	"errors"
	"example.com/othello-records/internal/domain"
	"example.com/othello-records/internal/rules"
)

type StateCheck struct {
	Valid       bool
	Status      domain.MatchStatus
	Score       domain.Score
	Expected    domain.Score
	TurnAllowed bool
	Terminal    bool
}

func (s *MatchService) CheckState(ctx context.Context, matchID string) (StateCheck, error) {
	match, err := s.matches.Get(ctx, matchID)
	if err != nil {
		return StateCheck{}, err
	}
	expected := match.Board.Counts()
	turnAllowed := match.Turn == domain.Black || match.Turn == domain.White
	terminal := rules.IsTerminal(match.Board, match.Turn)
	return StateCheck{Valid: expected == match.Score && (turnAllowed || match.Status == domain.StatusCompleted), Status: match.Status, Score: match.Score, Expected: expected, TurnAllowed: turnAllowed, Terminal: terminal}, nil
}

func (s *MatchService) CompleteAndRecord(ctx context.Context, matchID string, board domain.Board) (domain.MatchRecord, error) {
	match, err := s.Complete(ctx, matchID, board)
	if err != nil {
		return domain.MatchRecord{}, err
	}
	if err := (NewPlayerService(s.players)).RecordOutcome(ctx, match); err != nil {
		return domain.MatchRecord{}, err
	}
	return match, nil
}

func (s *MatchService) EnsurePlayable(ctx context.Context, matchID string) error {
	match, err := s.matches.Get(ctx, matchID)
	if err != nil {
		return err
	}
	if match.Status != domain.StatusActive {
		return errors.New("only active matches are playable")
	}
	if match.Turn == domain.Empty {
		return errors.New("active match has no turn")
	}
	if !rules.HasAnyMove(match.Board, match.Turn) && !rules.HasAnyMove(match.Board, match.Turn.Opponent()) {
		return errors.New("active match has no legal moves")
	}
	return nil
}

func (s *MatchService) CopyMatch(ctx context.Context, sourceID, targetID string, sequence int) error {
	if sourceID == targetID {
		return errors.New("source and target must differ")
	}
	source, err := s.matches.Get(ctx, sourceID)
	if err != nil {
		return err
	}
	source.ID, source.Sequence, source.Status = targetID, sequence, domain.StatusActive
	source.Winner = domain.Empty
	return s.Start(ctx, source)
}
