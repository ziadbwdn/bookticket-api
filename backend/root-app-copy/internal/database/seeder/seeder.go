
package seeder

import (
	"fmt"
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
		ID:           utils.NewBinaryUUID(),
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
		ID:           utils.NewBinaryUUID(),
		Username:     "user",
		Email:        "user@example.com",
		PasswordHash: string(hashedPassword),
		FullName:     "Regular User",
		Role:         entities.RoleUser,
		IsActive:     true,
		EmailVerified: true,
	}
	db.Create(&regularUser)

	// --- Seed additional 8 users ---
	for i := 1; i <= 8; i++ {
		user := entities.User{
			ID:           utils.NewBinaryUUID(),
			Username:     fmt.Sprintf("user%d", i+1),
			Email:        fmt.Sprintf("user%d@example.com", i+1),
			PasswordHash: string(hashedPassword),
			FullName:     fmt.Sprintf("User %d", i+1),
			Role:         entities.RoleUser,
			IsActive:     true,
			EmailVerified: true,
		}
		db.Create(&user)
	}

	price1, _ := utils.StringToGormDecimal("199.99")
	price2, _ := utils.StringToGormDecimal("89.50")
	totalPrice1, _ := utils.StringToGormDecimal("199.99")
	totalPrice2, _ := utils.StringToGormDecimal("179.00")

	// --- Seed Events ---
	event1 := entities.Event{
		ID:          utils.NewBinaryUUID(),
		Name:        "Tech Conference 2025",
		Description: "An annual conference for tech enthusiasts.",
		Category:    "Technology",
		Venue:       "Convention Center",
		StartDate:   time.Now().AddDate(0, 1, 0),
		EndDate:     time.Now().AddDate(0, 1, 2),
		Capacity:    500,
		Price:       price1,
		Status:      entities.EventStatusActive,
		CreatedBy:   adminUser.ID,
	}
	db.Create(&event1)

	event2 := entities.Event{
		ID:          utils.NewBinaryUUID(),
		Name:        "Music Festival",
		Description: "A weekend of live music from various artists.",
		Category:    "Music",
		Venue:       "Open Air Field",
		StartDate:   time.Now().AddDate(0, 2, 15),
		EndDate:     time.Now().AddDate(0, 2, 17),
		Capacity:    2000,
		Price:       price2,
		Status:      entities.EventStatusActive,
		CreatedBy:   adminUser.ID,
	}
	db.Create(&event2)

	// --- Seed additional 8 events ---
	for i := 1; i <= 8; i++ {
		price, _ := utils.StringToGormDecimal(fmt.Sprintf("%.2f", 100.00+float64(i)*10))
		event := entities.Event{
			ID:          utils.NewBinaryUUID(),
			Name:        fmt.Sprintf("Event %d", i+2),
			Description: fmt.Sprintf("Description for event %d", i+2),
			Category:    "Category",
			Venue:       fmt.Sprintf("Venue %d", i+2),
			StartDate:   time.Now().AddDate(0, int(time.Month(i+2)), 1),
			EndDate:     time.Now().AddDate(0, int(time.Month(i+2)), 3),
			Capacity:    100 + i*50,
			Price:       price,
			Status:      entities.EventStatusActive,
			CreatedBy:   adminUser.ID,
		}
		db.Create(&event)
	}

	// --- Seed Tickets ---
	ticket1 := entities.Ticket{
		ID:         utils.NewBinaryUUID(),
		EventID:    event1.ID,
		UserID:     regularUser.ID,
		TicketCode: "TCKT-TECH-001",
		Quantity:   1,
		UnitPrice:  event1.Price,
		TotalPrice: totalPrice1,
		Status:     entities.TicketStatusActive,
	}
	db.Create(&ticket1)

	ticket2 := entities.Ticket{
		ID:         utils.NewBinaryUUID(),
		EventID:    event2.ID,
		UserID:     regularUser.ID,
		TicketCode: "TCKT-MUSIC-001",
		Quantity:   2,
		UnitPrice:  event2.Price,
		TotalPrice: totalPrice2,
		Status:     entities.TicketStatusActive,
	}
	db.Create(&ticket2)

	// --- Seed additional 8 tickets ---
	var events []entities.Event
	db.Find(&events)

	var users []entities.User
	db.Where("role = ?", entities.RoleUser).Find(&users)

	for i := 1; i <= 8; i++ {
		totalPrice, _ := utils.StringToGormDecimal(fmt.Sprintf("%.2f", 100.00+float64(i)*20))
		ticket := entities.Ticket{
			ID:         utils.NewBinaryUUID(),
			EventID:    events[i+1].ID,
			UserID:     users[i].ID,
			TicketCode: fmt.Sprintf("TCKT-00%d", i+2),
			Quantity:   1,
			UnitPrice:  events[i+1].Price,
			TotalPrice: totalPrice,
			Status:     entities.TicketStatusActive,
		}
		db.Create(&ticket)
	}

	return nil
}
