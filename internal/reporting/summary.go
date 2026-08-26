package reporting

import (
	"example.com/othello-records/internal/domain"
	"sort"
)

type MatchReport struct {
	MatchID      string
	Players      string
	ScoreLine    string
	WinnerLabel  string
	StatusLabel  string
	BoardDensity float64
}

func BuildMatchReport(match domain.MatchRecord) MatchReport {
	return MatchReport{MatchID: match.ID, Players: match.BlackPlayer + " vs " + match.WhitePlayer,
		ScoreLine: formatScore(match.Score), WinnerLabel: winnerLabel(match.Winner), StatusLabel: statusLabel(match.Status),
		BoardDensity: density(match.Score)}
}

func BuildReports(matches []domain.MatchRecord) []MatchReport {
	reports := make([]MatchReport, 0, len(matches))
	for _, match := range matches {
		reports = append(reports, BuildMatchReport(match))
	}
	return reports
}

func formatScore(score domain.Score) string {
	return formatInt(score.Black) + "-" + formatInt(score.White)
}

func winnerLabel(disc domain.Disc) string {
	if disc == domain.Black {
		return "black"
	}
	if disc == domain.White {
		return "white"
	}
	return "draw"
}

func statusLabel(status domain.MatchStatus) string {
	if status == domain.StatusCompleted {
		return "completed"
	}
	if status == domain.StatusAbandoned {
		return "abandoned"
	}
	return "active"
}

func density(score domain.Score) float64 {
	if score.Total() == 0 {
		return 0
	}
	return float64(score.Total()) / float64(domain.BoardSize*domain.BoardSize)
}

func formatInt(value int) string {
	if value == 0 {
		return "0"
	}
	if value < 0 {
		return "-" + formatInt(-value)
	}
	digits := make([]byte, 0, 4)
	for value > 0 {
		digits = append([]byte{byte('0' + value%10)}, digits...)
		value /= 10
	}
	return string(digits)
}

func SortReports(reports []MatchReport) {
	sort.SliceStable(reports, func(i, j int) bool { return reports[i].MatchID < reports[j].MatchID })
}
