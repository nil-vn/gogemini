package repo

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"gogemini/internal/config"
	"gogemini/internal/service"

	_ "modernc.org/sqlite"
)

func MustOpen(cfg config.Config) *sql.DB {
	db, err := sql.Open(cfg.DBDriver, cfg.DBDSN)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	db.SetMaxOpenConns(cfg.DBMaxOpen)
	db.SetMaxIdleConns(cfg.DBMaxIdle)
	db.SetConnMaxLifetime(time.Duration(cfg.DBMaxLifeMs) * time.Millisecond)

	if err := ensureSchema(db); err != nil {
		log.Fatalf("ensure schema: %v", err)
	}

	if err := ensureDefaultAdmin(db, cfg); err != nil {
		log.Fatalf("ensure default admin: %v", err)
	}

	return db
}

func ensureSchema(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT NOT NULL UNIQUE, role TEXT DEFAULT 'admin', email TEXT UNIQUE, password_hash TEXT NOT NULL, status TEXT DEFAULT 'active', created_date DATETIME DEFAULT CURRENT_TIMESTAMP);`,
		`CREATE TABLE IF NOT EXISTS config (id INTEGER PRIMARY KEY AUTOINCREMENT, key TEXT NOT NULL UNIQUE, value TEXT);`,
		`CREATE TABLE IF NOT EXISTS car (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, branch TEXT, model TEXT, vin TEXT, color TEXT, traded_company TEXT, imported_date TEXT, inspection_from TEXT, inspection_to TEXT, year_of_manufacture TEXT, purchase_price INTEGER DEFAULT 0, selling_price INTEGER DEFAULT 0, status TEXT, note TEXT, license_plate_no TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP, car_situation TEXT);`,
		`CREATE TABLE IF NOT EXISTS customer (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, gender TEXT, birth_day TEXT, facebook TEXT, phone TEXT, address TEXT, license_img TEXT, gallery_id INTEGER, lead_source TEXT, status TEXT, note TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP);`,
		"CREATE TABLE IF NOT EXISTS `transaction` (id INTEGER PRIMARY KEY AUTOINCREMENT, purchase_date TEXT, selling_price INTEGER DEFAULT 0, status TEXT, note TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP, customer_id INTEGER, deposit_amount INTEGER DEFAULT 0, FOREIGN KEY(customer_id) REFERENCES customer(id));",
		`CREATE TABLE IF NOT EXISTS transaction_item (id INTEGER PRIMARY KEY AUTOINCREMENT, transaction_id INTEGER NOT NULL, name TEXT NOT NULL, price INTEGER DEFAULT 0, FOREIGN KEY(transaction_id) REFERENCES ` + "`transaction`" + `(id) ON DELETE CASCADE);`,
		`CREATE TABLE IF NOT EXISTS transaction_car (transaction_id INTEGER NOT NULL, car_id INTEGER NOT NULL, PRIMARY KEY(transaction_id, car_id), FOREIGN KEY(transaction_id) REFERENCES ` + "`transaction`" + `(id) ON DELETE CASCADE, FOREIGN KEY(car_id) REFERENCES car(id) ON DELETE CASCADE);`,
		`CREATE TABLE IF NOT EXISTS car_image (id INTEGER PRIMARY KEY AUTOINCREMENT, car_id INTEGER NOT NULL, file_path TEXT NOT NULL, created_at DATETIME DEFAULT CURRENT_TIMESTAMP, FOREIGN KEY(car_id) REFERENCES car(id) ON DELETE CASCADE);`,
		`CREATE TABLE IF NOT EXISTS customer_image (id INTEGER PRIMARY KEY AUTOINCREMENT, customer_id INTEGER NOT NULL, file_path TEXT NOT NULL, created_at DATETIME DEFAULT CURRENT_TIMESTAMP, FOREIGN KEY(customer_id) REFERENCES customer(id) ON DELETE CASCADE);`,
	}
	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}

func ensureDefaultAdmin(db *sql.DB, cfg config.Config) error {
	username := strings.TrimSpace(cfg.DefaultAdminUsername)
	password := strings.TrimSpace(cfg.DefaultAdminPassword)
	if username == "" || password == "" {
		return nil
	}

	email := strings.TrimSpace(cfg.DefaultAdminEmail)
	if email == "" {
		email = username + "@local"
	}
	role := strings.TrimSpace(cfg.DefaultAdminRole)
	if role == "" {
		role = "admin"
	}
	status := strings.TrimSpace(cfg.DefaultAdminStatus)
	if status == "" {
		status = "active"
	}

	hash, err := service.BuildWerkzeugPasswordHash(password)
	if err != nil {
		return err
	}

	var id int64
	err = db.QueryRow(`SELECT id FROM users WHERE username = ? OR email = ? LIMIT 1`, username, email).Scan(&id)
	if err == sql.ErrNoRows {
		_, err = db.Exec(`INSERT INTO users (username, role, email, password_hash, status) VALUES (?, ?, ?, ?, ?)`, username, role, email, hash, status)
		return err
	}
	if err != nil {
		return err
	}
	_, err = db.Exec(`UPDATE users SET username=?, email=?, password_hash=?, role=?, status=? WHERE id=?`, username, email, hash, role, status, id)
	return err
}
