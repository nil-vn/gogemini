package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		if query != "" {
			path += "?" + query
		}

		c.Next()
		sample := c.Writer.Status() >= 500 || rand.Intn(100) < 20
		if !sample {
			return
		}

		requestID := c.Writer.Header().Get("X-Request-ID")
		logger.Info("http_request",
			slog.String("request_id", requestID),
			slog.String("method", c.Request.Method),
			slog.String("path", redactPII(path)),
			slog.Int("status", c.Writer.Status()),
			slog.Int64("latency_ms", time.Since(start).Milliseconds()),
			slog.String("client_ip_hash", hashPII(c.ClientIP())),
		)
	}
}

func redactPII(v string) string {
	v = strings.ReplaceAll(v, "@", "[at]")
	return v
}

func hashPII(v string) string {
	h := sha256.Sum256([]byte(v))
	return hex.EncodeToString(h[:8])
}
