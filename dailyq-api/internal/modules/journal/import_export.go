package journal

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// ExportCSV returns all journal data for the user as CSV bytes.
func (s *Service) ExportCSV(ctx context.Context, userID string) ([]byte, error) {
	tags, err := s.repo.ListTags(ctx, userID)
	if err != nil {
		return nil, err
	}
	tasks, err := s.repo.ListAllTasks(ctx, userID)
	if err != nil {
		return nil, err
	}
	notes, err := s.repo.ListAllNotes(ctx, userID)
	if err != nil {
		return nil, err
	}
	return EncodeExport(tags, tasks, notes)
}

// ImportCSV upserts tags, tasks, and notes from a CSV export. Rows with the same
// id as an existing entity are updated instead of duplicated.
func (s *Service) ImportCSV(ctx context.Context, userID string, data []byte) (*ImportResult, error) {
	rows, err := ParseImport(data)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCSV, err.Error())
	}

	result := &ImportResult{Errors: []ImportRowError{}}
	tagByName := make(map[string]string)

	existingTags, err := s.repo.ListTags(ctx, userID)
	if err != nil {
		return nil, err
	}
	for i := range existingTags {
		tagByName[strings.ToLower(existingTags[i].Name)] = existingTags[i].ID
	}

	// Tag rows first so names and colors are defined before items reference them.
	for i, row := range rows {
		if row.Type != "tag" {
			continue
		}
		rowNum := i + 2 // header is row 1
		if err := s.importTagRow(ctx, userID, row, tagByName, result); err != nil {
			result.Errors = append(result.Errors, ImportRowError{Row: rowNum, Message: err.Error()})
			result.Skipped++
		}
	}

	for i, row := range rows {
		if row.Type == "tag" {
			continue
		}
		rowNum := i + 2
		var err error
		switch row.Type {
		case "task":
			err = s.importTaskRow(ctx, userID, row, tagByName, result)
		case "note":
			err = s.importNoteRow(ctx, userID, row, tagByName, result)
		default:
			err = fmt.Errorf("unknown type %q", row.Type)
		}
		if err != nil {
			result.Errors = append(result.Errors, ImportRowError{Row: rowNum, Message: err.Error()})
			result.Skipped++
		}
	}

	return result, nil
}

func (s *Service) importTagRow(ctx context.Context, userID string, row csvRow, tagByName map[string]string, result *ImportResult) error {
	if err := validateUUID(row.ID); err != nil {
		return err
	}
	name := strings.TrimSpace(row.Task)
	if name == "" {
		return fmt.Errorf("tag name is required in task column")
	}
	color, err := validateTagColor(row.Color)
	if err != nil {
		return err
	}
	if row.Note != "" || len(row.Tags) > 0 {
		return fmt.Errorf("tag rows must not include note or tags columns")
	}

	_, err = s.repo.GetTag(ctx, userID, row.ID)
	exists := err == nil
	if err != nil && err != ErrTagNotFound {
		return err
	}

	taken, err := s.repo.TagNameTaken(ctx, userID, name, row.ID)
	if err != nil {
		return err
	}
	if taken {
		return fmt.Errorf("tag name %q is already taken", name)
	}

	tag := &Tag{ID: row.ID, UserID: userID, Name: name, Color: color}
	if err := s.repo.SaveTag(ctx, tag); err != nil {
		return err
	}

	tagByName[strings.ToLower(name)] = row.ID
	if exists {
		result.Updated.Tags++
	} else {
		result.Created.Tags++
	}
	return nil
}

func (s *Service) importTaskRow(ctx context.Context, userID string, row csvRow, tagByName map[string]string, result *ImportResult) error {
	if err := validateUUID(row.ID); err != nil {
		return err
	}
	title := strings.TrimSpace(row.Task)
	if title == "" {
		return fmt.Errorf("task title is required")
	}
	if strings.TrimSpace(row.Note) != "" {
		return fmt.Errorf("task rows must not include note content")
	}
	if err := validateDate(row.Date); err != nil {
		return err
	}

	done := false
	if row.Done != nil {
		done = *row.Done
	}

	tags, err := s.resolveTagNames(ctx, userID, row.Tags, tagByName, result)
	if err != nil {
		return err
	}

	_, err = s.repo.GetTask(ctx, userID, row.ID)
	exists := err == nil
	if err != nil && err != ErrTaskNotFound {
		return err
	}

	task := &Task{
		ID:     row.ID,
		UserID: userID,
		Title:  title,
		Done:   done,
		Date:   row.Date,
		Tags:   tags,
	}
	if err := s.repo.SaveTask(ctx, task); err != nil {
		return err
	}

	if exists {
		result.Updated.Tasks++
	} else {
		result.Created.Tasks++
	}
	return nil
}

func (s *Service) importNoteRow(ctx context.Context, userID string, row csvRow, tagByName map[string]string, result *ImportResult) error {
	if err := validateUUID(row.ID); err != nil {
		return err
	}
	body := strings.TrimSpace(row.Note)
	if body == "" {
		return fmt.Errorf("note body is required")
	}
	if strings.TrimSpace(row.Task) != "" {
		return fmt.Errorf("note rows must not include task title")
	}
	if row.Done != nil {
		return fmt.Errorf("note rows must not include done value")
	}
	if err := validateDate(row.Date); err != nil {
		return err
	}

	tags, err := s.resolveTagNames(ctx, userID, row.Tags, tagByName, result)
	if err != nil {
		return err
	}

	_, err = s.repo.GetNote(ctx, userID, row.ID)
	exists := err == nil
	if err != nil && err != ErrNoteNotFound {
		return err
	}

	note := &Note{
		ID:     row.ID,
		UserID: userID,
		Body:   body,
		Date:   row.Date,
		Tags:   tags,
	}
	if err := s.repo.SaveNote(ctx, note); err != nil {
		return err
	}

	if exists {
		result.Updated.Notes++
	} else {
		result.Created.Notes++
	}
	return nil
}

// resolveTagNames maps tag names to Tag records, auto-creating missing tags.
func (s *Service) resolveTagNames(ctx context.Context, userID string, names []string, tagByName map[string]string, result *ImportResult) ([]Tag, error) {
	if len(names) == 0 {
		return nil, nil
	}

	seen := make(map[string]struct{}, len(names))
	tags := make([]Tag, 0, len(names))

	for _, name := range names {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		if id, ok := tagByName[key]; ok {
			tag, err := s.repo.GetTag(ctx, userID, id)
			if err != nil {
				return nil, err
			}
			tags = append(tags, *tag)
			continue
		}

		existing, err := s.repo.GetTagByName(ctx, userID, trimmed)
		if err == nil {
			tagByName[key] = existing.ID
			tags = append(tags, *existing)
			continue
		}
		if err != ErrTagNotFound {
			return nil, err
		}

		tag := &Tag{
			ID:     uuid.NewString(),
			UserID: userID,
			Name:   trimmed,
			Color:  "neutral",
		}
		if err := s.repo.CreateTag(ctx, tag); err != nil {
			return nil, err
		}
		tagByName[key] = tag.ID
		result.Created.Tags++
		tags = append(tags, *tag)
	}

	return tags, nil
}
