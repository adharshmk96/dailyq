package journal

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
)

// recentNotesLimit caps the "Recent notes" list on the overview page.
const recentNotesLimit = 5

// isoLayout is the wire format for every date the journal stores.
const isoLayout = "2006-01-02"

// Service holds the journal business logic.
type Service struct {
	repo   Repository
	logger *slog.Logger
	now    func() time.Time
}

// NewService builds the journal service.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger, now: time.Now}
}

// Overview assembles everything the dashboard overview page renders.
func (s *Service) Overview(ctx context.Context, userID string) (*OverviewResponse, error) {
	now := s.now()
	today := isoDate(now)
	// "This week" is a rolling 7-day window ending today.
	weekStart := isoDate(now.AddDate(0, 0, -6))

	totalTasks, err := s.repo.CountTasks(ctx, userID)
	if err != nil {
		return nil, err
	}
	completedToday, err := s.repo.CountTasksDoneOn(ctx, userID, today)
	if err != nil {
		return nil, err
	}
	openGeneral, err := s.repo.CountOpenGeneralTasks(ctx, userID)
	if err != nil {
		return nil, err
	}
	notesThisWeek, err := s.repo.CountNotesSince(ctx, userID, weekStart)
	if err != nil {
		return nil, err
	}

	tags, err := s.repo.ListTags(ctx, userID)
	if err != nil {
		return nil, err
	}
	todayTasks, err := s.repo.ListTasksForDate(ctx, userID, today)
	if err != nil {
		return nil, err
	}
	recentNotes, err := s.repo.ListRecentNotes(ctx, userID, recentNotesLimit)
	if err != nil {
		return nil, err
	}

	res := &OverviewResponse{
		Today: today,
		Stats: OverviewStats{
			TotalTasks:     totalTasks,
			CompletedToday: completedToday,
			OpenGeneral:    openGeneral,
			NotesThisWeek:  notesThisWeek,
		},
		Tags:        toTagResponses(tags),
		TodayTasks:  toTaskResponses(todayTasks),
		RecentNotes: toNoteResponses(recentNotes),
	}

	return res, nil
}

// ---- tags ----

// ListTags returns the user's tags, alphabetically.
func (s *Service) ListTags(ctx context.Context, userID string) ([]TagResponse, error) {
	tags, err := s.repo.ListTags(ctx, userID)
	if err != nil {
		return nil, err
	}
	return toTagResponses(tags), nil
}

// CreateTag adds a tag, rejecting a name the user already uses.
func (s *Service) CreateTag(ctx context.Context, userID string, req CreateTagRequest) (*TagResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrInvalidRequest
	}

	taken, err := s.repo.TagNameTaken(ctx, userID, name, "")
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, ErrTagNameTaken
	}

	tag := &Tag{
		ID:     uuid.NewString(),
		UserID: userID,
		Name:   name,
		Color:  req.Color,
	}
	if err := s.repo.CreateTag(ctx, tag); err != nil {
		return nil, err
	}

	res := NewTagResponse(tag)
	return &res, nil
}

// UpdateTag renames or recolours one of the user's tags.
func (s *Service) UpdateTag(ctx context.Context, userID, tagID string, req UpdateTagRequest) (*TagResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrInvalidRequest
	}

	taken, err := s.repo.TagNameTaken(ctx, userID, name, tagID)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, ErrTagNameTaken
	}

	if err := s.repo.UpdateTag(ctx, userID, tagID, name, req.Color); err != nil {
		return nil, err
	}

	tag, err := s.repo.GetTag(ctx, userID, tagID)
	if err != nil {
		return nil, err
	}

	res := NewTagResponse(tag)
	return &res, nil
}

// DeleteTag removes a tag and detaches it from every item that carried it.
func (s *Service) DeleteTag(ctx context.Context, userID, tagID string) error {
	return s.repo.DeleteTag(ctx, userID, tagID)
}

// ---- entries ----

// Entries returns the tasks and notes for one date, or for the undated
// "general" bucket when date is nil.
func (s *Service) Entries(ctx context.Context, userID string, date *string) (*EntriesResponse, error) {
	if err := validateDate(date); err != nil {
		return nil, err
	}

	tasks, err := s.repo.ListTasks(ctx, userID, date)
	if err != nil {
		return nil, err
	}
	notes, err := s.repo.ListNotes(ctx, userID, date)
	if err != nil {
		return nil, err
	}

	return &EntriesResponse{
		Date:  date,
		Tasks: toTaskResponses(tasks),
		Notes: toNoteResponses(notes),
	}, nil
}

// Dates lists every day the user has at least one item on.
func (s *Service) Dates(ctx context.Context, userID string) (*DatesResponse, error) {
	dates, err := s.repo.ListDatesWithItems(ctx, userID)
	if err != nil {
		return nil, err
	}
	if dates == nil {
		dates = []string{}
	}
	return &DatesResponse{Dates: dates}, nil
}

// ---- tasks ----

