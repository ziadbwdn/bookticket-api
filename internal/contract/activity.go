package contract

import (
	"context"
	"time"

	"root-app/internal/api/dto"
	"root-app/internal/entities"
	"root-app/internal/utils"
)

// UserActivityRepository defines the interface for interacting with user activity data storage.
type UserActivityRepository interface {
	Create(ctx context.Context, activity *entities.Activity) error
	ListUserActivities(ctx context.Context, filter dto.ActivityFilterRequest) ([]entities.Activity, int64, error)
	GetByUserID(ctx context.Context, userID utils.BinaryUUID, limit, offset int) ([]*entities.Activity, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*entities.Activity, error)
	GetByResourceType(ctx context.Context, resourceType string, limit, offset int) ([]*entities.Activity, error) 
	GetByActionType(ctx context.Context, actionType string, limit, offset int) ([]*entities.Activity, error) 
	GetFailedLoginAttempts(ctx context.Context, timeWindow time.Duration, limit int) ([]*entities.Activity, error)
	GetUserActivitySummary(ctx context.Context, userID utils.BinaryUUID, days int) (*dto.UserActivitySummaryResponse, error) 
	DeleteOldActivities(ctx context.Context, olderThan time.Duration) error
}

// UserActivityService defines the interface for interacting with user activity service logic
type UserActivityService interface {
	LogUserActivity(ctx context.Context, userID, username, actionType, resourceType string, resourceID *string, ipAddress, details, oldValue, newValue *string) error
	ListUserActivities(ctx context.Context, filter dto.ActivityFilterRequest) (*dto.UserActivityListResponse, error)
	GetUserActivitySummary(ctx context.Context, userID string, days int) (*dto.UserActivitySummaryResponse, error)
	GetSecurityAlerts(ctx context.Context, timeWindow time.Duration, limit int) ([]dto.SecurityAlertResponse, error)
	CleanOldActivities(ctx context.Context, olderThan time.Duration) error
}