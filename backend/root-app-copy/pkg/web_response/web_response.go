package web_response

import (
	"root-app/internal/exception" // Adjust import path if needed
	"net/http"

	"github.com/gin-gonic/gin"
)

// WebResponse is the standardized structure for all API responses.
type WebResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}

// ErrorDetail is the standardized structure for the nested error object.
type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Success sends a standardized success response.
func RespondWithSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, WebResponse{
		Code:   code,
		Status: http.StatusText(code),
		Data:   data,
	})
}

// Error sends a standardized error response from an AppError.
func HandleAppError(c *gin.Context, err *exception.AppError) {
	status := err.HTTPStatus()

	errorDetail := ErrorDetail{
		Code:    string(err.Code),
		Message: err.Message,
		Details: err.Details,
	}

	c.AbortWithStatusJSON(status, WebResponse{
		Code:   status,
		Status: http.StatusText(status),
		Error:  errorDetail,
	})
}