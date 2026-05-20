package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gogemini/internal/config"
)

func setupAuthDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", "file:test-auth.db?mode=memory&cache=shared")
	if err != nil { t.Fatal(err) }
	_, err = db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, role TEXT, email TEXT, password_hash TEXT, status TEXT);
INSERT INTO users (username, role, email, password_hash, status) VALUES ('admin','admin','admin@example.com','pbkdf2:sha256:260000$salt$JNgGCs8M26YUVygOUxb8y4Na0Irn5m9wGLZi4/4+rF4=','active');`)
	if err != nil { t.Fatal(err) }
	return db
}

func TestAuthLoginLogoutAndGuard(t *testing.T) {
	db := setupAuthDB(t)
	defer db.Close()
	r := NewRouter(config.Config{CORSOrigin:"*", AuthSecret:"test-secret", Environment:"test"}, db)

	body := []byte(`{"login":"admin","password":"password123"}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK { t.Fatalf("login code=%d body=%s", w.Code, w.Body.String()) }
	cookies := w.Result().Cookies()
	if len(cookies) == 0 || cookies[0].Name != "session" || !cookies[0].HttpOnly { t.Fatalf("missing hardened session cookie") }

	w = httptest.NewRecorder()
	adminReq := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
	adminReq.AddCookie(cookies[0])
	r.ServeHTTP(w, adminReq)
	if w.Code != http.StatusOK { t.Fatalf("authed admin code=%d", w.Code) }

	w = httptest.NewRecorder()
	logoutReq := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	r.ServeHTTP(w, logoutReq)
	if w.Code != http.StatusNoContent { t.Fatalf("logout code=%d", w.Code) }
}

func TestAuthLockoutAfterFailedAttempts(t *testing.T) {
	db := setupAuthDB(t)
	defer db.Close()
	r := NewRouter(config.Config{CORSOrigin:"*", AuthSecret:"test-secret", Environment:"test"}, db)

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		bad := map[string]string{"login":"admin", "password":"wrong"}
		b, _ := json.Marshal(bad)
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized { t.Fatalf("attempt %d expected 401 got %d", i+1, w.Code) }
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader([]byte(`{"login":"admin","password":"password123"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests { t.Fatalf("expected lockout 429 got %d", w.Code) }
}

func TestAuthLoginProductionOverHTTPStillSetsUsableSessionCookie(t *testing.T) {
	db := setupAuthDB(t)
	defer db.Close()
	r := NewRouter(config.Config{CORSOrigin:"*", AuthSecret:"test-secret", Environment:"production"}, db)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader([]byte(`{"login":"admin","password":"password123"}`)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK { t.Fatalf("login code=%d body=%s", w.Code, w.Body.String()) }

	cookies := w.Result().Cookies()
	if len(cookies) == 0 { t.Fatal("missing session cookie") }
	if cookies[0].Secure {
		t.Fatal("session cookie should not be Secure for plain HTTP requests")
	}

	w = httptest.NewRecorder()
	adminReq := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
	adminReq.AddCookie(cookies[0])
	r.ServeHTTP(w, adminReq)
	if w.Code != http.StatusOK { t.Fatalf("authed admin code=%d", w.Code) }
}
