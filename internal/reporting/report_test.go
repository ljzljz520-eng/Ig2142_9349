package reporting

import (
	"example.com/othello-records/internal/domain"
	"testing"
)

func TestLeaderboardOrdersByPoints(t *testing.T) {
	players := []domain.PlayerProfile{{ID: "a", DisplayName: "Alpha", Games: 2, Wins: 2}, {ID: "b", DisplayName: "Beta", Games: 2, Wins: 1, Draws: 1}}
	entries := BuildLeaderboard(players)
	if len(entries) != 2 || entries[0].PlayerID != "a" {
		t.Fatalf("entries = %#v", entries)
	}
	if WinnerText(domain.Black) != "Black wins" {
		t.Fatal("winner text mismatch")
	}
}
