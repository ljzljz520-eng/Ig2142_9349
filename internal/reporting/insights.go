package reporting

import (
	"example.com/othello-records/internal/domain"
	"example.com/othello-records/internal/rules"
)

type Insight struct {
	MatchID       string
	Leader        string
	Margin        int
	BlackMobility int
	WhiteMobility int
	Terminal      bool
}

func BuildInsight(match domain.MatchRecord) Insight {
	analysis := rules.Analyze(match.Board, match.Turn)
	leader := domain.FormatWinner(rules.LeadingPlayer(match.Score))
	margin := rules.ScoreDifference(match.Score, domain.Black)
	if margin < 0 {
		margin = -margin
	}
	return Insight{MatchID: match.ID, Leader: leader, Margin: margin, BlackMobility: analysis.BlackMoves, WhiteMobility: analysis.WhiteMoves, Terminal: analysis.IsTerminal}
}

func BuildInsights(matches []domain.MatchRecord) []Insight {
	insights := make([]Insight, 0, len(matches))
	for _, match := range matches {
		insights = append(insights, BuildInsight(match))
	}
	return insights
}
