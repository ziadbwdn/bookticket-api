package entities

import (
	"time"
	"root-app/internal/utils"
)

type EventStatus string

const (
	EventStatusActive   EventStatus = "active"
	EventStatusOngoing  EventStatus = "ongoing"
	EventStatusFinished EventStatus = "finished"
	EventStatusCancelled EventStatus = "cancelled"
)

type Event struct {
	ID          utils.BinaryUUID `gorm:"primaryKey;type:binary(16)" json:"id"`
	Name        string           `gorm:"size:255;not null;uniqueIndex" json:"name"`
	Description string           `gorm:"type:text" json:"description"`
	Category    string           `gorm:"size:100;not null" json:"category"`
	Venue       string           `gorm:"size:255;not null" json:"venue"`
	StartDate   time.Time        `gorm:"not null" json:"startDate"`
	EndDate     time.Time        `gorm:"not null" json:"endDate"`
	Capacity    int              `gorm:"not null;check:capacity >= 0" json:"capacity"`
	Price       *utils.GormDecimal `gorm:"type:decimal(10,2);not null" json:"price"`
	Status      EventStatus      `gorm:"size:20;not null;default:'active'" json:"status"`
	IsActive    bool             `gorm:"default:true" json:"isActive"`
	CreatedBy   utils.BinaryUUID `gorm:"type:binary(16);not null" json:"createdBy"`
	CreatedAt   time.Time        `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time        `gorm:"autoUpdateTime" json:"updatedAt"`

	// Relations
	Creator *User    `gorm:"foreignKey:CreatedBy;references:ID" json:"creator,omitempty"`
	Tickets []Ticket `gorm:"foreignKey:EventID;references:ID" json:"tickets,omitempty"`
}