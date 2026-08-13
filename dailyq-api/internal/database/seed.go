package database

import (
	"errors"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"dailyq-api/internal/modules/auth"

	"github.com/google/uuid"
)

// DefaultSeedUser is the development account created by Seed.
const (
	DefaultSeedEmail    = "demo@dailyq.local"
	DefaultSeedPassword = "password123"
)

// Seed inserts development data. It is idempotent.
func Seed(db *gorm.DB, log *slog.Logger) error {
	var existing auth.User
	err := db.First(&existing, "email = ?", DefaultSeedEmail).Error
	if err == nil {
		log.Info("seed skipped, user already exists", "email", DefaultSeedEmail)
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(DefaultSeedPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &auth.User{
		ID:           uuid.NewString(),
		Email:        DefaultSeedEmail,
		Name:         "Demo User",
		PasswordHash: string(hash),
	}
	if err := db.Create(user).Error; err != nil {
		return err
	}

	log.Info("seed user created", "email", user.Email, "password", DefaultSeedPassword)
	return nil
}
