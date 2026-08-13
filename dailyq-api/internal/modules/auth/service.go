package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"dailyq-api/internal/config"
)

// SessionContext carries per-request metadata recorded on a session.
type SessionContext struct {
	UserAgent string
	IP        string
}

// Service holds the auth business logic.
type Service struct {
	repo   Repository
	cfg    config.AuthConfig
	base   string
	logger *slog.Logger
	now    func() time.Time
}

// NewService builds the auth service. baseURL is used to render reset links.
func NewService(repo Repository, cfg config.AuthConfig, baseURL string, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		cfg:    cfg,
		base:   strings.TrimRight(baseURL, "/"),
		logger: logger,
		now:    time.Now,
	}
}

// Register creates an account and logs the user in immediately.
func (s *Service) Register(ctx context.Context, req RegisterRequest, sc SessionContext) (*AuthResponse, error) {
	email := normalizeEmail(req.Email)

	if _, err := s.repo.GetUserByEmail(ctx, email); err == nil {
		return nil, ErrEmailTaken
	} else if !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}

	if err := s.validatePassword(req.Password); err != nil {
		return nil, err
	}

	hash, err := hashPassword(req.Password, s.cfg.BcryptCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:           newID(),
		Email:        email,
		Name:         strings.TrimSpace(req.Name),
		PasswordHash: hash,
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	s.logger.Info("user registered", "user_id", user.ID, "email", user.Email)

	return s.issueSession(ctx, user, sc)
}

// Login verifies credentials and issues a new session + JWT.
func (s *Service) Login(ctx context.Context, req LoginRequest, sc SessionContext) (*AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, normalizeEmail(req.Email))
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !checkPassword(user.PasswordHash, req.Password) {
		return nil, ErrInvalidCredentials
	}

	return s.issueSession(ctx, user, sc)
}

// Logout revokes the session the current token is bound to.
func (s *Service) Logout(ctx context.Context, sessionID string) error {
	if err := s.repo.RevokeSession(ctx, sessionID, s.now()); err != nil {
		return err
	}
	s.logger.Info("session revoked", "session_id", sessionID)
	return nil
}

// ChangePassword updates the password and revokes every other session.
func (s *Service) ChangePassword(ctx context.Context, userID, currentSessionID string, req ChangePasswordRequest) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if !checkPassword(user.PasswordHash, req.CurrentPassword) {
		return ErrInvalidCredentials
	}
	if req.CurrentPassword == req.NewPassword {
		return ErrSamePassword
	}
	if err := s.validatePassword(req.NewPassword); err != nil {
		return err
	}

	hash, err := hashPassword(req.NewPassword, s.cfg.BcryptCost)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateUserPassword(ctx, user.ID, hash); err != nil {
		return err
	}

	// Every other session is signed out; the caller keeps the session they
	// used to make the change.
	if err := s.repo.RevokeUserSessionsExcept(ctx, user.ID, currentSessionID, s.now()); err != nil {
		return err
	}

	s.logger.Info("password changed", "user_id", user.ID)
	return nil
}

// ForgotPassword issues a reset token and logs the reset URL. It never reveals
// whether the email exists.
func (s *Service) ForgotPassword(ctx context.Context, req ForgotPasswordRequest) error {
	email := normalizeEmail(req.Email)

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			s.logger.Info("password reset requested for unknown email", "email", email)
			return nil
		}
		return err
	}

	now := s.now()
	if err := s.repo.InvalidateUserResetTokens(ctx, user.ID, now); err != nil {
		return err
	}

	token, tokenHash, err := generateResetToken()
	if err != nil {
		return err
	}

	record := &PasswordResetToken{
		ID:        newID(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: now.Add(s.cfg.ResetTokenTTL),
	}
	if err := s.repo.CreateResetToken(ctx, record); err != nil {
		return err
	}

	s.logger.Info("password reset link generated",
		"user_id", user.ID,
		"email", user.Email,
		"expires_at", record.ExpiresAt,
		"reset_url", s.resetURL(token),
	)

	return nil
}

// ResetPassword consumes a reset token and sets a new password, revoking all
// existing sessions.
func (s *Service) ResetPassword(ctx context.Context, token string, req ResetPasswordRequest) error {
	if strings.TrimSpace(token) == "" {
		return ErrInvalidResetToken
	}

	record, err := s.repo.GetResetTokenByHash(ctx, hashToken(token))
	if err != nil {
		return err
	}

	now := s.now()
	if !record.IsUsable(now) {
		return ErrInvalidResetToken
	}
	if err := s.validatePassword(req.Password); err != nil {
		return err
	}

	hash, err := hashPassword(req.Password, s.cfg.BcryptCost)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateUserPassword(ctx, record.UserID, hash); err != nil {
		return err
	}
	if err := s.repo.MarkResetTokenUsed(ctx, record.ID, now); err != nil {
		return err
	}
	if err := s.repo.RevokeUserSessions(ctx, record.UserID, now); err != nil {
		return err
	}

	s.logger.Info("password reset", "user_id", record.UserID)
	return nil
}

// Me returns the authenticated user.
func (s *Service) Me(ctx context.Context, userID string) (*User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

// Authenticate validates a bearer token and its backing session.
func (s *Service) Authenticate(ctx context.Context, tokenString string) (*Claims, *Session, error) {
	claims, err := parseToken(s.cfg.JWTSecret, tokenString)
	if err != nil {
		return nil, nil, ErrInvalidToken
	}

	session, err := s.repo.GetSessionByID(ctx, claims.ID)
	if err != nil {
		return nil, nil, err
	}
	if !session.IsActive(s.now()) || session.UserID != claims.Subject {
		return nil, nil, ErrSessionExpired
	}

	return claims, session, nil
}

func (s *Service) issueSession(ctx context.Context, user *User, sc SessionContext) (*AuthResponse, error) {
	now := s.now()

	session := &Session{
		ID:        newID(),
		UserID:    user.ID,
		UserAgent: sc.UserAgent,
		IP:        sc.IP,
		ExpiresAt: now.Add(s.cfg.SessionTTL),
	}
	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	token, expiresAt, err := signToken(
		s.cfg.JWTSecret, s.cfg.JWTIssuer,
		user.ID, user.Email, session.ID,
		s.cfg.AccessTokenTTL, now,
	)
	if err != nil {
		return nil, err
	}

	s.logger.Info("session created", "user_id", user.ID, "session_id", session.ID)

	return &AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		SessionID: session.ID,
		User:      NewUserResponse(user),
	}, nil
}

func (s *Service) validatePassword(password string) error {
	min := s.cfg.PasswordMinLength
	if min <= 0 {
		min = 8
	}
	if len(password) < min {
		return &Error{
			Status:  ErrWeakPassword.Status,
			Code:    ErrWeakPassword.Code,
			Message: fmt.Sprintf("password must be at least %d characters", min),
		}
	}
	return nil
}

func (s *Service) resetURL(token string) string {
	path := s.cfg.ResetPasswordPath
	if path == "" {
		path = "/reset-password"
	}
	return fmt.Sprintf("%s%s?token=%s", s.base, path, url.QueryEscape(token))
}
