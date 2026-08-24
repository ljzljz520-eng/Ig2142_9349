package matches

import (
	"context"
	"errors"
	"example.com/othello-records/internal/domain"
	"example.com/othello-records/internal/rules"
	"example.com/othello-records/internal/storage"
	"fmt"
)

type MatchService struct {
	matches *storage.MatchRepository
	players *storage.PlayerRepository
	history *HistoryWriter
}

func NewService(matches *storage.MatchRepository, players *storage.PlayerRepository, history *HistoryWriter) *MatchService {
	return &MatchService{matches: matches, players: players, history: history}
}

func (s *MatchService) Start(ctx context.Context, match domain.MatchRecord) error {
	if match.Status == "" {
		match.Status = domain.StatusActive
	}
	if match.Board == (domain.Board{}) {
		match.Board = domain.NewBoard()
	}
	if err := match.Validate(); err != nil {
		return err
	}
	match.Score = match.Board.Counts()
	if match.Turn == domain.Empty {
		match.Turn = domain.Black
	}
	if !rules.HasAnyMove(match.Board, match.Turn) {
		return errors.New("starting player has no legal move")
	}
	if err := s.matches.Save(ctx, match); err != nil {
		return err
	}
	snapshot := domain.BoardSnapshot{ID: match.ID + "-start", MatchID: match.ID, Ply: 0, Board: match.Board, Turn: match.Turn}
	return s.history.SaveSnapshot(ctx, snapshot)
}

func (s *MatchService) SubmitMove(ctx context.Context, matchID string, move domain.Move, eventID string) (domain.MoveResult, error) {
	match, err := s.matches.Get(ctx, matchID)
	if err != nil {
		return domain.MoveResult{}, err
	}
	if match.Status != domain.StatusActive {
		return domain.MoveResult{}, errors.New("match is not active")
	}
	if err := s.ensurePlayerTurn(match, move); err != nil {
		return domain.MoveResult{}, err
	}
	result, err := rules.ApplyMove(match.Board, move)
	if err != nil {
		return domain.MoveResult{}, err
	}
	match.Board = result.Board
	match.Score = result.Score
	match.Turn = rules.NextTurn(match.Board, move.Player)
	if match.Turn == domain.Empty || rules.IsTerminal(match.Board, match.Turn) {
		match.Status = domain.StatusCompleted
		match.Winner = match.Score.Winner()
	}
	if err := s.matches.Save(ctx, match); err != nil {
		return domain.MoveResult{}, err
	}
	event := domain.MatchEvent{ID: eventID, MatchID: matchID, Ply: match.Sequence + 1, PlayerID: playerForDisc(match, move.Player), Row: move.Row, Column: move.Column, Action: "place", CreatedBy: "match-service"}
	if err := s.history.SaveEvent(ctx, event); err != nil {
		return domain.MoveResult{}, err
	}
	snapshot := domain.BoardSnapshot{ID: fmt.Sprintf("%s-%d", matchID, event.Ply), MatchID: matchID, Ply: event.Ply, Board: match.Board, Turn: match.Turn}
	if err := s.history.SaveSnapshot(ctx, snapshot); err != nil {
		return domain.MoveResult{}, err
	}
	return result, nil
}

func (s *MatchService) Complete(ctx context.Context, matchID string, finalBoard domain.Board) (domain.MatchRecord, error) {
	match, err := s.matches.Get(ctx, matchID)
	if err != nil {
		return domain.MatchRecord{}, err
	}
	if match.Status != domain.StatusActive {
		return domain.MatchRecord{}, errors.New("match is not active")
	}
	match = match.WithBoard(finalBoard).Finish()
	if err := s.history.SaveCompleted(ctx, match); err != nil {
		return domain.MatchRecord{}, err
	}
	return match, nil
}

func (s *MatchService) ensurePlayerTurn(match domain.MatchRecord, move domain.Move) error {
	if move.Player != match.Turn {
		return errors.New("move does not match current turn")
	}
	if move.Player != domain.Black && move.Player != domain.White {
		return errors.New("unsupported player color")
	}
	return nil
}

func playerForDisc(match domain.MatchRecord, disc domain.Disc) string {
	if disc == domain.Black {
		return match.BlackPlayer
	}
	return match.WhitePlayer
}

func (s *MatchService) Get(ctx context.Context, id string) (domain.MatchRecord, error) {
	return s.matches.Get(ctx, id)
}

func (s *MatchService) History(ctx context.Context, filter domain.MatchFilter) ([]domain.MatchRecord, error) {
	return s.matches.List(ctx, filter)
}

func (s *MatchService) Events(ctx context.Context, id string) ([]domain.MatchEvent, error) {
	return s.history.events.List(ctx, id)
}

func (s *MatchService) Snapshots(ctx context.Context, id string) ([]domain.BoardSnapshot, error) {
	return s.history.snapshots.List(ctx, id)
}
