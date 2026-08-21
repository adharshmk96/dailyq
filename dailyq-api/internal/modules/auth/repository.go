package auth

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// Repository is the persistence boundary for the auth module.
type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUserPassword(ctx context.Context, userID, passwordHash string) error
	UpdateUserProfile(ctx context.Context, userID, name, email string) error

	CreateSession(ctx context.Context, session *Session) error
	GetSessionByID(ctx context.Context, id string) (*Session, error)
	RevokeSession(ctx context.Context, id string, at time.Time) error
	RevokeUserSessions(ctx context.Context, userID string, at time.Time) error
	RevokeUserSessionsExcept(ctx context.Context, userID, keepSessionID string, at time.Time) error

	CreateResetToken(ctx context.Context, token *PasswordResetToken) error
	GetResetTokenByHash(ctx context.Context, hash string) (*PasswordResetToken, error)
	MarkResetTokenUsed(ctx context.Context, id string, at time.Time) error
	InvalidateUserResetTokens(ctx context.Context, userID string, at time.Time) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository returns a gorm-backed Repository.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *repository) GetUserByID(ctx context.Context, id string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, wrapNotFound(err, ErrUserNotFound)
	}
	return &user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return nil, wrapNotFound(err, ErrUserNotFound)
	}
	return &user, nil
}

func (r *repository) UpdateUserPassword(ctx context.Context, userID, passwordHash string) error {
	return r.db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", userID).
		Update("password_hash", passwordHash).Error
}

func (r *repository) UpdateUserProfile(ctx context.Context, userID, name, email string) error {
	return r.db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", userID).
		Updates(map[string]any{
			"name":  name,
			"email": email,
		}).Error
}

func (r *repository) CreateSession(ctx context.Context, session *Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *repository) GetSessionByID(ctx context.Context, id string) (*Session, error) {
	var session Session
	if err := r.db.WithContext(ctx).First(&session, "id = ?", id).Error; err != nil {
		return nil, wrapNotFound(err, ErrSessionExpired)
	}
	return &session, nil
}

func (r *repository) RevokeSession(ctx context.Context, id string, at time.Time) error {
	return r.db.WithContext(ctx).
		Model(&Session{}).
		Where("id = ? AND revoked_at IS NULL", id).
		Update("revoked_at", at).Error
}

func (r *repository) RevokeUserSessions(ctx context.Context, userID string, at time.Time) error {
	return r.db.WithContext(ctx).
		Model(&Session{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", at).Error
}

func (r *repository) RevokeUserSessionsExcept(ctx context.Context, userID, keepSessionID string, at time.Time) error {
	return r.db.WithContext(ctx).
		Model(&Session{}).
		Where("user_id = ? AND id <> ? AND revoked_at IS NULL", userID, keepSessionID).
		Update("revoked_at", at).Error
}

func (r *repository) CreateResetToken(ctx context.Context, token *PasswordResetToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *repository) GetResetTokenByHash(ctx context.Context, hash string) (*PasswordResetToken, error) {
	var token PasswordResetToken
	if err := r.db.WithContext(ctx).First(&token, "token_hash = ?", hash).Error; err != nil {
		return nil, wrapNotFound(err, ErrInvalidResetToken)
	}
	return &token, nil
}

func (r *repository) MarkResetTokenUsed(ctx context.Context, id string, at time.Time) error {
	return r.db.WithContext(ctx).
		Model(&PasswordResetToken{}).
		Where("id = ?", id).
		Update("used_at", at).Error
}

func (r *repository) InvalidateUserResetTokens(ctx context.Context, userID string, at time.Time) error {
	return r.db.WithContext(ctx).
		Model(&PasswordResetToken{}).
		Where("user_id = ? AND used_at IS NULL", userID).
		Update("used_at", at).Error
}

func wrapNotFound(err error, notFound *Error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return notFound
	}
	return err
}
