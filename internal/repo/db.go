package repo

import (
	"database/sql"
	"log"

	"gogemini/internal/config"

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

	return db
}
