package matches

import (
	"context"
	"example.com/othello-records/internal/domain"
)

type QueryService struct{ matches *MatchService }

func NewQueryService(matches *MatchService) *QueryService { return &QueryService{matches: matches} }

func (q *QueryService) History(ctx context.Context, filter domain.MatchFilter) ([]domain.MatchSummary, error) {
	records, err := q.matches.History(ctx, filter)
	if err != nil {
		return nil, err
	}
	summaries := make([]domain.MatchSummary, 0, len(records))
	for _, record := range records {
		summaries = append(summaries, domain.MatchSummary{MatchID: record.ID, BlackPlayer: record.BlackPlayer, WhitePlayer: record.WhitePlayer, Score: record.Score, Winner: record.Winner, Status: record.Status})
	}
	return summaries, nil
}

func (q *QueryService) Completed(ctx context.Context, playerID string) ([]domain.MatchSummary, error) {
	return q.History(ctx, domain.MatchFilter{PlayerID: playerID, Status: domain.StatusCompleted})
}
