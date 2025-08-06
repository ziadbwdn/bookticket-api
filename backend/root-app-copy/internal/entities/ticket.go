package entities

import (
	"time"
	"root-app/internal/utils"
)

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