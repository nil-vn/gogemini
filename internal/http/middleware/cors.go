package middleware

import "github.com/gin-gonic/gin"

func CORS(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestOrigin := c.Request.Header.Get("Origin")
		allowOrigin := origin
		if origin == "*" && requestOrigin != "" {
			allowOrigin = requestOrigin
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Vary", "Origin")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
