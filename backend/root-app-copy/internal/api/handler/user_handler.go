package handler

import (
	"root-app/internal/api/dto"
	"root-app/internal/exception"
	"root-app/internal/contract"
	"root-app/pkg/gin_helper"
	"root-app/pkg/web_response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService contract.UserService
}

func NewAuthHandler(authService contract.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request body", err.Error())
		web_response.HandleAppError(c, appErr)
		return
	}

    // Get client IP address
    ipAddress := c.ClientIP()

    // Pass ipAddress to the service
	profile, appErr := h.authService.Register(c.Request.Context(), req, ipAddress)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}

	web_response.RespondWithSuccess(c, http.StatusCreated, profile)
}

// Login godoc
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request body", err.Error())
		web_response.HandleAppError(c, appErr)
		return
	}

    // Get client IP address
    ipAddress := c.ClientIP()

    // Pass ipAddress to the service
	tokenResponse, appErr := h.authService.Login(c.Request.Context(), req, ipAddress)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}

	web_response.RespondWithSuccess(c, http.StatusOK, tokenResponse)
}

// Logout godoc
// @Router /auth/logout [post] 
func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request: missing refresh_token in body", err.Error())
		web_response.HandleAppError(c, appErr)
		return
	}

	ipAddress := c.ClientIP()
	appErr := h.authService.Logout(c.Request.Context(), req.RefreshToken, ipAddress)
	
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}

	web_response.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Logout successful"})
}

// refresh handler
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	// 1. Bind the incoming JSON request to the RefreshRequest DTO.
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request: missing refresh_token in body", err.Error())
		web_response.HandleAppError(c, appErr)
		return
	}

	// 3. Call the auth service to perform the token refresh logic.
	tokenResponse, appErr := h.authService.RefreshToken(c.Request.Context(), req)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}

	web_response.RespondWithSuccess(c, http.StatusOK, tokenResponse)
}

// GetProfile retrieves the authenticated user's profile.
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, appErr := gin_helper.GetUserIDFromContext(c)
	if appErr != nil {
		// GetUserIDFromContext already calls Abort, but we use our handler for a consistent response body
		web_response.HandleAppError(c, appErr) // CORRECTED
		return
	}

	profile, appErr := h.authService.GetUserProfile(c.Request.Context(), userID)
	if appErr != nil {
		web_response.HandleAppError(c, appErr) // CORRECTED
		return
	}

	web_response.RespondWithSuccess(c, http.StatusOK, profile) // CORRECTED
}

// UpdateProfile handles updating the authenticated user's profile.
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, appErr := gin_helper.GetUserIDFromContext(c)
	if appErr != nil {
		web_response.HandleAppError(c, appErr) // CORRECTED
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := exception.NewValidationError("Invalid request body", err.Error())
		web_response.HandleAppError(c, appErr) // CORRECTED
		return
	}

	profile, appErr := h.authService.UpdateUserProfile(c.Request.Context(), userID, req)
	if appErr != nil {
		web_response.HandleAppError(c, appErr) // CORRECTED
		return
	}

	web_response.RespondWithSuccess(c, http.StatusOK, profile) // CORRECTED
}
