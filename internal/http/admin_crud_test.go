package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gogemini/internal/config"
	"gogemini/internal/service"
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
		"CREATE TABLE `transaction` (id INTEGER PRIMARY KEY AUTOINCREMENT, customer_id INTEGER, status TEXT, selling_price INTEGER, purchase_date TEXT, note TEXT, created_at DATETIME);",
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
	req.AddCookie(&http.Cookie{Name: "session", Value: service.BuildSessionToken(1, "test-secret")})
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestAdminCRUDFlows(t *testing.T) {
	db := setupAdminTestDB(t)
	defer db.Close()
	r := NewRouter(config.Config{CORSOrigin: "*", AuthSecret: "test-secret", Environment: "test"}, db)

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

func TestAdminSearchParity(t *testing.T) {
	db := setupAdminTestDB(t)
	defer db.Close()
	r := NewRouter(config.Config{CORSOrigin: "*", AuthSecret: "test-secret", Environment: "test"}, db)

	seed := []string{
		`INSERT INTO users (username, role, email, password_hash, status) VALUES ('search_user','admin','search@example.com','x','active');`,
		`INSERT INTO car (name, branch, model, vin, status, car_situation, selling_price) VALUES ('Search Car','Honda','Type R','VIN-SEARCH','available','new',100);`,
		`INSERT INTO customer (name, phone, address, status) VALUES ('Search Customer','000','HCM','active');`,
		"INSERT INTO `transaction` (customer_id, status, selling_price, purchase_date, note) VALUES (1,'search-status',200,'2026-01-01','search note');",
	}
	for _, s := range seed {
		if _, err := db.Exec(s); err != nil {
			t.Fatal(err)
		}
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, authedReq(http.MethodGet, "/api/admin/search?q=", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("empty search code=%d body=%s", w.Code, w.Body.String())
	}
	var emptyResp map[string][]map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &emptyResp); err != nil {
		t.Fatal(err)
	}
	for module, items := range emptyResp {
		if len(items) != 0 {
			t.Fatalf("expected empty results for module=%s when q is empty, got=%d", module, len(items))
		}
	}

	w = httptest.NewRecorder()
	r.ServeHTTP(w, authedReq(http.MethodGet, "/api/admin/search?q=search", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("search code=%d body=%s", w.Code, w.Body.String())
	}
	var searchResp map[string][]map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &searchResp); err != nil {
		t.Fatal(err)
	}
	for _, module := range []string{"users", "cars", "customers", "transactions"} {
		if len(searchResp[module]) == 0 {
			t.Fatalf("expected non-empty search result for module=%s", module)
		}
	}
}

func TestAdminDashboardMetricsParity(t *testing.T) {
	db := setupAdminTestDB(t)
	defer db.Close()
	r := NewRouter(config.Config{CORSOrigin: "*", AuthSecret: "test-secret", Environment: "test"}, db)

	now := time.Now().UTC()
	currentMonth := time.Date(now.Year(), now.Month(), 15, 10, 0, 0, 0, time.UTC).Format(time.RFC3339)
	prevMonth := time.Date(now.Year(), now.Month(), 15, 10, 0, 0, 0, time.UTC).AddDate(0, -1, 0).Format(time.RFC3339)

	seed := []string{
		`INSERT INTO car (name, branch, model, vin, status, car_situation, selling_price) VALUES ('C4 Car','Toyota','Cross','VIN-C4','available','new',100);`,
		`INSERT INTO customer (name, phone, address, status) VALUES ('C4 Customer','111','HN','active');`,
		fmt.Sprintf("INSERT INTO `transaction` (customer_id, status, selling_price, purchase_date, note, created_at) VALUES (1,'PAID',1200,'2026-01-01','seed paid','%s');", currentMonth),
		fmt.Sprintf("INSERT INTO `transaction` (customer_id, status, selling_price, purchase_date, note, created_at) VALUES (1,'DEPOSITED',800,'2026-01-02','seed deposited','%s');", prevMonth),
		fmt.Sprintf("INSERT INTO `transaction` (customer_id, status, selling_price, purchase_date, note, created_at) VALUES (1,'CANCELLED',5000,'2026-01-03','excluded','%s');", currentMonth),
	}
	for _, s := range seed {
		if _, err := db.Exec(s); err != nil {
			t.Fatal(err)
		}
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, authedReq(http.MethodGet, "/api/admin/dashboard", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("dashboard code=%d body=%s", w.Code, w.Body.String())
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if int(resp["users"].(float64)) != 1 || int(resp["cars"].(float64)) != 1 || int(resp["customers"].(float64)) != 1 || int(resp["transactions"].(float64)) != 3 {
		t.Fatalf("unexpected count metrics: %+v", resp)
	}
	if int(resp["total_revenue"].(float64)) != 2000 {
		t.Fatalf("expected total_revenue=2000, got=%v", resp["total_revenue"])
	}
	if _, ok := resp["revenue_"+time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format("01/2006")]; !ok {
		t.Fatalf("missing current month revenue key: %+v", resp)
	}
}

func itoa(v int64) string { return fmt.Sprintf("%d", v) }
