package gin_helper

import (
	"root-app/internal/exception"
	"root-app/internal/entities"
	"root-app/internal/utils"
	"root-app/pkg/web_response"
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext extracts the userID from Gin's context.
// It returns the BinaryUUID and an AppError if the userID is not found or is in an invalid format.

// Define context keys
const (
	UserIDContextKey   = "userID"
	UserRoleContextKey = "userRole" // New key for user role
	UsernameContextKey  = "username"
	IPAddressContextKey = "ipAddress"
)

func GetUserIDFromContext(c *gin.Context) (utils.BinaryUUID, *exception.AppError) {
	val, exists := c.Get(UserIDContextKey)
	if !exists {
		appErr := exception.NewAuthError("User ID not found in context. Authentication middleware missing or failed.")
		c.AbortWithStatus(http.StatusUnauthorized) // Or Internal Server Error if unexpected
		return utils.BinaryUUID{}, appErr
	}
	userID, ok := val.(utils.BinaryUUID)
	if !ok {
		appErr := exception.NewInternalError("Invalid user ID type in context.", nil)
		c.AbortWithStatus(http.StatusInternalServerError)
		return utils.BinaryUUID{}, appErr
	}
	return userID, nil
}

func ParseIDFromContext(c *gin.Context, paramName, resourceType string) (utils.BinaryUUID, *exception.AppError) {
	idStr := c.Param(paramName)
	if idStr == "" {
		appErr := exception.NewValidationError(fmt.Sprintf("%s ID is required", resourceType), fmt.Sprintf("Missing path parameter '%s'", paramName))
		web_response.HandleAppError(c, appErr)
		return utils.BinaryUUID{}, appErr
	}
	id, err := utils.ParseBinaryUUID(idStr)
	if err != nil {
		appErr := exception.NewValidationError(fmt.Sprintf("Invalid %s ID format", resourceType), err.Error())
		web_response.HandleAppError(c, appErr)
		return utils.BinaryUUID{}, appErr
	}
	return id, nil
}

// GetUserRoleFromContext extracts the user role from Gin's context.
func GetUserRoleFromContext(c *gin.Context) (entities.UserRole, *exception.AppError) {
	val, exists := c.Get(UserRoleContextKey)
	if !exists {
		appErr := exception.NewAuthError("User role not found in context. Authentication middleware missing or failed.")
		c.AbortWithStatus(http.StatusUnauthorized)
		return "", appErr
	}
	userRoleStr, ok := val.(string) // Role is stored as string by AuthMiddleware
	if !ok {
		appErr := exception.NewInternalError("Invalid user role type in context.", nil)
		c.AbortWithStatus(http.StatusInternalServerError)
		return "", appErr
	}
	userRole := entities.UserRole(userRoleStr) // Cast string back to UserRole
	return userRole, nil
}

func GetUsernameFromContext(c *gin.Context) (string, *exception.AppError) {
	val, exists := c.Get(UsernameContextKey)
	if !exists {
		// This is less critical than UserID, so we might not want to abort, but return an error.
		return "", exception.NewInternalError("Username not found in context", nil)
	}
	username, ok := val.(string)
	if !ok {
		return "", exception.NewInternalError("Invalid username type in context", nil)
	}
	return username, nil
}

// GetIPAddressFromContext extracts the IP address from Gin's context.
func GetIPAddressFromContext(c *gin.Context) string {
	val, exists := c.Get(IPAddressContextKey)
	if !exists {
		// Fallback to Gin's direct method if not set in context
		return c.ClientIP()
	}
	ip, ok := val.(string)
	if !ok {
		return c.ClientIP() // Fallback
	}
	return ip
}