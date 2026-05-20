package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RecoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				requestID := c.Writer.Header().Get("X-Request-ID")
				if requestID == "" {
					if v, ok := c.Get("request_id"); ok {
						if value, ok := v.(string); ok {
							requestID = value
						}
					}
				}

				log.Printf("panic recovered request_id=%s method=%s path=%s err=%v", requestID, c.Request.Method, c.Request.URL.Path, rec)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error":      "internal_server_error",
					"request_id": requestID,
				})
			}
		}()

		c.Next()
	}
}
