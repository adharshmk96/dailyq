package journal

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// Repository is the persistence boundary for the journal module.
type Repository interface {
	ListTags(ctx context.Context, userID string) ([]Tag, error)

	CountTasks(ctx context.Context, userID string) (int64, error)
	CountTasksDoneOn(ctx context.Context, userID, date string) (int64, error)
	CountOpenGeneralTasks(ctx context.Context, userID string) (int64, error)
	CountNotesSince(ctx context.Context, userID, since string) (int64, error)

	ListTasksForDate(ctx context.Context, userID, date string) ([]Task, error)
	ListRecentNotes(ctx context.Context, userID string, limit int) ([]Note, error)

	GetTask(ctx context.Context, userID, id string) (*Task, error)
	SetTaskDone(ctx context.Context, userID, id string, done bool) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository returns a gorm-backed Repository.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) ListTags(ctx context.Context, userID string) ([]Tag, error) {
	var tags []Tag
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("name asc").
		Find(&tags).Error
	return tags, err
}

func (r *repository) CountTasks(ctx context.Context, userID string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&Task{}).Where("user_id = ?", userID).Count(&n).Error
	return n, err
}

func (r *repository) CountTasksDoneOn(ctx context.Context, userID, date string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&Task{}).
		Where("user_id = ? AND date = ? AND done = ?", userID, date, true).
		Count(&n).Error
	return n, err
}

func (r *repository) CountOpenGeneralTasks(ctx context.Context, userID string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&Task{}).
		Where("user_id = ? AND date IS NULL AND done = ?", userID, false).
		Count(&n).Error
	return n, err
}

func (r *repository) CountNotesSince(ctx context.Context, userID, since string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&Note{}).
		Where("user_id = ? AND date IS NOT NULL AND date >= ?", userID, since).
		Count(&n).Error
	return n, err
}

func (r *repository) ListTasksForDate(ctx context.Context, userID, date string) ([]Task, error) {
	var tasks []Task
	err := r.db.WithContext(ctx).
		Preload("Tags").
		Where("user_id = ? AND date = ?", userID, date).
		Order("done asc, created_at asc").
		Find(&tasks).Error
	return tasks, err
}

func (r *repository) ListRecentNotes(ctx context.Context, userID string, limit int) ([]Note, error) {
	var notes []Note
	err := r.db.WithContext(ctx).
		Preload("Tags").
		Where("user_id = ?", userID).
		Order("date desc, created_at desc").
		Limit(limit).
		Find(&notes).Error
	return notes, err
}

func (r *repository) GetTask(ctx context.Context, userID, id string) (*Task, error) {
	var task Task
	err := r.db.WithContext(ctx).
		Preload("Tags").
		First(&task, "id = ? AND user_id = ?", id, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return &task, nil
}

func (r *repository) SetTaskDone(ctx context.Context, userID, id string, done bool) error {
	res := r.db.WithContext(ctx).Model(&Task{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("done", done)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}
