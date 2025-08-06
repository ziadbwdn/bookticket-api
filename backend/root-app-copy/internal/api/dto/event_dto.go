package dto

import (
	"root-app/internal/entities"
	"root-app/internal/utils"
	"time"
)

type EventFilterRequest struct {
	Category  *string				`json:"category,omitempty" form:"category"`  
	Status    *string				`json:"status" binding:"required"`
	IPAddress *string    			`json:"ipAddress,omitempty" form:"ipAddress"` 
	StartDate *time.Time 			`json:"startDate,omitempty" form:"startDate" time_format:"2006-01-02T15:04:05Z07:00"` 
	EndDate   *time.Time 			`json:"endDate,omitempty" form:"endDate" time_format:"2006-01-02T15:04:05Z07:00"` 
	Search    *string				`json:"search,omitempty" form:"search"`
	Page      int					`json:"page" form:"page,default=1"`
	PageSize  int 					`json:"pageSize" form:"pageSize,default=10"`
}

// CreateStationRequest represents the request body for creating a new station.
type CreateEventRequest struct {
	Name        string           	`json:"name" binding:"required"`
	Description string           	`json:"description" binding:"required"`
	Category    string           	`json:"category" binding:"required"`
	Venue       string           	`json:"venue" binding:"required"`
	StartDate   time.Time        	`json:"start_date" binding:"required"`
	EndDate     time.Time        	`json:"end_date" binding:"required"`
	Capacity    int              	`json:"capacity" binding:"omitempty"`
	Price       string 				`json:"price" binding:"required"`
	Status      string		      	`json:"status" binding:"omitempty"`
}

// UpdateEventStatusRequest represents the request body for updating an existing station.
type UpdateEventRequest struct {
	Name        *string           	`json:"name" binding:"omitempty"`
	Description *string           	`json:"description" binding:"required"`
	Category    *string           	`json:"category" binding:"required"`
	Venue       *string           	`json:"venue" binding:"required"`
	StartDate   *time.Time        	`json:"start_date" binding:"required"`
	EndDate     *time.Time        	`json:"end_date" binding:"required"`
	Capacity    *int              	`json:"capacity" binding:"omitempty"`
	Price       *string				`json:"price" binding:"required"`
	Status      *string     	 	`json:"status" binding:"omitempty"`
}

// StationResponse represents the response body for station details.
type EventResponse struct {
	ID          utils.BinaryUUID 	`json:"id"`
	Name        string           	`json:"name" binding:"required"`
	Description string           	`json:"description" binding:"required"`
	Category    string           	`json:"category" binding:"required"`
	Venue       string           	`json:"venue" binding:"required"`
	StartDate   time.Time        	`json:"start_date" binding:"required"`
	EndDate     time.Time        	`json:"end_date" binding:"required"`
	Capacity    int              	`json:"capacity" binding:"omitempty"`
	Price       string 				`json:"price" binding:"required"`
	Status      string		      	`json:"status" binding:"required"`
	CreatedBy   utils.BinaryUUID 	`json:"createdBy"`
	CreatedAt   time.Time        	`json:"created_at"`
	UpdatedAt   time.Time        	`json:"updated_at,omitempty"`
}

// map event to response
func MapEventToResponse(event *entities.Event) *EventResponse {
	if event == nil {
		return nil
	}

	return &EventResponse{
		ID:				event.ID,
		Name:       	event.Name,
		Description:	event.Description,
		Category:    	event.Category,
		Venue:       	event.Venue,
		StartDate:   	event.StartDate,
		EndDate:     	event.EndDate,
		Capacity:    	event.Capacity,
		Price:       	utils.GormDecimalToString(event.Price),
		Status:      	string(event.Status),
		CreatedAt:      event.CreatedAt,
		UpdatedAt:      event.UpdatedAt,
	}
}

// ListStationsResponse (optional, for listing multiple stations)
type ListEventResponse struct {
	Event []EventResponse `json:"events"`
	Total    int          `json:"total"`
}