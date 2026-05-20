package http

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"gogemini/internal/config"
	"gogemini/internal/http/middleware"
)

func NewRouter(cfg config.Config, db *sql.DB) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.RecoverPanic())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.CORS(cfg.CORSOrigin))

	registerAdminRoutes(r, db, cfg)

	r.GET("/healthz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "app": "up", "db": "down", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "app": "up", "db": "up"})
	})

	return r
}
