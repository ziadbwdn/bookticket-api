package middleware

import (
	"root-app/internal/exception"
	"root-app/internal/contract"
	"root-app/pkg/web_response"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware struct holds the authentication service dependency.
type AuthMiddleware struct {
	authService contract.UserService
}

// NewAuthMiddleware creates and returns a new instance of AuthMiddleware.
func NewAuthMiddleware(authService contract.UserService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// Handle returns a Gin middleware handler function.
func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			appErr := exception.NewAuthError("Authorization header required")
			web_response.HandleAppError(c, appErr)
			c.Abort() // Stop processing the request
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			appErr := exception.NewAuthError("Invalid authorization format. Expected 'Bearer <token>'")
			web_response.HandleAppError(c, appErr)
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Use the AuthService to verify the token
		userID, userRole, appErr := m.authService.VerifyToken(c.Request.Context(), tokenString)
		if appErr != nil {
			web_response.HandleAppError(c, appErr)
			c.Abort()
			return
		}

		// Set user context for downstream handlers
		c.Set("userID", userID)
		c.Set("userRole", userRole)
		c.Next() // Proceed to the next handler in the chain
	}
}