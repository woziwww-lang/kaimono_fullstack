package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Meta struct {
	Count  int `json:"count,omitempty"`
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type APIResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Meta  *Meta       `json:"meta,omitempty"`
	Error *APIError   `json:"error,omitempty"`
}

const (
	ErrInvalidArgument = "INVALID_ARGUMENT"
	ErrNotFound        = "NOT_FOUND"
	ErrInternal        = "INTERNAL_ERROR"
	ErrUnauthorized    = "UNAUTHORIZED"
)

func OK(c *gin.Context, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, APIResponse{Data: data, Meta: meta})
}

func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, APIResponse{Error: &APIError{Code: code, Message: message}})
}
