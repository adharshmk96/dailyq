package journal

import (
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"dailyq-api/internal/modules/auth"
)

func TestImportCSVWithEmptyNoteIDs(t *testing.T) {
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

	userID := "11111111-1111-1111-1111-111111111111"
	user := &auth.User{ID: userID, Email: "test@example.com", Name: "Test", PasswordHash: "hash"}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	csv := strings.Join([]string{
		"type,id,date,task,note,done,tags,color",
		"tag,22222222-2222-2222-2222-222222222222,,work,,,,primary",
		"note,33333333-3333-3333-3333-333333333333,2026-08-13,,existing note,,work,",
		"note,,2026-05-20,,new note without id,,work,",
		"note,,2026-05-21,,another new note,,work,",
	}, "\n")

	svc := NewService(NewRepository(db), slog.Default())
	result, err := svc.ImportCSV(context.Background(), userID, []byte(csv))
	if err != nil {
		t.Fatalf("import: %v", err)
	}

	if result.Skipped > 0 {
		t.Fatalf("expected no skipped rows, got %d: %+v", result.Skipped, result.Errors)
	}
	if result.Created.Tags != 1 {
		t.Fatalf("expected 1 tag created, got %d", result.Created.Tags)
	}
	if result.Created.Notes != 3 {
		t.Fatalf("expected 3 notes created, got %d", result.Created.Notes)
	}
}
