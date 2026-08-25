package httpapi

import (
	"aromaatelier/internal/model"
	"net/url"
	"strconv"
	"strings"
)

func queryFromURL(values url.Values) model.SearchQuery {
	limit := 20
	if raw := values.Get("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 && parsed < 100 {
			limit = parsed
		}
	}
	return model.SearchQuery{Text: strings.TrimSpace(values.Get("q")), Status: model.RecordStatus(strings.TrimSpace(values.Get("status"))), Tag: strings.TrimSpace(values.Get("tag")), Limit: limit}
}

func contentType(format string) string {
	switch strings.ToLower(format) {
	case "csv":
		return "text/csv"
	case "markdown", "md":
		return "text/markdown"
	default:
		return "application/json"
	}
}

func validAction(action string) bool {
	switch action {
	case "submit", "approve", "archive", "address":
		return true
	default:
		return false
	}
}

func splitPath(path string) []string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	result := []string{}
	for _, part := range parts {
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}
