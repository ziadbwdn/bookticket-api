package contract

import (
	"root-app/internal/api/dto"
	"root-app/internal/exception"
	"root-app/internal/entities"
	"root-app/internal/utils"
	"time"
	"context"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entities.User) *exception.AppError
	GetUserByID(ctx context.Context, id utils.BinaryUUID) (*entities.User, *exception.AppError)
	GetUserByUsername(ctx context.Context, username string) (*entities.User, *exception.AppError)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, *exception.AppError)
	UpdateUser(ctx context.Context, user *entities.User) *exception.AppError
	DeleteUser(ctx context.Context, id utils.BinaryUUID) *exception.AppError
	SaveRefreshToken(ctx context.Context, userID utils.BinaryUUID, tokenHash string, expiresAt time.Time) *exception.AppError
	// UpdateRefreshToken(ctx context.Context, userID utils.BinaryUUID, refreshTokenHash string, expiresAt time.Time) *exception.AppError
	ClearRefreshToken(ctx context.Context, userID utils.BinaryUUID) *exception.AppError
	GetUserByPasswordResetTokenHash(ctx context.Context, tokenHash string) (*entities.User, *exception.AppError)
	ClearPasswordResetToken(ctx context.Context, userID utils.BinaryUUID) *exception.AppError
	UpdateLastLogin(ctx context.Context, userID utils.BinaryUUID) *exception.AppError
	GrantPermission(ctx context.Context, userID, projectID utils.BinaryUUID, permission string) *exception.AppError
}

// AuthService defines the interface for authentication-related operations.
type UserService interface {
	// Add ipAddress parameter
	Register(ctx context.Context, req dto.RegisterRequest, ipAddress string) (*dto.ProfileResponse, *exception.AppError)
	// Add ipAddress parameter
	Login(ctx context.Context, req dto.LoginRequest, ipAddress string) (*dto.TokenResponse, *exception.AppError)

	VerifyToken(ctx context.Context, tokenString string) (utils.BinaryUUID, string, *exception.AppError)
	GetUserProfile(ctx context.Context, userID utils.BinaryUUID) (*dto.ProfileResponse, *exception.AppError)
	UpdateUserProfile(ctx context.Context, userID utils.BinaryUUID, req dto.UpdateProfileRequest) (*dto.ProfileResponse, *exception.AppError)
	GetUserDetailsForLogging(ctx context.Context, userID utils.BinaryUUID) (username string, appErr *exception.AppError)
	RefreshToken(ctx context.Context, req dto.RefreshRequest) (*dto.RefreshResponse, *exception.AppError)
	SendPasswordReset(ctx context.Context, email string) *exception.AppError
	ResetPassword(ctx context.Context, token, newPassword string) *exception.AppError
	VerifyEmail(ctx context.Context, token string) *exception.AppError
	Logout(ctx context.Context, token string, ipAddress string) *exception.AppError
}