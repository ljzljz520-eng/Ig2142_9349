package reporting

import (
	"encoding/json"
	"example.com/othello-records/internal/domain"
)

type HistoryExport struct {
	Matches []MatchReport      `json:"matches"`
	Players []LeaderboardEntry `json:"players"`
}

func ExportJSON(matches []domain.MatchRecord, players []domain.PlayerProfile) ([]byte, error) {
	report := HistoryExport{Matches: BuildReports(matches), Players: BuildLeaderboard(players)}
	return json.MarshalIndent(report, "", "  ")
}

func WinnerText(winner domain.Disc) string {
	switch winner {
	case domain.Black:
		return "Black wins"
	case domain.White:
		return "White wins"
	default:
		return "Draw"
	}
}
