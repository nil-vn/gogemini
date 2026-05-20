package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		requestID := c.Writer.Header().Get("X-Request-ID")
		if requestID == "" {
			if v, ok := c.Get("request_id"); ok {
				if value, ok := v.(string); ok {
					requestID = value
				}
			}
		}

		if query != "" {
			path = path + "?" + query
		}

		log.Printf("request_id=%s method=%s path=%s status=%d latency=%s client_ip=%s", requestID, c.Request.Method, path, status, latency, c.ClientIP())
	}
}
