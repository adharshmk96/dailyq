package database

import (
	"log/slog"

	"gorm.io/gorm"

	"dailyq-api/internal/modules/auth"
	"dailyq-api/internal/modules/journal"
)

// models lists every entity managed by auto-migration.
func models() []any {
	return []any{
		&auth.User{},
		&auth.Session{},
		&auth.PasswordResetToken{},
		&journal.Tag{},
		&journal.Task{},
		&journal.Note{},
	}
}

// Migrate brings the schema up to date.
func Migrate(db *gorm.DB, log *slog.Logger) error {
	if err := db.AutoMigrate(models()...); err != nil {
		return err
	}
	log.Info("migrations applied", "models", len(models()))
	return nil
}
