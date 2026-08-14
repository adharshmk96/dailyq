package journal

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// Repository is the persistence boundary for the journal module.
type Repository interface {
	ListTags(ctx context.Context, userID string) ([]Tag, error)
	GetTag(ctx context.Context, userID, id string) (*Tag, error)
	CreateTag(ctx context.Context, tag *Tag) error
	UpdateTag(ctx context.Context, userID, id, name, color string) error
	DeleteTag(ctx context.Context, userID, id string) error
	TagNameTaken(ctx context.Context, userID, name, exceptID string) (bool, error)
	// TagsByIDs returns only the tags among ids that the user actually owns.
	TagsByIDs(ctx context.Context, userID string, ids []string) ([]Tag, error)
	GetTagByName(ctx context.Context, userID, name string) (*Tag, error)

	ListAllTasks(ctx context.Context, userID string) ([]Task, error)
	ListAllNotes(ctx context.Context, userID string) ([]Note, error)

	SaveTag(ctx context.Context, tag *Tag) error
	SaveTask(ctx context.Context, task *Task) error
	SaveNote(ctx context.Context, note *Note) error

	CountTasks(ctx context.Context, userID string) (int64, error)
	CountTasksDoneOn(ctx context.Context, userID, date string) (int64, error)
	CountOpenGeneralTasks(ctx context.Context, userID string) (int64, error)
	CountNotesSince(ctx context.Context, userID, since string) (int64, error)

	// ListTasks / ListNotes scope to one date, or to the undated "general"
	// bucket when date is nil.
	ListTasks(ctx context.Context, userID string, date *string) ([]Task, error)
	ListNotes(ctx context.Context, userID string, date *string) ([]Note, error)
	ListTasksForDate(ctx context.Context, userID, date string) ([]Task, error)
	ListRecentNotes(ctx context.Context, userID string, limit int) ([]Note, error)
	ListDatesWithItems(ctx context.Context, userID string) ([]string, error)

	GetTask(ctx context.Context, userID, id string) (*Task, error)
	CreateTask(ctx context.Context, task *Task) error
	UpdateTask(ctx context.Context, userID, id, title string, date *string, tags []Tag) error
	DeleteTask(ctx context.Context, userID, id string) error
	SetTaskDone(ctx context.Context, userID, id string, done bool) error

	GetNote(ctx context.Context, userID, id string) (*Note, error)
	CreateNote(ctx context.Context, note *Note) error
	UpdateNote(ctx context.Context, userID, id, body string, date *string, tags []Tag) error
	DeleteNote(ctx context.Context, userID, id string) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository returns a gorm-backed Repository.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ---- tags ----

func (r *repository) ListTags(ctx context.Context, userID string) ([]Tag, error) {
	var tags []Tag
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("name asc").
		Find(&tags).Error
	return tags, err
}

func (r *repository) GetTag(ctx context.Context, userID, id string) (*Tag, error) {
	var tag Tag
	err := r.db.WithContext(ctx).First(&tag, "id = ? AND user_id = ?", id, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTagNotFound
		}
		return nil, err
	}
	return &tag, nil
}

func (r *repository) CreateTag(ctx context.Context, tag *Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

func (r *repository) UpdateTag(ctx context.Context, userID, id, name, color string) error {
	res := r.db.WithContext(ctx).Model(&Tag{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]any{"name": name, "color": color})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrTagNotFound
	}
	return nil
}

// DeleteTag removes the tag and detaches it from every task and note, so the
// join rows do not outlive it.
func (r *repository) DeleteTag(ctx context.Context, userID, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM task_tags WHERE tag_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM note_tags WHERE tag_id = ?", id).Error; err != nil {
			return err
		}

		res := tx.Where("id = ? AND user_id = ?", id, userID).Delete(&Tag{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrTagNotFound
		}
		return nil
	})
}

func (r *repository) TagNameTaken(ctx context.Context, userID, name, exceptID string) (bool, error) {
	q := r.db.WithContext(ctx).Model(&Tag{}).
		Where("user_id = ? AND name = ? COLLATE NOCASE", userID, name)
	if exceptID != "" {
		q = q.Where("id <> ?", exceptID)
	}

	var n int64
	if err := q.Count(&n).Error; err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *repository) TagsByIDs(ctx context.Context, userID string, ids []string) ([]Tag, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var tags []Tag
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND id IN ?", userID, ids).
		Find(&tags).Error
	return tags, err
}

