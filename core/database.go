package core

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global GORM database instance (singleton).
var DB *gorm.DB

// InitDB initialises the SQLite connection, applies PRAGMA settings for
// performance and foreign-key enforcement, then runs AutoMigrate for every
// registered model.
func InitDB(cfg *Config, models ...interface{}) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(cfg.DBName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Enable foreign-key constraints (off by default in SQLite).
	db.Exec("PRAGMA foreign_keys = ON")
	// Use Write-Ahead Logging for better concurrent read/write performance.
	db.Exec("PRAGMA journal_mode = WAL")

	// Auto-migrate all provided models.
	if err := db.AutoMigrate(models...); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	DB = db
	log.Println("Database initialised:", cfg.DBName)
	return db
}
