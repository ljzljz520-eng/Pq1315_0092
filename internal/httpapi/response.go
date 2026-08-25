package httpapi

import (
	"aromaatelier/internal/model"
	"encoding/json"
	"net/http"
)

type envelope struct {
	Data  any            `json:"data,omitempty"`
	Error string         `json:"error,omitempty"`
	Meta  map[string]any `json:"meta,omitempty"`
}

func writeData(w http.ResponseWriter, status int, value any, meta map[string]any) {
	writeJSON(w, status, envelope{Data: value, Meta: meta})
}

func writeRecord(w http.ResponseWriter, status int, record model.Record) {
	writeData(w, status, record, map[string]any{"version": record.Version, "status": record.Status})
}

func decodeRecord(r *http.Request) (model.Record, error) {
	var record model.Record
	err := json.NewDecoder(r.Body).Decode(&record)
	return record, err
}

func requestID(r *http.Request) string {
	if value := r.Header.Get("X-Request-ID"); value != "" {
		return value
	}
	return "local-request"
}

func withRequestID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Request-ID", requestID(r))
}
