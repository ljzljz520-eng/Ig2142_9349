package matches

import (
	"context"
	"errors"
	"example.com/othello-records/internal/domain"
	"example.com/othello-records/internal/rules"
)

type ReplayFrame struct {
	Ply       int
	Board     domain.Board
	Turn      domain.Disc
	Score     domain.Score
	EventID   string
	LegalNext int
}

func (s *MatchService) Replay(ctx context.Context, matchID string) ([]ReplayFrame, error) {
	match, err := s.matches.Get(ctx, matchID)
	if err != nil {
		return nil, err
	}
	snapshots, err := s.Snapshots(ctx, matchID)
	if err != nil {
		return nil, err
	}
	if len(snapshots) == 0 {
		return nil, errors.New("match has no snapshots")
	}
	frames := make([]ReplayFrame, 0, len(snapshots))
	for _, snapshot := range snapshots {
		frames = append(frames, ReplayFrame{Ply: snapshot.Ply, Board: snapshot.Board, Turn: snapshot.Turn, Score: snapshot.Board.Counts(), LegalNext: len(rules.LegalMoves(snapshot.Board, snapshot.Turn))})
	}
	if frames[len(frames)-1].Ply == 0 && match.Status == domain.StatusCompleted {
		frames[len(frames)-1].Score = match.Score
	}
	return frames, nil
}

func (s *MatchService) ValidateReplay(ctx context.Context, matchID string) error {
	frames, err := s.Replay(ctx, matchID)
	if err != nil {
		return err
	}
	for index := 1; index < len(frames); index++ {
		if frames[index].Ply <= frames[index-1].Ply {
			return errors.New("replay ply order is invalid")
		}
	}
	return nil
}
