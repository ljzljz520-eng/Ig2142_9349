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
}

func NewHistoryWriter(matches *storage.MatchRepository, snapshots *storage.SnapshotRepository, events *storage.EventRepository) *HistoryWriter {
	return &HistoryWriter{matches: matches, snapshots: snapshots, events: events}
}

// SaveCompleted persists the final win/loss record for a match. The score and
// winner are derived from the match's final board so the stored history always
// reflects this match's own disc counts and winner, never a previous match's.
func (w *HistoryWriter) SaveCompleted(ctx context.Context, match domain.MatchRecord) error {
	final := match
	final.Status = domain.StatusCompleted
	final.Score = final.Board.Counts()
	final.Winner = final.Score.Winner()
	if err := w.matches.Save(ctx, final); err != nil {
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
