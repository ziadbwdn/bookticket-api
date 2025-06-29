package database

import (
	"root-app/internal/config" 
	"root-app/internal/exception"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger" // Import GORM logger
)

// InitDB initializes the GORM database connection.
func InitDB(cfg *config.Config) (*gorm.DB, *exception.AppError) {
	// Construct the DSN (Data Source Name) for MySQL using values from the Config struct.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Open the GORM database connection. Set logger.Info to get detailed SQL logs and errors from GORM.
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), 
	})
	if err != nil {
		log.Printf("ERROR: GORM failed to open database connection: %v", err)
		return nil, exception.NewDatabaseError("DB connection failed", err)
	}
	log.Println("DEBUG: GORM database connection opened.")

	// Get the underlying sql.DB instance from GORM for pinging.
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("ERROR: Failed to get underlying SQL DB instance: %v", err)
		return nil, exception.NewDatabaseError("DB instance retrieval failed", err)
	}
	log.Println("DEBUG: Underlying SQL DB instance retrieved.")

	// Ping the database to verify the connection is alive.
	if err := sqlDB.Ping(); err != nil {
		log.Printf("ERROR: Database ping failed: %v", err) // Log the underlying error
		return nil, exception.NewDatabaseError("DB ping failed", err)
	}
	log.Println("DEBUG: Database ping successful.")

	// Return the GORM database instance and nil for the AppError on success.
	return db, nil
}

