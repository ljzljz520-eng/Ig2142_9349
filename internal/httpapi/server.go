package httpapi

import (
	"context"
	"example.com/othello-records/internal/matches"
	"net/http"
)

type Server struct {
	mux     *http.ServeMux
	matches *matches.MatchService
	players *matches.PlayerService
	queries *matches.QueryService
}

func NewServer(matchService *matches.MatchService, playerService *matches.PlayerService, queryService *matches.QueryService) *Server {
	server := &Server{mux: http.NewServeMux(), matches: matchService, players: playerService, queries: queryService}
	server.routes()
	return server
}

func (s *Server) Handler() http.Handler { return s.logging(s.mux) }

func (s *Server) routes() {
	s.mux.HandleFunc("/health", s.health)
	s.mux.HandleFunc("/players", s.playersEndpoint)
	s.mux.HandleFunc("/matches", s.matchesEndpoint)
	s.mux.HandleFunc("/moves", s.moveEndpoint)
	s.mux.HandleFunc("/history", s.historyEndpoint)
	s.mux.HandleFunc("/match", s.matchDetailsEndpoint)
	s.mux.HandleFunc("/events", s.eventsEndpoint)
}

func (s *Server) health(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) Serve(ctx context.Context, address string) error {
	server := &http.Server{Addr: address, Handler: s.Handler()}
	go func() { <-ctx.Done(); _ = server.Shutdown(context.Background()) }()
	return server.ListenAndServe()
}

func (s *Server) logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("X-Othello-Service", "history")
		next.ServeHTTP(writer, request)
	})
}
