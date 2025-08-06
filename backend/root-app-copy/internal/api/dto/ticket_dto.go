package dto

import (
	"root-app/internal/entities"
	"root-app/internal/utils"
	"time"
)

/*
Ticket Entities:

type TicketStatus string

const (
	TicketStatusActive    TicketStatus = "active"
	TicketStatusUsed      TicketStatus = "used"
	TicketStatusCancelled TicketStatus = "cancelled"
	TicketStatusExpired   TicketStatus = "expired"
)

type Ticket struct {
	ID            utils.BinaryUUID `gorm:"primaryKey;type:binary(16)" json:"id"`
	EventID       utils.BinaryUUID `gorm:"type:binary(16);not null" json:"eventId"`
	UserID        utils.BinaryUUID `gorm:"type:binary(16);not null" json:"userId"`
	TicketCode    string           `gorm:"size:50;not null;uniqueIndex" json:"ticketCode"`
	Quantity      int              `gorm:"not null;check:quantity > 0" json:"quantity"`
	UnitPrice     *utils.GormDecimal `gorm:"type:decimal(10,2);not null" json:"unitPrice"`
	TotalPrice    *utils.GormDecimal `gorm:"type:decimal(10,2);not null" json:"totalPrice"`
	Status        TicketStatus     `gorm:"size:20;not null;default:'active'" json:"status"`
	PurchaseDate  time.Time        `gorm:"autoCreateTime" json:"purchaseDate"`
	CancelledAt   *time.Time       `json:"cancelledAt,omitempty"`
	CancelReason  string           `gorm:"size:500" json:"cancelReason,omitempty"`
	CreatedAt     time.Time        `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time        `gorm:"autoUpdateTime" json:"updatedAt"`

	// Relations
	Event *Event `gorm:"foreignKey:EventID;references:ID" json:"event,omitempty"`
	User  *User  `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}
*/

// CreateStationRequest represents the request body for creating a new station.
type BookTicketRequest struct {
	EventID        	string           	`json:"event" binding:"required"`
	UserID 			utils.BinaryUUID    `json:"user_id" binding:"required"`
	TicketCode    	string           	`json:"ticket_code" binding:"required"`
	Quantity       	string           	`json:"quantity" binding:"required"`
	UnitPrice   	string	        	`json:"unit_price" binding:"required"`
	TotalPrice     	string	        	`json:"total_price" binding:"required"`
	Status    		int              	`json:"status" binding:"omitempty"`
	PurchaseDate	*time.Time 			`json:"purchase_date" binding:"required"`
}
/*
type TicketFilter struct {
	EventID utils.BinaryUUID
	UserID  utils.BinaryUUID
	Status  string
	Page    int
	Limit   int
}
*/

type TicketFilterRequest struct {
	EventID		*string				`json:"category,omitempty" form:"category"`  
	UserID    	*string				`json:"status" binding:"required"`
	IPAddress 	*string    			`json:"ipAddress,omitempty" form:"ipAddress"` 
	Status    	*string				`json:"status,omitempty" form:"status"`
	Page      	int					`json:"page" form:"page,default=1"`
	PageSize  	int 				`json:"pageSize" form:"pageSize,default=10"`
}

// UpdateEventStatusRequest represents the request body for updating an existing station.
type UpdateTicketStatusRequest struct {
	EventID        	*string           	`json:"event_id" binding:"required"`
	UserID 			*string           	`json:"user_id" binding:"required"`
	TicketCode    	*string           	`json:"ticket_code" binding:"required"`
	Status    		*int              	`json:"status" binding:"omitempty"`
	PurchaseDate	*string 			`json:"purchase_date" binding:"required"`
}

// BookTicketResponse represents the response body for station details.
type TicketResponse struct {
	ID          	utils.BinaryUUID 	`json:"id"`
	EventID        	utils.BinaryUUID           	`json:"event_id"`
	UserID 			utils.BinaryUUID    `json:"user_id"`
	TicketCode    	string           	`json:"ticket_code"`
	Quantity       	int		           	`json:"quantity"`
	UnitPrice   	string	        	`json:"unit_price"`
	TotalPrice     	string     		   	`json:"total_price"`
	Status    		string              `json:"status" default: "order booked, waiting for payment"`
	PurchaseDate  	time.Time        	`json:"purchase_date"`
	CancelledAt   	*time.Time       	`json:"cancelled_at" default: "no cancellation yet"`
	CancelReason  	string           	`json:"cancel_reason, default:"no cancellation reason yet"`
	CreatedAt     	time.Time        	`json:"created_at"`
	UpdatedAt     	time.Time        	`json:"updated_at"`
}

// UpdateTicketStatusResponse check
type UpdateTicketStatusResponse struct {
	Event        	*string           	`json:"event_id""`
	UserID 			*string           	`json:"user_id"`
	TicketCode    	*string           	`json:"ticket_code"`
	Status    		*int              	`json:"status"`
	PurchaseDate	*string 			`json:"purchase_date"`
	CancelledAt   	*time.Time       	`json:"cancelled_at"`
	CancelReason  	*string           	`json:"cancel_reason"`
}


// Map Ticket to response
func MapTicketToResponse(ticket *entities.Ticket) *TicketResponse {
	if ticket == nil {
		return nil
	}

	return &TicketResponse{
		ID:				ticket.ID,
		EventID:       	ticket.EventID,
		UserID:			ticket.UserID,
		TicketCode:    	ticket.TicketCode,
		Quantity:       ticket.Quantity,
		UnitPrice:     	utils.GormDecimalToString(ticket.UnitPrice),
		TotalPrice:    	utils.GormDecimalToString(ticket.TotalPrice),
		Status:       	string(ticket.Status),
		PurchaseDate:   ticket.PurchaseDate,
		CancelledAt:   	ticket.CancelledAt,
		CancelReason:  	ticket.CancelReason,
		CreatedAt:      ticket.CreatedAt,
		UpdatedAt:      ticket.UpdatedAt,
	}
}

// ListStationsResponse (optional, for listing multiple stations)
type ListTicketResponses struct {
	Ticket []TicketResponse `json:"tickets"`
	Total    int          `json:"total"`
}