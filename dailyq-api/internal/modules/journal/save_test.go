package journal

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"dailyq-api/internal/modules/auth"
)

func TestSaveNoteCreatesNewRecordWithPresetID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		t.Fatalf("pragma: %v", err)
	}
	if err := db.AutoMigrate(&auth.User{}, &Tag{}, &Task{}, &Note{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	userID := uuid.NewString()
	if err := db.Create(&auth.User{ID: userID, Email: "a@b.com", Name: "A", PasswordHash: "x"}).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	tagID := uuid.NewString()
	if err := db.Create(&Tag{ID: tagID, UserID: userID, Name: "agoda", Color: "primary"}).Error; err != nil {
		t.Fatalf("create tag: %v", err)
	}

	noteID := uuid.NewString()
	date := "2026-05-20"
	repo := NewRepository(db)

	t.Run("without tags", func(t *testing.T) {
		note := &Note{ID: noteID, UserID: userID, Body: "hello", Date: &date}
		if err := repo.SaveNote(context.Background(), note); err != nil {
			t.Fatalf("save note: %v", err)
		}
		var count int64
		db.Model(&Note{}).Count(&count)
		if count != 1 {
			t.Fatalf("expected 1 note, got %d", count)
		}
	})

	t.Run("with partial tag via SaveNote", func(t *testing.T) {
		noteID3 := uuid.NewString()
		note := &Note{
			ID:     noteID3,
			UserID: userID,
			Body:   "tagged",
			Date:   &date,
			Tags:   []Tag{{ID: tagID}},
		}
		if err := repo.SaveNote(context.Background(), note); err != nil {
			t.Fatalf("save note with tags: %v", err)
		}
	})
}
