package journal

import (
	"time"

	"gorm.io/gorm"

	"dailyq-api/internal/modules/auth"
)

// Tag is a user-owned label that can be attached to tasks and notes.
type Tag struct {
	ID        string         `gorm:"primaryKey;size:36" json:"id"`
	UserID    string         `gorm:"index;size:36;not null" json:"-"`
	Name      string         `gorm:"size:64;not null" json:"name"`
	Color     string         `gorm:"size:16;not null;default:neutral" json:"color"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User auth.User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

// Task is a to-do item. Date is nil for "general" (undated) items.
type Task struct {
	ID        string         `gorm:"primaryKey;size:36" json:"id"`
	UserID    string         `gorm:"index;size:36;not null" json:"-"`
	Title     string         `gorm:"size:512;not null" json:"title"`
	Done      bool           `gorm:"not null;default:false" json:"done"`
	Date      *string        `gorm:"index;size:10" json:"date"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User auth.User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Tags []Tag     `gorm:"many2many:task_tags;constraint:OnDelete:CASCADE" json:"-"`
}

// Note is a free-form entry. Date is nil for "general" (undated) items.
type Note struct {
	ID        string         `gorm:"primaryKey;size:36" json:"id"`
	UserID    string         `gorm:"index;size:36;not null" json:"-"`
	Body      string         `gorm:"type:text;not null" json:"body"`
	Date      *string        `gorm:"index;size:10" json:"date"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User auth.User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Tags []Tag     `gorm:"many2many:note_tags;constraint:OnDelete:CASCADE" json:"-"`
}

// ---- request payloads ----

type ToggleTaskRequest struct {
	Done *bool `json:"done" binding:"required"`
}

// ---- response payloads ----

type TagResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type TaskResponse struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Done   bool     `json:"done"`
	Date   *string  `json:"date"`
	TagIDs []string `json:"tag_ids"`
}

type NoteResponse struct {
	ID     string   `json:"id"`
	Body   string   `json:"body"`
	Date   *string  `json:"date"`
	TagIDs []string `json:"tag_ids"`
}

// OverviewStats are the counters shown on the dashboard overview cards.
type OverviewStats struct {
	TotalTasks     int64 `json:"total_tasks"`
	CompletedToday int64 `json:"completed_today"`
	OpenGeneral    int64 `json:"open_general"`
	NotesThisWeek  int64 `json:"notes_this_week"`
}

// OverviewResponse is everything the overview page renders in one round trip.
type OverviewResponse struct {
	Today       string         `json:"today"`
	Stats       OverviewStats  `json:"stats"`
	Tags        []TagResponse  `json:"tags"`
	TodayTasks  []TaskResponse `json:"today_tasks"`
	RecentNotes []NoteResponse `json:"recent_notes"`
}

func NewTagResponse(t *Tag) TagResponse {
	return TagResponse{ID: t.ID, Name: t.Name, Color: t.Color}
}

func NewTaskResponse(t *Task) TaskResponse {
	return TaskResponse{
		ID:     t.ID,
		Title:  t.Title,
		Done:   t.Done,
		Date:   t.Date,
		TagIDs: tagIDs(t.Tags),
	}
}

func NewNoteResponse(n *Note) NoteResponse {
	return NoteResponse{
		ID:     n.ID,
		Body:   n.Body,
		Date:   n.Date,
		TagIDs: tagIDs(n.Tags),
	}
}

func tagIDs(tags []Tag) []string {
	ids := make([]string, 0, len(tags))
	for i := range tags {
		ids = append(ids, tags[i].ID)
	}
	return ids
}
