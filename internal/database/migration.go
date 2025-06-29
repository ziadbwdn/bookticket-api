package database

import (
	"root-app/internal/entities"

	"gorm.io/gorm"
)

// RunMigrations executes the GORM auto-migration for all application entities, centralize migration logic
func RunMigrations(db *gorm.DB) error {
	// GORM's AutoMigrate will create tables, add missing columns, and indices.
	return db.AutoMigrate(
		&entities.User{},
		&entities.Event{},
		&entities.Ticket{},
		&entities.Activity{},
	)
}