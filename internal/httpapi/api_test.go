package httpapi

import (
	"context"
	"example.com/othello-records/internal/matches"
	"example.com/othello-records/internal/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthEndpointReturnsServiceStatus(t *testing.T) {
	database, err := storage.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	if err := database.Initialize(context.Background()); err != nil {
		t.Fatal(err)
	}
	players := storage.NewPlayerRepository(database)
	service := matches.NewService(storage.NewMatchRepository(database), players, matches.NewHistoryWriter(storage.NewMatchRepository(database), storage.NewSnapshotRepository(database), storage.NewEventRepository(database)))
	server := NewServer(service, matches.NewPlayerService(players), matches.NewQueryService(service))
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d", response.Code)
	}
	if !strings.Contains(response.Body.String(), "ok") {
		t.Fatalf("body = %s", response.Body.String())
	}
}
