package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gogemini/internal/config"
)

func TestAPIParityCriticalEndpoints(t *testing.T) {
	db := setupAdminTestDB(t)
	defer db.Close()
	r := NewRouter(config.Config{CORSOrigin: "*", AuthSecret: "test-secret", Environment: "test"}, db)

	seed := []string{
		`INSERT INTO car (name, branch, model, vin, status, car_situation, selling_price) VALUES ('Parity Car','Honda','Civic','VIN-E2','available','new',100);`,
		`INSERT INTO customer (name, phone, address, status) VALUES ('Parity Customer','123','HCM','active');`,
		"INSERT INTO `transaction` (customer_id, status, selling_price, purchase_date, note, created_at) VALUES (1,'PAID',500,'2026-01-01','parity',CURRENT_TIMESTAMP);",
	}
	for _, s := range seed {
		if _, err := db.Exec(s); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("dashboard response shape", func(t *testing.T) {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, authedReq(http.MethodGet, "/api/admin/dashboard", nil))
		if w.Code != http.StatusOK {
			t.Fatalf("unexpected code=%d body=%s", w.Code, w.Body.String())
		}
		var body map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		for _, key := range []string{"users", "cars", "customers", "transactions", "total_revenue"} {
			if _, ok := body[key]; !ok {
				t.Fatalf("missing key %s in dashboard response: %+v", key, body)
			}
		}
	})

	t.Run("list envelope parity", func(t *testing.T) {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, authedReq(http.MethodGet, "/api/admin/users?page=1&page_size=10&sort=id&order=asc", nil))
		if w.Code != http.StatusOK {
			t.Fatalf("unexpected code=%d body=%s", w.Code, w.Body.String())
		}
		var body map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		for _, key := range []string{"items", "total", "page", "page_size", "sort", "order"} {
			if _, ok := body[key]; !ok {
				t.Fatalf("missing key %s in list response: %+v", key, body)
			}
		}
	})

	t.Run("search empty q returns empty results", func(t *testing.T) {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, authedReq(http.MethodGet, "/api/admin/search?q=", nil))
		if w.Code != http.StatusOK {
			t.Fatalf("unexpected code=%d body=%s", w.Code, w.Body.String())
		}
		var body map[string][]map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		for module, items := range body {
			if len(items) != 0 {
				t.Fatalf("module=%s expected empty results, got=%d", module, len(items))
			}
		}
	})

	t.Run("system settings payload parity", func(t *testing.T) {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, authedReq(http.MethodGet, "/api/admin/system", nil))
		if w.Code != http.StatusOK {
			t.Fatalf("unexpected code=%d body=%s", w.Code, w.Body.String())
		}
		var body map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		for _, key := range []string{"currency", "theme", "language"} {
			if _, ok := body[key]; !ok {
				t.Fatalf("missing key %s in system settings response: %+v", key, body)
			}
		}
	})
}
