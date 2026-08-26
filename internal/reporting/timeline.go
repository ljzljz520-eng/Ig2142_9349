package reporting

import (
	"example.com/othello-records/internal/domain"
	"strings"
)

type TimelineRow struct {
	Ply        int
	Action     string
	PlayerID   string
	Coordinate string
	Score      domain.Score
}

func BuildTimeline(events []domain.MatchEvent, snapshots []domain.BoardSnapshot) []TimelineRow {
	rows := make([]TimelineRow, 0, len(events))
	for _, event := range events {
		score := domain.Score{}
		for _, snapshot := range snapshots {
			if snapshot.Ply == event.Ply {
				score = snapshot.Board.Counts()
				break
			}
		}
		move := domain.Move{Row: event.Row, Column: event.Column}
		rows = append(rows, TimelineRow{Ply: event.Ply, Action: event.Action, PlayerID: event.PlayerID, Coordinate: move.Key(), Score: score})
	}
	return rows
}

func TimelineText(rows []TimelineRow) string {
	parts := make([]string, 0, len(rows))
	for _, row := range rows {
		parts = append(parts, row.Coordinate+"="+domain.FormatScore(row.Score))
	}
	return strings.Join(parts, ", ")
}

func PerformanceBand(score domain.Score) string {
	difference := score.Black - score.White
	if difference < 0 {
		difference = -difference
	}
	if difference >= 20 {
		return "decisive"
	}
	if difference >= 8 {
		return "close"
	}
	return "balanced"
}

func GroupByWinner(matches []domain.MatchRecord) map[domain.Disc][]domain.MatchRecord {
	grouped := map[domain.Disc][]domain.MatchRecord{domain.Black: {}, domain.White: {}, domain.Empty: {}}
	for _, match := range matches {
		grouped[match.Winner] = append(grouped[match.Winner], match)
	}
	return grouped
}
