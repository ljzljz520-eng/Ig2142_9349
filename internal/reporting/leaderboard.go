package reporting

import (
	"example.com/othello-records/internal/domain"
	"sort"
)

type LeaderboardEntry struct {
	PlayerID    string
	DisplayName string
	Games       int
	Wins        int
	WinRate     float64
	Points      int
}

func BuildLeaderboard(players []domain.PlayerProfile) []LeaderboardEntry {
	entries := make([]LeaderboardEntry, 0, len(players))
	for _, player := range players {
		entries = append(entries, LeaderboardEntry{PlayerID: player.ID, DisplayName: player.DisplayName, Games: player.Games,
			Wins: player.Wins, WinRate: winRate(player), Points: points(player)})
	}
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].Points == entries[j].Points {
			return entries[i].DisplayName < entries[j].DisplayName
		}
		return entries[i].Points > entries[j].Points
	})
	return entries
}

func winRate(player domain.PlayerProfile) float64 {
	if player.Games == 0 {
		return 0
	}
	return float64(player.Wins) / float64(player.Games)
}

func points(player domain.PlayerProfile) int { return player.Wins*3 + player.Draws }