func (r *repository) GetTagByName(ctx context.Context, userID, name string) (*Tag, error) {
	var tag Tag
	err := r.db.WithContext(ctx).
		First(&tag, "user_id = ? AND name = ? COLLATE NOCASE", userID, name).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTagNotFound
		}
		return nil, err
	}
	return &tag, nil
}

func (r *repository) ListAllTasks(ctx context.Context, userID string) ([]Task, error) {
	var tasks []Task
	err := r.db.WithContext(ctx).
		Preload("Tags").
		Where("user_id = ?", userID).
		Order("date asc, created_at asc").
		Find(&tasks).Error
	return tasks, err
}

func (r *repository) ListAllNotes(ctx context.Context, userID string) ([]Note, error) {
	var notes []Note
	err := r.db.WithContext(ctx).
		Preload("Tags").
		Where("user_id = ?", userID).
		Order("date asc, created_at asc").
		Find(&notes).Error
	return notes, err
}

func (r *repository) SaveTag(ctx context.Context, tag *Tag) error {
	return r.db.WithContext(ctx).Save(tag).Error
}

func (r *repository) SaveTask(ctx context.Context, task *Task) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(task).Error; err != nil {
			return err
		}
		return tx.Model(task).Association("Tags").Replace(task.Tags)
	})
}

func (r *repository) SaveNote(ctx context.Context, note *Note) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(note).Error; err != nil {
			return err
		}
		return tx.Model(note).Association("Tags").Replace(note.Tags)
	})
}

// ---- counters ----

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

// ---- listing ----

// scopeDate narrows a query to one date, or to the undated bucket when nil.
func scopeDate(q *gorm.DB, date *string) *gorm.DB {
	if date == nil {
		return q.Where("date IS NULL")
	}
	return q.Where("date = ?", *date)
}

func (r *repository) ListTasks(ctx context.Context, userID string, date *string) ([]Task, error) {
	var tasks []Task
	q := r.db.WithContext(ctx).
		Preload("Tags").
		Where("user_id = ?", userID)
	err := scopeDate(q, date).
		Order("done asc, created_at desc").
		Find(&tasks).Error
	return tasks, err
}

func (r *repository) ListNotes(ctx context.Context, userID string, date *string) ([]Note, error) {
	var notes []Note
	q := r.db.WithContext(ctx).
		Preload("Tags").
		Where("user_id = ?", userID)
	err := scopeDate(q, date).
		Order("created_at desc").
		Find(&notes).Error
	return notes, err
}

func (r *repository) ListTasksForDate(ctx context.Context, userID, date string) ([]Task, error) {
	return r.ListTasks(ctx, userID, &date)
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

func (r *repository) ListDatesWithItems(ctx context.Context, userID string) ([]string, error) {
	var dates []string
	err := r.db.WithContext(ctx).Raw(`
		SELECT date FROM tasks
		WHERE user_id = ? AND date IS NOT NULL AND deleted_at IS NULL
		UNION
		SELECT date FROM notes
		WHERE user_id = ? AND date IS NOT NULL AND deleted_at IS NULL
		ORDER BY date asc
	`, userID, userID).Scan(&dates).Error
	return dates, err
}

// ---- tasks ----

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

func (r *repository) CreateTask(ctx context.Context, task *Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *repository) UpdateTask(ctx context.Context, userID, id, title string, date *string, tags []Tag) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Task{}).
			Where("id = ? AND user_id = ?", id, userID).
			Updates(map[string]any{"title": title, "date": date})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrTaskNotFound
		}
		return tx.Model(&Task{ID: id}).Association("Tags").Replace(tags)
	})
}

func (r *repository) DeleteTask(ctx context.Context, userID, id string) error {
	res := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&Task{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
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

// ---- notes ----

func (r *repository) GetNote(ctx context.Context, userID, id string) (*Note, error) {
	var note Note
	err := r.db.WithContext(ctx).
		Preload("Tags").
		First(&note, "id = ? AND user_id = ?", id, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoteNotFound
		}
		return nil, err
	}
	return &note, nil
}

func (r *repository) CreateNote(ctx context.Context, note *Note) error {
	return r.db.WithContext(ctx).Create(note).Error
}

func (r *repository) UpdateNote(ctx context.Context, userID, id, body string, date *string, tags []Tag) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Note{}).
			Where("id = ? AND user_id = ?", id, userID).
			Updates(map[string]any{"body": body, "date": date})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrNoteNotFound
		}
		return tx.Model(&Note{ID: id}).Association("Tags").Replace(tags)
	})
}

func (r *repository) DeleteNote(ctx context.Context, userID, id string) error {
	res := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&Note{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNoteNotFound
	}
	return nil
}
