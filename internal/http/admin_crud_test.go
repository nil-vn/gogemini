package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gogemini/internal/config"
)

func setupAdminTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", "file:test-admin.db?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	stmts := []string{
		`CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, role TEXT, email TEXT, password_hash TEXT, status TEXT);`,
		`CREATE TABLE car (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, branch TEXT, model TEXT, vin TEXT, status TEXT, car_situation TEXT, selling_price INTEGER);`,
		`CREATE TABLE customer (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, phone TEXT, address TEXT, status TEXT);`,
		"CREATE TABLE `transaction` (id INTEGER PRIMARY KEY AUTOINCREMENT, customer_id INTEGER, status TEXT, selling_price INTEGER, purchase_date TEXT, note TEXT);",
		`CREATE TABLE config (key TEXT PRIMARY KEY, value TEXT);`,
		`INSERT INTO users (username, role, email, password_hash, status) VALUES ('admin','admin','a@a.com','pbkdf2:sha256:260000$salt$hash','active');`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			t.Fatal(err)
		}
	}
	return db
}

func authedReq(method, path string, body []byte) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.AddCookie(&http.Cookie{Name: "session", Value: "ok"})
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestAdminCRUDFlows(t *testing.T) {
	db := setupAdminTestDB(t)
	defer db.Close()
	r := NewRouter(config.Config{CORSOrigin: "*"}, db)

	cases := []struct{ name, base, create, update string }{
		{"users", "/api/admin/users", `{"username":"u1","password_hash":"x"}`, `{"username":"u2"}`},
		{"cars", "/api/admin/cars", `{"name":"c1"}`, `{"name":"c2"}`},
		{"customers", "/api/admin/customers", `{"name":"cus1"}`, `{"name":"cus2"}`},
		{"transactions", "/api/admin/transactions", `{"customer_id":1,"status":"new"}`, `{"customer_id":1,"status":"done"}`},
	}

	for _, tc := range cases {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, authedReq(http.MethodPost, tc.base, []byte(tc.create)))
		if w.Code != http.StatusCreated {
			t.Fatalf("%s create code=%d body=%s", tc.name, w.Code, w.Body.String())
		}
		var created map[string]int64
		_ = json.Unmarshal(w.Body.Bytes(), &created)
		id := created["id"]

		w = httptest.NewRecorder()
		r.ServeHTTP(w, authedReq(http.MethodGet, tc.base+"?page=1&page_size=10", nil))
		if w.Code != http.StatusOK {
			t.Fatalf("%s list code=%d", tc.name, w.Code)
		}

		w = httptest.NewRecorder()
		r.ServeHTTP(w, authedReq(http.MethodGet, tc.base+"/"+itoa(id), nil))
		if w.Code != http.StatusOK {
			t.Fatalf("%s detail code=%d", tc.name, w.Code)
		}

		w = httptest.NewRecorder()
		r.ServeHTTP(w, authedReq(http.MethodPut, tc.base+"/"+itoa(id), []byte(tc.update)))
		if w.Code != http.StatusNoContent {
			t.Fatalf("%s update code=%d body=%s", tc.name, w.Code, w.Body.String())
		}

		w = httptest.NewRecorder()
		r.ServeHTTP(w, authedReq(http.MethodDelete, tc.base+"/"+itoa(id), nil))
		if w.Code != http.StatusNoContent {
			t.Fatalf("%s delete code=%d", tc.name, w.Code)
		}
	}
}

func itoa(v int64) string { return fmt.Sprintf("%d", v) }
