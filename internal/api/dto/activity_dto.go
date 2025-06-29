package dto

import (
	"time"
)

// LogActivityRequest represents the request body for logging a new user activity.
type LogActivityRequest struct {
	Username     string  `json:"username" binding:"required"`
	ActionType   string  `json:"action_type" binding:"required"`
	ResourceType string  `json:"resource_type" binding:"required"`
	ResourceID   *string `json:"resource_id,omitempty"`
	IPAddress    *string `json:"ip_address,omitempty"`
	Details      *string `json:"details,omitempty"`
	OldValue     *string `json:"old_value,omitempty"`
	NewValue     *string `json:"new_value,omitempty"`
}

// UserActivityResponse represents the response structure for a single user activity entry.
type UserActivityResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	Username     string    `json:"username"`
	ActionType   string    `json:"actionType"`
	ResourceType string    `json:"resourceType"`
	ResourceID   *string   `json:"resourceId,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
	IPAddress    *string   `json:"ipAddress,omitempty"`
	Details      *string   `json:"details,omitempty"`
	OldValue     *string   `json:"oldValue,omitempty"`
	NewValue     *string   `json:"newValue,omitempty"`
}

// UserActivityListResponse represents a paginated list of user activities.
type UserActivityListResponse struct {
	Activities []UserActivityResponse `json:"activities"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"pageSize"`
	TotalPages int                    `json:"totalPages"`
}

// UserActivitySummaryResponse provides aggregated activity data for reporting.
type UserActivitySummaryResponse struct {
	UserID               string    `json:"userId"`
	TotalActivities      int64     `json:"totalActivities"`
	LoginCount           int64     `json:"loginCount"`
	CreateOperations     int64     `json:"createOperations"`
	UpdateOperations     int64     `json:"updateOperations"`
	DeleteOperations     int64     `json:"deleteOperations"`
	ReportGenerations    int64     `json:"reportGenerations"`
	LastActivity         time.Time `json:"lastActivity"`
	MostAccessedResource string    `json:"mostAccessedResource"`
}

// ActivityFilterRequest represents the request parameters for filtering user activities.
type ActivityFilterRequest struct {
	UserID       *string    `json:"userId,omitempty" form:"userId"`
	ActionType   *string    `json:"actionType,omitempty" form:"actionType"`
	ResourceType *string    `json:"resourceType,omitempty" form:"resourceType"` 
	ResourceID   *string    `json:"resourceId,omitempty" form:"resourceId"`    
	IPAddress    *string    `json:"ipAddress,omitempty" form:"ipAddress"`
	StartDate    *time.Time `json:"startDate,omitempty" form:"startDate" time_format:"2006-01-02T15:04:05Z07:00"`
	EndDate      *time.Time `json:"endDate,omitempty" form:"endDate" time_format:"2006-01-02T15:04:05Z07:00"`
	SearchTerm   *string    `json:"searchTerm,omitempty" form:"searchTerm"`
	Page     int `json:"page" form:"page,default=1"`
	PageSize int `json:"pageSize" form:"pageSize,default=10"`
}

// SecurityAlertResponse represents a specific type of user activity that indicates a potential security concern.
type SecurityAlertResponse struct {
	AlertID         string                `json:"alertId"`
	UserID          string                `json:"userId"`
	Username        string                `json:"username"`
	ActionType      string                `json:"actionType"`
	Timestamp       time.Time             `json:"timestamp"`
	IPAddress       *string               `json:"ipAddress,omitempty"`
	Description     string                `json:"description"`
	Severity        string                `json:"severity"`
	RelatedActivity *UserActivityResponse `json:"relatedActivity,omitempty"`
}