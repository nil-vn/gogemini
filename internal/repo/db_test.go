package repo

import (
	"database/sql"
	"testing"

	"gogemini/internal/config"
	"gogemini/internal/service"

	_ "modernc.org/sqlite"
)

func setupRepoTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", "file:test-repo.db?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	if err := ensureSchema(db); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestEnsureDefaultAdminCreatesNewUser(t *testing.T) {
	db := setupRepoTestDB(t)
	defer db.Close()

	cfg := config.Config{
		DefaultAdminUsername: "root",
		DefaultAdminPassword: "secret123",
		DefaultAdminEmail:    "root@example.com",
		DefaultAdminRole:     "admin",
		DefaultAdminStatus:   "active",
	}

	if err := ensureDefaultAdmin(db, cfg); err != nil {
		t.Fatalf("ensureDefaultAdmin: %v", err)
	}

	var username, email, role, status, hash string
	if err := db.QueryRow(`SELECT username, email, role, status, password_hash FROM users WHERE username = ?`, "root").Scan(&username, &email, &role, &status, &hash); err != nil {
		t.Fatalf("query user: %v", err)
	}

	if username != "root" || email != "root@example.com" || role != "admin" || status != "active" {
		t.Fatalf("unexpected user data username=%s email=%s role=%s status=%s", username, email, role, status)
	}
	if !service.CheckWerkzeugPasswordHash(hash, "secret123") {
		t.Fatalf("stored hash does not match password")
	}
}

func TestEnsureDefaultAdminUpdatesMatchedByEmailAndUsername(t *testing.T) {
	db := setupRepoTestDB(t)
	defer db.Close()

	if _, err := db.Exec(`INSERT INTO users (username, role, email, password_hash, status) VALUES ('oldadmin', 'editor', 'admin@example.com', 'pbkdf2:sha256:260000$salt$JNgGCs8M26YUVygOUxb8y4Na0Irn5m9wGLZi4/4+rF4=', 'inactive')`); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	cfg := config.Config{
		DefaultAdminUsername: "newadmin",
		DefaultAdminPassword: "newpass123",
		DefaultAdminEmail:    "admin@example.com",
		DefaultAdminRole:     "admin",
		DefaultAdminStatus:   "active",
	}

	if err := ensureDefaultAdmin(db, cfg); err != nil {
		t.Fatalf("ensureDefaultAdmin: %v", err)
	}

	var username, email, role, status, hash string
	if err := db.QueryRow(`SELECT username, email, role, status, password_hash FROM users WHERE email = ?`, "admin@example.com").Scan(&username, &email, &role, &status, &hash); err != nil {
		t.Fatalf("query updated user: %v", err)
	}

	if username != "newadmin" {
		t.Fatalf("username not updated, got %s", username)
	}
	if role != "admin" || status != "active" || email != "admin@example.com" {
		t.Fatalf("unexpected updated fields role=%s status=%s email=%s", role, status, email)
	}
	if !service.CheckWerkzeugPasswordHash(hash, "newpass123") {
		t.Fatalf("stored hash does not match updated password")
	}
}
