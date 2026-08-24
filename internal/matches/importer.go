package matches

import (
	"context"
	"errors"
	"example.com/othello-records/internal/domain"
)

type MatchBatch struct{ Records []domain.MatchRecord }

func (s *MatchService) ImportBatch(ctx context.Context, batch MatchBatch) (int, error) {
	if len(batch.Records) == 0 {
		return 0, errors.New("batch cannot be empty")
	}
	count := 0
	for _, record := range batch.Records {
		if record.Status == domain.StatusCompleted {
			if err := s.history.SaveCompleted(ctx, record); err != nil {
				return count, err
			}
		} else if err := s.Start(ctx, record); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func (s *MatchService) Abandon(ctx context.Context, matchID string) error {
	match, err := s.matches.Get(ctx, matchID)
	if err != nil {
		return err
	}
	if match.Status != domain.StatusActive {
		return errors.New("only active matches can be abandoned")
	}
	match.Status = domain.StatusAbandoned
	return s.matches.Save(ctx, match)
}

func (s *MatchService) Rename(ctx context.Context, matchID, label string) error {
	if label == "" {
		return errors.New("match label is required")
	}
	match, err := s.matches.Get(ctx, matchID)
	if err != nil {
		return err
	}
	match.Label = label
	return s.matches.Save(ctx, match)
}
