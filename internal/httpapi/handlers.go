package httpapi

import (
	"example.com/othello-records/internal/domain"
	"net/http"
)

type playerRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Rank string `json:"rank"`
}

type matchRequest struct {
	ID       string `json:"id"`
	Black    string `json:"black_player"`
	White    string `json:"white_player"`
	Label    string `json:"label"`
	Sequence int    `json:"sequence"`
}

type moveRequest struct {
	MatchID string `json:"match_id"`
	EventID string `json:"event_id"`
	Row     int    `json:"row"`
	Column  int    `json:"column"`
	Player  string `json:"player"`
}

func (s *Server) playersEndpoint(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		players, err := s.players.List(request.Context(), false)
		if err != nil {
			writeError(writer, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(writer, http.StatusOK, players)
	case http.MethodPost:
		var payload playerRequest
		if err := decodeJSON(request, &payload); err != nil {
			writeError(writer, http.StatusBadRequest, err.Error())
			return
		}
		player, err := s.players.Register(request.Context(), payload.ID, payload.Name, payload.Rank)
		if err != nil {
			writeError(writer, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(writer, http.StatusCreated, player)
	default:
		writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) matchesEndpoint(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		records, err := s.matches.History(request.Context(), domain.MatchFilter{})
		if err != nil {
			writeError(writer, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(writer, http.StatusOK, records)
	case http.MethodPost:
		var payload matchRequest
		if err := decodeJSON(request, &payload); err != nil {
			writeError(writer, http.StatusBadRequest, err.Error())
			return
		}
		match := domain.MatchRecord{ID: payload.ID, BlackPlayer: payload.Black, WhitePlayer: payload.White, Label: payload.Label, Sequence: payload.Sequence, Status: domain.StatusActive}
		if err := s.matches.Start(request.Context(), match); err != nil {
			writeError(writer, http.StatusBadRequest, err.Error())
			return
		}
		writeJSON(writer, http.StatusCreated, match)
	default:
		writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) historyEndpoint(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	filter := filterFromRequest(request)
	summaries, err := s.queries.History(request.Context(), filter)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(writer, http.StatusOK, summaries)
}

func (s *Server) moveEndpoint(writer http.ResponseWriter, request *http.Request) {
	var payload moveRequest
	if err := decodeJSON(request, &payload); err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	player, err := domain.ParseDisc(payload.Player)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	result, err := s.matches.SubmitMove(request.Context(), payload.MatchID, domain.Move{Row: payload.Row, Column: payload.Column, Player: player}, payload.EventID)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(writer, http.StatusOK, result)
}
