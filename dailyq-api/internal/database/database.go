package database

import (
	"log/slog"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"dailyq-api/internal/config"
)

// Open connects to the SQLite database described by cfg.
func Open(cfg config.DatabaseConfig, log *slog.Logger) (*gorm.DB, error) {
	level := gormlogger.Silent
	if cfg.LogMode {
		level = gormlogger.Info
	}

	db, err := gorm.Open(sqlite.Open(cfg.DSN), &gorm.Config{
		Logger: gormlogger.Default.LogMode(level),
	})
	if err != nil {
		return nil, err
	}

	// SQLite needs foreign keys enabled per connection.
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return nil, err
	}

	log.Info("database connected", "dsn", cfg.DSN)
	return db, nil
}

// Close releases the underlying connection pool.
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
