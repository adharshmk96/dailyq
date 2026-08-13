package journal

import (
	"context"
	"log/slog"
	"time"
)

// recentNotesLimit caps the "Recent notes" list on the overview page.
const recentNotesLimit = 5

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
		Tags:        make([]TagResponse, 0, len(tags)),
		TodayTasks:  make([]TaskResponse, 0, len(todayTasks)),
		RecentNotes: make([]NoteResponse, 0, len(recentNotes)),
	}
	for i := range tags {
		res.Tags = append(res.Tags, NewTagResponse(&tags[i]))
	}
	for i := range todayTasks {
		res.TodayTasks = append(res.TodayTasks, NewTaskResponse(&todayTasks[i]))
	}
	for i := range recentNotes {
		res.RecentNotes = append(res.RecentNotes, NewNoteResponse(&recentNotes[i]))
	}

	return res, nil
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

func isoDate(t time.Time) string {
	return t.Format("2006-01-02")
}
