package http

import "github.com/gin-gonic/gin"

type apiError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

type errorEnvelope struct {
	Error apiError `json:"error"`
}

func writeError(c *gin.Context, status int, code, message string) {
	rid, _ := c.Get("request_id")
	requestID, _ := rid.(string)
	c.AbortWithStatusJSON(status, errorEnvelope{Error: apiError{Code: code, Message: message, RequestID: requestID}})
}
