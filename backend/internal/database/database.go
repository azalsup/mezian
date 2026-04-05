// Package database initializes the SQLite connection via GORM and applies migrations.
package database

import (
	"fmt"
	"mezian/internal/config"
	"mezian/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect opens the SQLite database and performs schema auto-migration.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	logLevel := logger.Silent
	if cfg.Server.Mode == "debug" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("ouverture SQLite: %w", err)
	}

	// Enable foreign keys (SQLite disables them by default)
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("récupération sql.DB: %w", err)
	}
	if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("PRAGMA foreign_keys: %w", err)
	}
	if _, err := sqlDB.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return nil, fmt.Errorf("PRAGMA journal_mode: %w", err)
	}
	if _, err := sqlDB.Exec("PRAGMA busy_timeout = 5000"); err != nil {
		return nil, fmt.Errorf("PRAGMA busy_timeout: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migration: %w", err)
	}
	return db, nil
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.OTPCode{},
		&models.RefreshToken{},
		&models.Category{},
		&models.AttributeDefinition{},
		&models.Shop{},
		&models.Ad{},
		&models.AdAttribute{},
		&models.Media{},
	)
}
