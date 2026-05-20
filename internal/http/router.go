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
	r.Use(middleware.RequestID(), middleware.RecoverPanic(), middleware.RequestLogger(), middleware.CORS(cfg.CORSOrigin))
	registerAdminRoutes(r, db, cfg)

	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	r.GET("/livez", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "live"}) })
	r.GET("/readyz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			writeError(c, http.StatusServiceUnavailable, "DEPENDENCY_DB_DOWN", "database dependency not ready")
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready", "dependencies": gin.H{"db": "up"}})
	})
	return r
}
