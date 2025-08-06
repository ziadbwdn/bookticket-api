package entities

import (
	"time"

	"root-app/internal/utils"
)

// Activity represents an activity log entry for audit trails
type Activity struct {
	ID           utils.BinaryUUID `gorm:"primaryKey;type:binary(16)" json:"id"`    // Unique ID for the activity log entry
	UserID       utils.BinaryUUID `gorm:"type:binary(16);not null" json:"user_id"` // ID of the user who performed the action (matches User.ID)
	Username     string           `gorm:"size:50;not null" json:"username"`        // Username at the time of the action (denormalized for convenience)
	ActionType   string           `gorm:"size:50;not null" json:"action_type"`     // Type of action performed (e.g., "login", "create_project", "update_station")
	ResourceType string           `gorm:"size:50;not null" json:"resource_type"`   // Type of resource affected (e.g., "Project", "Station")
	ResourceID   *string          `gorm:"size:36" json:"resource_id,omitempty"`    // ID of the resource affected (e.g., Project ID, Station ID), optional if action is not resource-specific (e.g., login)
	Timestamp    time.Time        `gorm:"autoCreateTime" json:"timestamp"`         // When the action occurred
	IPAddress    *string          `gorm:"size:45" json:"ip_address,omitempty"`     // IP address from which the action originated, optional
	Details      *string          `gorm:"type:text" json:"details,omitempty"`      // Additional human-readable details about the action, optional
	OldValue     *string          `gorm:"type:json" json:"old_value,omitempty"`    // Optional: JSON string of the resource state *before* update
	NewValue     *string          `gorm:"type:json" json:"new_value,omitempty"`    // Optional: JSON string of the resource state *after* update
}

// ActionType constants for common user activities.
// This list can be expanded as needed to be more granular.
const (
	ActionTypeLogin                	= "login"
	ActionTypeLogout               	= "logout"
	ActionTypeRegisterUser         	= "register_user"
	ActionTypeUpdateProfile        	= "update_profile"
	ActionTypePasswordResetRequest 	= "password_reset_request" // When a user requests a password reset link
	ActionTypePasswordReset      	= "password_reset"         // When a user successfully resets their password
	ActionTypeEmailVerified      	= "email_verified"         // When a user successfully verifies their email
	ActionTypeTokenRefresh       	= "token_refresh"          // When a user refreshes their access token
	ActionTypeCreateEvent        	= "create_event" 
	ActionTypeUpdateEvent        	= "update_event"
	ActionTypeGetAllEvents        	= "view_events"
	ActionTypeDeleteEvent        	= "delete_event"
	ActionTypePurchaseTicket		= "purchase-ticket"
	ActionTypeUpdateTicketStatus	= "update_ticket_status"
	ActionTypeViewAllTicket			= "view_tickets"
	ActionTypeDeleteTicket			= "delete_ticket"
	ActionTypeSummaryReport       	= "summary_report"
	ActionTypeTicketEventReport		= "ticket_event_report"
	ActionTypeFailedLogin          	= "failed_login_attempt"
	ActionTypeUnauthorizedAccess   	= "unauthorized_access"
	ActionTypeRoleChange           	= "role_change"
)

// ResourceType constants for affected entities.
const (
	ResourceTypeUser	= "User"
	ResourceTypeEvent 	= "Event" // change
	ResourceTypeTicket 	= "Ticket" // change
	ResourceTypeReport	= "Report" // 
)

// ActivityLogContext holds metadata for an action, gathered by the handler.
type ActivityLogContext struct {
	UserID    string
	Username  string
	IPAddress string
}

// UserActivitySummary provides aggregated activity data for reporting.
type UserActivitySummary struct {
	UserID               utils.BinaryUUID `json:"userId"`
	TotalActivities      int64            `json:"totalActivities"`
	LoginCount           int64            `json:"loginCount"`
	CreateOperations     int64            `json:"createOperations"`
	UpdateOperations     int64            `json:"updateOperations"`
	DeleteOperations     int64            `json:"deleteOperations"`
	ReportGenerations    int64            `json:"reportGenerations"`
	LastActivity         time.Time        `json:"lastActivity"`
	MostAccessedResource string           `json:"mostAccessedResource"`
}