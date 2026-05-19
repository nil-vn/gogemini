package main

import (
	"log"

	"gogemini/internal/config"
	"gogemini/internal/http"
	"gogemini/internal/repo"
)

func main() {
	cfg := config.MustLoad()

	db := repo.MustOpen(cfg)
	defer db.Close()

	r := http.NewRouter(cfg, db)

	log.Printf("server starting on %s", cfg.ServerAddr)
	if err := r.Run(cfg.ServerAddr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
