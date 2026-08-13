package database

import (
	"errors"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"dailyq-api/internal/modules/auth"
	"dailyq-api/internal/modules/journal"

	"github.com/google/uuid"
)

// DefaultSeedUser is the development account created by Seed.
const (
	DefaultSeedEmail    = "demo@dailyq.local"
	DefaultSeedPassword = "password123"
)

// Seed inserts development data. It is idempotent.
func Seed(db *gorm.DB, log *slog.Logger) error {
	user, err := seedUser(db, log)
	if err != nil {
		return err
	}
	return seedJournal(db, log, user.ID)
}

func seedUser(db *gorm.DB, log *slog.Logger) (*auth.User, error) {
	var existing auth.User
	err := db.First(&existing, "email = ?", DefaultSeedEmail).Error
	if err == nil {
		log.Info("seed user already exists", "email", DefaultSeedEmail)
		return &existing, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(DefaultSeedPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &auth.User{
		ID:           uuid.NewString(),
		Email:        DefaultSeedEmail,
		Name:         "Demo User",
		PasswordHash: string(hash),
	}
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	log.Info("seed user created", "email", user.Email, "password", DefaultSeedPassword)
	return user, nil
}

// seedJournal gives the demo user something to look at on the dashboard. It
// only runs when the user has no tasks yet.
func seedJournal(db *gorm.DB, log *slog.Logger, userID string) error {
	var taskCount int64
	if err := db.Model(&journal.Task{}).Where("user_id = ?", userID).Count(&taskCount).Error; err != nil {
		return err
	}
	if taskCount > 0 {
		log.Info("seed journal skipped, tasks already exist", "user_id", userID)
		return nil
	}

	tags := map[string]*journal.Tag{
		"work":     {ID: uuid.NewString(), UserID: userID, Name: "Work", Color: "primary"},
		"personal": {ID: uuid.NewString(), UserID: userID, Name: "Personal", Color: "success"},
		"ideas":    {ID: uuid.NewString(), UserID: userID, Name: "Ideas", Color: "warning"},
		"health":   {ID: uuid.NewString(), UserID: userID, Name: "Health", Color: "info"},
	}

	today := isoDay(0)
	tomorrow := isoDay(1)
	yesterday := isoDay(-1)

	tasks := []struct {
		title string
		done  bool
		date  *string
		tags  []*journal.Tag
	}{
		{"Review weekly goals", false, nil, []*journal.Tag{tags["work"]}},
		{"Read design notes", true, nil, []*journal.Tag{tags["ideas"]}},
		{"Clear open general items", false, nil, nil},
		{"Morning stretch", true, &today, []*journal.Tag{tags["health"]}},
		{"Write journal entry", false, &today, []*journal.Tag{tags["personal"]}},
		{"Plan tomorrow focus", false, &today, []*journal.Tag{tags["work"]}},
		{"Team standup prep", false, &tomorrow, []*journal.Tag{tags["work"]}},
		{"Ship dashboard UI", true, &yesterday, []*journal.Tag{tags["work"]}},
	}

	notes := []struct {
		body string
		date *string
		tags []*journal.Tag
	}{
		{"Ideas for the habit tracker — keep it simple.", nil, []*journal.Tag{tags["ideas"]}},
		{"Felt focused after the morning walk.", &today, []*journal.Tag{tags["health"]}},
		{"Sketch calendar chips for days with items.", &today, []*journal.Tag{tags["ideas"]}},
		{"Tomorrow: try a shorter standup agenda.", &tomorrow, []*journal.Tag{tags["work"]}},
		{"Reflection: shipping UI before API was the right call.", &yesterday, []*journal.Tag{tags["work"]}},
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, tag := range tags {
			if err := tx.Create(tag).Error; err != nil {
				return err
			}
		}
		for _, t := range tasks {
			task := &journal.Task{
				ID:     uuid.NewString(),
				UserID: userID,
				Title:  t.title,
				Done:   t.done,
				Date:   t.date,
				Tags:   derefTags(t.tags),
			}
			if err := tx.Create(task).Error; err != nil {
				return err
			}
		}
		for _, n := range notes {
			note := &journal.Note{
				ID:     uuid.NewString(),
				UserID: userID,
				Body:   n.body,
				Date:   n.date,
				Tags:   derefTags(n.tags),
			}
			if err := tx.Create(note).Error; err != nil {
				return err
			}
		}

		log.Info("seed journal created", "user_id", userID, "tasks", len(tasks), "notes", len(notes))
		return nil
	})
}

func derefTags(tags []*journal.Tag) []journal.Tag {
	out := make([]journal.Tag, 0, len(tags))
	for _, t := range tags {
		out = append(out, *t)
	}
	return out
}

func isoDay(offset int) string {
	return time.Now().AddDate(0, 0, offset).Format("2006-01-02")
}
