package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gogemini/internal/config"
	ghttp "gogemini/internal/http"
	"gogemini/internal/repo"
)

func main() {
	cfg := config.MustLoad()
	db := repo.MustOpen(cfg)
	r := ghttp.NewRouter(cfg, db)
	srv := &http.Server{Addr: cfg.ServerAddr, Handler: r, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		log.Printf("server starting on %s", cfg.ServerAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	_ = db.Close()
}
