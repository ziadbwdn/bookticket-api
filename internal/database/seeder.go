package database

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"root-app/internal/entities"
	"root-app/internal/utils"
)

// RunSeeder populates the database with initial data for development and testing
func RunSeeder(db *gorm.DB) error {
	// Check if data already exists to prevent duplicate entries
	var userCount int64
	db.Model(&entities.User{}).Count(&userCount)
	if userCount > 0 {
		return nil // Data already seeded
	}

	// --- Seed Users ---
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	adminUser := entities.User{
		ID:           utils.NewUUID(),
		Username:     "admin",
		Email:        "admin@example.com",
		PasswordHash: string(hashedPassword),
		FullName:     "Admin User",
		Role:         entities.RoleAdmin,
		IsActive:     true,
		EmailVerified: true,
	}
	db.Create(&adminUser)

	regularUser := entities.User{
		ID:           utils.NewUUID(),
		Username:     "user",
		Email:        "user@example.com",
		PasswordHash: string(hashedPassword),
		FullName:     "Regular User",
		Role:         entities.RoleUser,
		IsActive:     true,
		EmailVerified: true,
	}
	db.Create(&regularUser)

	// --- Seed Events ---
	event1 := entities.Event{
		ID:          utils.NewUUID(),
		Name:        "Tech Conference 2025",
		Description: "An annual conference for tech enthusiasts.",
		Category:    "Technology",
		Venue:       "Convention Center",
		StartDate:   time.Now().AddDate(0, 1, 0),
		EndDate:     time.Now().AddDate(0, 1, 2),
		Capacity:    500,
		Price:       utils.Decimal(199.99),
		Status:      entities.EventStatusActive,
		CreatedBy:   adminUser.ID,
	}
	db.Create(&event1)

	event2 := entities.Event{
		ID:          utils.NewUUID(),
		Name:        "Music Festival",
		Description: "A weekend of live music from various artists.",
		Category:    "Music",
		Venue:       "Open Air Field",
		StartDate:   time.Now().AddDate(0, 2, 15),
		EndDate:     time.Now().AddDate(0, 2, 17),
		Capacity:    2000,
		Price:       utils.Decimal(89.50),
		Status:      entities.EventStatusActive,
		CreatedBy:   adminUser.ID,
	}
	db.Create(&event2)

	// --- Seed Tickets ---
	ticket1 := entities.Ticket{
		ID:         utils.NewUUID(),
		EventID:    event1.ID,
		UserID:     regularUser.ID,
		TicketCode: "TCKT-TECH-001",
		Quantity:   1,
		UnitPrice:  event1.Price,
		TotalPrice: utils.Decimal(199.99),
		Status:     entities.TicketStatusActive,
	}
	db.Create(&ticket1)

	ticket2 := entities.Ticket{
		ID:         utils.NewUUID(),
		EventID:    event2.ID,
		UserID:     regularUser.ID,
		TicketCode: "TCKT-MUSIC-001",
		Quantity:   2,
		UnitPrice:  event2.Price,
		TotalPrice: utils.Decimal(179.00),
		Status:     entities.TicketStatusActive,
	}
	db.Create(&ticket2)

	return nil
}
