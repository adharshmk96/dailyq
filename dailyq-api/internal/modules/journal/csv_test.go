package journal

import (
	"strings"
	"testing"
)

func TestEncodeParseRoundTrip(t *testing.T) {
	date := "2026-08-14"
	tags := []Tag{
		{ID: "11111111-1111-1111-1111-111111111111", Name: "work", Color: "primary"},
		{ID: "22222222-2222-2222-2222-222222222222", Name: "home", Color: "neutral"},
	}
	tasks := []Task{
		{
			ID:    "33333333-3333-3333-3333-333333333333",
			Title: "Buy milk",
			Done:  true,
			Date:  &date,
			Tags:  []Tag{tags[0]},
		},
	}
	notes := []Note{
		{
			ID:   "44444444-4444-4444-4444-444444444444",
			Body: "Line one\nLine two",
			Tags: []Tag{tags[0], tags[1]},
		},
	}

	encoded, err := EncodeExport(tags, tasks, notes)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	rows, err := ParseImport(encoded)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(rows) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(rows))
	}

	if rows[0].Type != "tag" || rows[0].Task != "work" || rows[0].Color != "primary" {
		t.Fatalf("unexpected tag row: %+v", rows[0])
	}
	if rows[1].Type != "tag" || rows[1].Task != "home" {
		t.Fatalf("unexpected second tag row: %+v", rows[1])
	}
	if rows[2].Type != "task" || rows[2].Task != "Buy milk" || rows[2].Done == nil || !*rows[2].Done {
		t.Fatalf("unexpected task row: %+v", rows[2])
	}
	if rows[3].Note != "Line one\nLine two" || len(rows[3].Tags) != 2 {
		t.Fatalf("unexpected note row: %+v", rows[3])
	}
}

func TestParseImportRejectsBadHeader(t *testing.T) {
	_, err := ParseImport([]byte("bad,header\n"))
	if err == nil || !strings.Contains(err.Error(), "invalid header") {
		t.Fatalf("expected invalid header error, got %v", err)
	}
}

func TestParseImportMultilineNote(t *testing.T) {
	csv := `type,id,date,task,note,done,tags,color
note,44444444-4444-4444-4444-444444444444,,,"hello
world",,,
`
	rows, err := ParseImport([]byte(csv))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(rows) != 1 || rows[0].Note != "hello\nworld" {
		t.Fatalf("unexpected note: %+v", rows[0])
	}
}
