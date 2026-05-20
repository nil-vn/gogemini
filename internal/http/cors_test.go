package http

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"gogemini/internal/config"
)

func TestCORSWildcardWithCredentialsReflectsRequestOrigin(t *testing.T) {
	db, err := sql.Open("sqlite", "file:test-cors.db?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	defer db.Close()

	r := NewRouter(config.Config{CORSOrigin: "*", AuthSecret: "test-secret", Environment: "test"}, db)

	req := httptest.NewRequest(http.MethodOptions, "/api/auth/login", nil)
	req.Header.Set("Origin", "http://192.168.1.51:5173")
	req.Header.Set("Access-Control-Request-Method", "POST")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204 preflight, got %d", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://192.168.1.51:5173" {
		t.Fatalf("allow-origin mismatch: %q", got)
	}
	if got := w.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Fatalf("allow-credentials mismatch: %q", got)
	}
}
