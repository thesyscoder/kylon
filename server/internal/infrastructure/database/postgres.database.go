/**
 * @File: database.go
 * @Title: Database Connection Management
 * @Description: Handles the establishment and configuration of the PostgreSQL database connection using GORM.
 * @Author: thesyscoder (github.com/thesyscoder)
 */

package database

import (
	"fmt"
	"log"

	// Import time for connection lifetime settings

	"github.com/thesyscoder/kylon/internal/infrastructure/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectPostgres establishes a connection to the PostgreSQL database using GORM.
// It applies connection pool settings from the provided configuration.
func ConnectPostgres(cfg *config.Config) (*gorm.DB, error) {
	// Construct the Data Source Name (DSN) for PostgreSQL connection.
	// TimeZone is set to Asia/Kolkata as per regional context.
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Kolkata",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.SslMode,
	)

	// Open a GORM database connection.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Set GORM logger to info level. In production, consider adjusting based on log level.
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		// Corrected error wrapping for better debugging.
		return nil, fmt.Errorf("[Postgres]: Failed to connect to database: %w", err)
	}

	// Get the underlying sql.DB object to configure connection pool settings.
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("[Database]: Failed to get generic database object: %w", err)
	}

	// Apply connection pool settings from configuration for optimal performance.
	sqlDB.SetMaxOpenConns(cfg.Database.MaxConnections)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnectionMaxLifetime)

	// Ping the database to verify the connection is alive.
	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("[Database]: Failed to ping database: %w", err)
	}

	log.Println("[Database]: Successfully connected to PostgreSQL.")
	return db, nil
}
