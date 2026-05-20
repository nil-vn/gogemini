package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gogemini/internal/config"
)

func TestHealthzReturnsOKWhenDBIsUp(t *testing.T) {
	db, err := sql.Open("sqlite", "file:test-health.db?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	r := NewRouter(config.Config{CORSOrigin: "*"}, db)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("parse response: %v", err)
	}
	if body["status"] != "ready" {
		t.Fatalf("unexpected body: %+v", body)
	}
}

func TestHealthzReturns503WhenDBIsDown(t *testing.T) {
	db, err := sql.Open("sqlite", "file:test-health-down.db?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	db.Close()

	r := NewRouter(config.Config{CORSOrigin: "*"}, db)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", w.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("parse response: %v", err)
	}
	errObj, ok := body["error"].(map[string]any)
	if !ok || errObj["code"] == "" {
		t.Fatalf("expected error envelope: %+v", body)
	}
}