// CreateTask adds a task to a date or to the general bucket.
func (s *Service) CreateTask(ctx context.Context, userID string, req CreateTaskRequest) (*TaskResponse, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, ErrInvalidRequest
	}
	if err := validateDate(req.Date); err != nil {
		return nil, err
	}

	tags, err := s.resolveTags(ctx, userID, req.TagIDs)
	if err != nil {
		return nil, err
	}

	task := &Task{
		ID:     uuid.NewString(),
		UserID: userID,
		Title:  title,
		Date:   req.Date,
		Tags:   tags,
	}
	if err := s.repo.CreateTask(ctx, task); err != nil {
		return nil, err
	}

	res := NewTaskResponse(task)
	return &res, nil
}

// UpdateTask rewrites a task's title, date and tag set.
func (s *Service) UpdateTask(ctx context.Context, userID, taskID string, req UpdateTaskRequest) (*TaskResponse, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, ErrInvalidRequest
	}
	if err := validateDate(req.Date); err != nil {
		return nil, err
	}

	tags, err := s.resolveTags(ctx, userID, req.TagIDs)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateTask(ctx, userID, taskID, title, req.Date, tags); err != nil {
		return nil, err
	}

	task, err := s.repo.GetTask(ctx, userID, taskID)
	if err != nil {
		return nil, err
	}

	res := NewTaskResponse(task)
	return &res, nil
}

// DeleteTask removes one of the user's tasks.
func (s *Service) DeleteTask(ctx context.Context, userID, taskID string) error {
	return s.repo.DeleteTask(ctx, userID, taskID)
}

// SetTaskDone flips the completion flag on one of the user's tasks.
func (s *Service) SetTaskDone(ctx context.Context, userID, taskID string, done bool) (*TaskResponse, error) {
	if err := s.repo.SetTaskDone(ctx, userID, taskID, done); err != nil {
		return nil, err
	}

	task, err := s.repo.GetTask(ctx, userID, taskID)
	if err != nil {
		return nil, err
	}

	res := NewTaskResponse(task)
	return &res, nil
}

// ---- notes ----

// CreateNote adds a note to a date or to the general bucket.
func (s *Service) CreateNote(ctx context.Context, userID string, req CreateNoteRequest) (*NoteResponse, error) {
	body := strings.TrimSpace(req.Body)
	if body == "" {
		return nil, ErrInvalidRequest
	}
	if err := validateDate(req.Date); err != nil {
		return nil, err
	}

	tags, err := s.resolveTags(ctx, userID, req.TagIDs)
	if err != nil {
		return nil, err
	}

	note := &Note{
		ID:     uuid.NewString(),
		UserID: userID,
		Body:   body,
		Date:   req.Date,
		Tags:   tags,
	}
	if err := s.repo.CreateNote(ctx, note); err != nil {
		return nil, err
	}

	res := NewNoteResponse(note)
	return &res, nil
}

// UpdateNote rewrites a note's body, date and tag set.
func (s *Service) UpdateNote(ctx context.Context, userID, noteID string, req UpdateNoteRequest) (*NoteResponse, error) {
	body := strings.TrimSpace(req.Body)
	if body == "" {
		return nil, ErrInvalidRequest
	}
	if err := validateDate(req.Date); err != nil {
		return nil, err
	}

	tags, err := s.resolveTags(ctx, userID, req.TagIDs)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateNote(ctx, userID, noteID, body, req.Date, tags); err != nil {
		return nil, err
	}

	note, err := s.repo.GetNote(ctx, userID, noteID)
	if err != nil {
		return nil, err
	}

	res := NewNoteResponse(note)
	return &res, nil
}

// DeleteNote removes one of the user's notes.
func (s *Service) DeleteNote(ctx context.Context, userID, noteID string) error {
	return s.repo.DeleteNote(ctx, userID, noteID)
}

// ---- helpers ----

// resolveTags loads the requested tags, refusing ids the user does not own so a
// caller cannot attach someone else's tag.
func (s *Service) resolveTags(ctx context.Context, userID string, ids []string) ([]Tag, error) {
	unique := dedupe(ids)
	if len(unique) == 0 {
		return nil, nil
	}

	tags, err := s.repo.TagsByIDs(ctx, userID, unique)
	if err != nil {
		return nil, err
	}
	if len(tags) != len(unique) {
		return nil, ErrUnknownTag
	}
	return tags, nil
}

func dedupe(ids []string) []string {
	seen := make(map[string]struct{}, len(ids))
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func validateDate(date *string) error {
	if date == nil {
		return nil
	}
	if _, err := time.Parse(isoLayout, *date); err != nil {
		return ErrInvalidDate
	}
	return nil
}

func toTagResponses(tags []Tag) []TagResponse {
	out := make([]TagResponse, 0, len(tags))
	for i := range tags {
		out = append(out, NewTagResponse(&tags[i]))
	}
	return out
}

func toTaskResponses(tasks []Task) []TaskResponse {
	out := make([]TaskResponse, 0, len(tasks))
	for i := range tasks {
		out = append(out, NewTaskResponse(&tasks[i]))
	}
	return out
}

func toNoteResponses(notes []Note) []NoteResponse {
	out := make([]NoteResponse, 0, len(notes))
	for i := range notes {
		out = append(out, NewNoteResponse(&notes[i]))
	}
	return out
}

func isoDate(t time.Time) string {
	return t.Format(isoLayout)
}
