package httpapi

import (
	"example.com/othello-records/internal/domain"
	"example.com/othello-records/internal/reporting"
	"net/http"
	"strconv"
)

func (s *Server) matchDetailsEndpoint(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	id := request.URL.Query().Get("id")
	if id == "" {
		writeError(writer, http.StatusBadRequest, "id is required")
		return
	}
	match, err := s.matches.Get(request.Context(), id)
	if err != nil {
		writeError(writer, http.StatusNotFound, err.Error())
		return
	}
	frames, frameErr := s.matches.Replay(request.Context(), id)
	if frameErr != nil {
		writeError(writer, http.StatusInternalServerError, frameErr.Error())
		return
	}
	writeJSON(writer, http.StatusOK, map[string]any{"match": reporting.BuildMatchReport(match), "frames": frames})
}

func (s *Server) eventsEndpoint(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	id := request.URL.Query().Get("match")
	events, err := s.matches.Events(request.Context(), id)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(writer, http.StatusOK, events)
}

func parseLimit(request *http.Request) int {
	value, err := strconv.Atoi(request.URL.Query().Get("limit"))
	if err != nil {
		return 50
	}
	if value < 1 {
		return 1
	}
	if value > 500 {
		return 500
	}
	return value
}

func filterFromRequest(request *http.Request) domain.MatchFilter {
	return domain.MatchFilter{PlayerID: request.URL.Query().Get("player"), Status: domain.MatchStatus(request.URL.Query().Get("status")), Limit: parseLimit(request)}
}
