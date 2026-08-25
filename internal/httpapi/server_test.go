package httpapi

import (
	"aromaatelier/internal/flow031"
	"aromaatelier/internal/store"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/http.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	New(flow031.New(db)).Handler().ServeHTTP(rec, req)
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "ok") {
		t.Fatalf("code=%d body=%s", rec.Code, rec.Body.String())
	}
}
