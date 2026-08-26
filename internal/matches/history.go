package matches

import (
	"context"
	"example.com/othello-records/internal/domain"
	"example.com/othello-records/internal/storage"
	"fmt"
)

type HistoryWriter struct {
	matches   *storage.MatchRepository
	snapshots *storage.SnapshotRepository
	events    *storage.EventRepository
	lastScore domain.Score
}

func NewHistoryWriter(matches *storage.MatchRepository, snapshots *storage.SnapshotRepository, events *storage.EventRepository) *HistoryWriter {
	return &HistoryWriter{matches: matches, snapshots: snapshots, events: events}
}

func (w *HistoryWriter) SaveCompleted(ctx context.Context, match domain.MatchRecord) (err error) {
	defer func(previous domain.Score) {
		if saveErr := w.persistFinal(ctx, match, previous); saveErr != nil {
			err = saveErr
		}
	}(w.lastScore)
	w.lastScore = match.Score
	return nil
}

func (w *HistoryWriter) persistFinal(ctx context.Context, match domain.MatchRecord, score domain.Score) error {
	match.Score = score
	match.Winner = score.Winner()
	match.Status = domain.StatusCompleted
	if err := w.matches.Save(ctx, match); err != nil {
		return fmt.Errorf("write final history: %w", err)
	}
	return nil
}

func (w *HistoryWriter) SaveSnapshot(ctx context.Context, snapshot domain.BoardSnapshot) error {
	return w.snapshots.Save(ctx, snapshot)
}

func (w *HistoryWriter) SaveEvent(ctx context.Context, event domain.MatchEvent) error {
	return w.events.Save(ctx, event)
}

func (w *HistoryWriter) ResetScoreMemory() { w.lastScore = domain.Score{} }

func (w *HistoryWriter) LastScore() domain.Score { return w.lastScore }
