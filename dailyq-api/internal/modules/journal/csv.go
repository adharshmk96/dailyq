package journal

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const (
	csvHeader      = "type,id,date,task,note,done,tags,color"
	csvTagSep      = "|"
	maxImportRows  = 10_000
	maxImportBytes = 5 << 20 // 5 MiB
)

var (
	validCSVHeader = strings.Split(csvHeader, ",")
	validTagColors = map[string]struct{}{
		"primary": {}, "success": {}, "warning": {},
		"info": {}, "error": {}, "neutral": {},
	}
)

// csvRow is one decoded import row (1-based data row number is tracked separately).
type csvRow struct {
	Type  string
	ID    string
	Date  *string
	Task  string
	Note  string
	Done  *bool
	Tags  []string
	Color string
}

// EncodeExport serialises all journal data into the canonical CSV format.
func EncodeExport(tags []Tag, tasks []Task, notes []Note) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	if err := w.Write(validCSVHeader); err != nil {
		return nil, err
	}

	tagNames := make(map[string]string, len(tags))
	for i := range tags {
		tagNames[tags[i].ID] = tags[i].Name
		if err := w.Write(tagRow(&tags[i])); err != nil {
			return nil, err
		}
	}

	for i := range tasks {
		if err := w.Write(taskRow(&tasks[i], tagNames)); err != nil {
			return nil, err
		}
	}

	for i := range notes {
		if err := w.Write(noteRow(&notes[i], tagNames)); err != nil {
			return nil, err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func tagRow(t *Tag) []string {
	return []string{
		"tag",
		t.ID,
		"",
		t.Name,
		"",
		"",
		"",
		t.Color,
	}
}

func taskRow(t *Task, tagNames map[string]string) []string {
	done := ""
	if t.Done {
		done = "true"
	}
	return []string{
		"task",
		t.ID,
		stringOrEmpty(t.Date),
		t.Title,
		"",
		done,
		joinTagNames(t.Tags, tagNames),
		"",
	}
}

func noteRow(n *Note, tagNames map[string]string) []string {
	return []string{
		"note",
		n.ID,
		stringOrEmpty(n.Date),
		"",
		n.Body,
		"",
		joinTagNames(n.Tags, tagNames),
		"",
	}
}

func joinTagNames(tags []Tag, names map[string]string) string {
	if len(tags) == 0 {
		return ""
	}
	parts := make([]string, 0, len(tags))
	for i := range tags {
		if name, ok := names[tags[i].ID]; ok {
			parts = append(parts, name)
		}
	}
	return strings.Join(parts, csvTagSep)
}

func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// ParseImport reads and validates a CSV export. Row numbers in errors are 1-based
// and include the header as row 1.
func ParseImport(data []byte) ([]csvRow, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("csv file is empty")
	}
	if len(data) > maxImportBytes {
		return nil, fmt.Errorf("csv file exceeds %d bytes", maxImportBytes)
	}

	r := csv.NewReader(bytes.NewReader(data))
	r.FieldsPerRecord = -1
	r.LazyQuotes = true

	header, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	if !headerMatches(header) {
		return nil, fmt.Errorf("invalid header: expected %q", csvHeader)
	}

	var rows []csvRow
	rowNum := 1

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		rowNum++
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", rowNum, err)
		}
		if len(rows) >= maxImportRows {
			return nil, fmt.Errorf("csv exceeds %d data rows", maxImportRows)
		}

		row, parseErr := parseRecord(record)
		if parseErr != nil {
			return nil, fmt.Errorf("row %d: %w", rowNum, parseErr)
		}
		if isBlankRow(row) {
			continue
		}
		rows = append(rows, row)
	}

	return rows, nil
}

func headerMatches(header []string) bool {
	if len(header) != len(validCSVHeader) {
		return false
	}
	for i, col := range validCSVHeader {
		if strings.TrimSpace(strings.ToLower(header[i])) != col {
			return false
		}
	}
	return true
}

func parseRecord(record []string) (csvRow, error) {
	padded := padRecord(record, len(validCSVHeader))
	row := csvRow{
		Type:  strings.ToLower(strings.TrimSpace(padded[0])),
		ID:    strings.TrimSpace(padded[1]),
		Task:  padded[3],
		Note:  padded[4],
		Color: strings.TrimSpace(padded[7]),
	}

	date := strings.TrimSpace(padded[2])
	if date != "" {
		row.Date = &date
	}

	doneRaw := strings.TrimSpace(strings.ToLower(padded[5]))
	if doneRaw != "" {
		done, err := strconv.ParseBool(doneRaw)
		if err != nil {
			return csvRow{}, fmt.Errorf("invalid done value %q", padded[5])
		}
		row.Done = &done
	}

	tagsRaw := strings.TrimSpace(padded[6])
	if tagsRaw != "" {
		for _, part := range strings.Split(tagsRaw, csvTagSep) {
			name := strings.TrimSpace(part)
			if name != "" {
				row.Tags = append(row.Tags, name)
			}
		}
	}

	return row, nil
}

func padRecord(record []string, n int) []string {
	if len(record) >= n {
		return record
	}
	out := make([]string, n)
	copy(out, record)
	return out
}

func isBlankRow(row csvRow) bool {
	return row.Type == "" && row.ID == "" && row.Task == "" && row.Note == ""
}

func validateUUID(id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("id must be a valid UUID")
	}
	return nil
}

func validateTagColor(color string) (string, error) {
	if color == "" {
		return "neutral", nil
	}
	c := strings.ToLower(color)
	if _, ok := validTagColors[c]; !ok {
		return "", fmt.Errorf("invalid color %q", color)
	}
	return c, nil
}
