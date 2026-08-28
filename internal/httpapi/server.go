package httpapi

import (
	"aromaatelier/internal/flow031"
	"aromaatelier/internal/model"
	"encoding/json"
	"net/http"
	"strings"
)

type Server struct {
	handler *flow031.Handler
	mux     *http.ServeMux
}

func New(handler *flow031.Handler) *Server {
	s := &Server{handler: handler, mux: http.NewServeMux()}
	s.mux.HandleFunc("/health", s.health)
	s.mux.HandleFunc("/records", s.records)
	s.mux.HandleFunc("/records/", s.record)
	return s
}

func (s *Server) Handler() http.Handler { return logging(s.mux) }

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) records(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rows, err := s.handler.Search(queryFromURL(r.URL.Query()))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, rows)
	case http.MethodPost:
		var record model.Record
		if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
			writeErrorStatus(w, err, http.StatusBadRequest)
			return
		}
		created, err := s.handler.Register(record, actor(r))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, created)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) record(w http.ResponseWriter, r *http.Request) {
	parts := splitPath(r.URL.Path)
	id := ""
	if len(parts) >= 2 {
		id = parts[1]
	}
	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodPost:
		var body struct {
			Action  string `json:"action"`
			Note    string `json:"note"`
			Address string `json:"address"`
			Reason  string `json:"reason"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeErrorStatus(w, err, http.StatusBadRequest)
			return
		}
		var result model.Record
		var err error
		if !validAction(body.Action) {
			writeErrorStatus(w, http.ErrNotSupported, http.StatusBadRequest)
			return
		}
		switch body.Action {
		case "submit":
			result, err = s.handler.Submit(id, actor(r))
		case "approve":
			result, err = s.handler.Approve(id, actor(r), body.Note)
		case "archive":
			result, err = s.handler.Archive(id, actor(r))
		case "address":
			result, err = s.handler.UpdateAddress(id, actor(r), body.Address, body.Reason)
		default:
			writeErrorStatus(w, http.ErrNotSupported, http.StatusBadRequest)
			return
		}
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, result)
	case http.MethodGet:
		payload, err := s.handler.ExportOne(id)
		if err != nil {
			writeError(w, err)
			return
		}
		w.Header().Set("Content-Type", contentType(r.URL.Query().Get("format")))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(payload)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func actor(r *http.Request) string {
	value := r.Header.Get("X-Actor")
	if strings.TrimSpace(value) == "" {
		return "anonymous"
	}
	return value
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, err error) {
	writeErrorStatus(w, err, http.StatusUnprocessableEntity)
}

func writeErrorStatus(w http.ResponseWriter, err error, status int) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { withRequestID(w, r); next.ServeHTTP(w, r) })
}
