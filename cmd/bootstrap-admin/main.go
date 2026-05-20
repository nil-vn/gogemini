package main

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"flag"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
	"golang.org/x/crypto/pbkdf2"
)

func main() {
	dbDriver := flag.String("db-driver", "sqlite", "database driver")
	dbDSN := flag.String("db-dsn", "file:app.db?cache=shared", "database dsn")
	username := flag.String("username", "", "default admin username")
	password := flag.String("password", "", "default admin password")
	email := flag.String("email", "", "default admin email")
	role := flag.String("role", "admin", "default role")
	status := flag.String("status", "active", "default status")
	flag.Parse()

	if *username == "" || *password == "" {
		log.Fatal("username and password are required")
	}

	db, err := sql.Open(*dbDriver, *dbDSN)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if *email == "" {
		*email = fmt.Sprintf("%s@local", *username)
	}

	hash, err := buildWerkzeugPBKDF2Hash(*password)
	if err != nil {
		log.Fatalf("build password hash: %v", err)
	}

	var id int64
	err = db.QueryRow(`SELECT id FROM users WHERE username = ? OR email = ? LIMIT 1`, *username, *email).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		_, err = db.Exec(`INSERT INTO users (username, role, email, password_hash, status) VALUES (?, ?, ?, ?, ?)`, *username, *role, *email, hash, *status)
		if err != nil {
			log.Fatalf("insert default admin user: %v", err)
		}
		log.Printf("created default admin user: %s", *username)
	case err != nil:
		log.Fatalf("query existing user: %v", err)
	default:
		_, err = db.Exec(`UPDATE users SET password_hash=?, role=?, status=? WHERE id=?`, hash, *role, *status, id)
		if err != nil {
			log.Fatalf("update existing admin user: %v", err)
		}
		log.Printf("updated existing admin user credentials: %s", *username)
	}
}

func buildWerkzeugPBKDF2Hash(password string) (string, error) {
	const iterations = 260000
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	saltB64 := base64.RawURLEncoding.EncodeToString(salt)
	derived := pbkdf2.Key([]byte(password), []byte(saltB64), iterations, sha256.Size, sha256.New)
	hash := base64.StdEncoding.EncodeToString(derived)
	return fmt.Sprintf("pbkdf2:sha256:%d$%s$%s", iterations, saltB64, hash), nil
}
