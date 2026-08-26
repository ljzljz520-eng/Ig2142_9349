package matches

import (
	"context"
	"errors"
	"example.com/othello-records/internal/domain"
	"example.com/othello-records/internal/storage"
)

type PlayerService struct{ repository *storage.PlayerRepository }

func NewPlayerService(repository *storage.PlayerRepository) *PlayerService {
	return &PlayerService{repository: repository}
}

func (s *PlayerService) Register(ctx context.Context, id, displayName, rank string) (domain.PlayerProfile, error) {
	normalID, err := domain.NormalizeID(id)
	if err != nil {
		return domain.PlayerProfile{}, err
	}
	name, err := domain.NormalizeDisplayName(displayName)
	if err != nil {
		return domain.PlayerProfile{}, err
	}
	if rank == "" {
		rank = "beginner"
	}
	player := domain.PlayerProfile{ID: normalID, DisplayName: name, Rank: rank, Active: true}
	if err := s.repository.Save(ctx, player); err != nil {
		return domain.PlayerProfile{}, err
	}
	return player, nil
}

func (s *PlayerService) Get(ctx context.Context, id string) (domain.PlayerProfile, error) {
	if id == "" {
		return domain.PlayerProfile{}, errors.New("player id is required")
	}
	return s.repository.Get(ctx, id)
}

func (s *PlayerService) List(ctx context.Context, activeOnly bool) ([]domain.PlayerProfile, error) {
	return s.repository.List(ctx, activeOnly)
}

func (s *PlayerService) Deactivate(ctx context.Context, id string) error {
	player, err := s.repository.Get(ctx, id)
	if err != nil {
		return err
	}
	player.Active = false
	return s.repository.Save(ctx, player)
}

func (s *PlayerService) RecordOutcome(ctx context.Context, match domain.MatchRecord) error {
	black, err := s.repository.Get(ctx, match.BlackPlayer)
	if err != nil {
		return err
	}
	white, err := s.repository.Get(ctx, match.WhitePlayer)
	if err != nil {
		return err
	}
	black.Games++
	white.Games++
	switch match.Winner {
	case domain.Black:
		black.Wins++
		white.Losses++
	case domain.White:
		white.Wins++
		black.Losses++
	default:
		black.Draws++
		white.Draws++
	}
	if err := s.repository.Save(ctx, black); err != nil {
		return err
	}
	return s.repository.Save(ctx, white)
}
