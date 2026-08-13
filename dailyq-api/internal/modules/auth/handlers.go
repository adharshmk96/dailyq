package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler adapts HTTP requests to the auth service.
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// POST /api/v1/auth/register — creates the account and logs the user in.
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if !bindJSON(c, &req) {
		return
	}

	res, err := h.svc.Register(c.Request.Context(), req, sessionContext(c))
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

// POST /api/v1/auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if !bindJSON(c, &req) {
		return
	}

	res, err := h.svc.Login(c.Request.Context(), req, sessionContext(c))
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// POST /api/v1/auth/logout — revokes the current session.
func (h *Handler) Logout(c *gin.Context) {
	if err := h.svc.Logout(c.Request.Context(), SessionIDFrom(c)); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// POST /api/v1/auth/change-password
func (h *Handler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if !bindJSON(c, &req) {
		return
	}

	err := h.svc.ChangePassword(c.Request.Context(), UserIDFrom(c), SessionIDFrom(c), req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed"})
}

// POST /api/v1/auth/forgot-password — always returns 202 so the endpoint does
// not disclose which emails are registered.
func (h *Handler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if !bindJSON(c, &req) {
		return
	}

	if err := h.svc.ForgotPassword(c.Request.Context(), req); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "if the email exists, a reset link has been sent"})
}

// POST /api/v1/auth/reset-password?token=...
func (h *Handler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if !bindJSON(c, &req) {
		return
	}

	if err := h.svc.ResetPassword(c.Request.Context(), c.Query("token"), req); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset"})
}

// GET /api/v1/auth/me
func (h *Handler) Me(c *gin.Context) {
	user, err := h.svc.Me(c.Request.Context(), UserIDFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, NewUserResponse(user))
}

func bindJSON(c *gin.Context, target any) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": &Error{
			Status:  http.StatusBadRequest,
			Code:    ErrInvalidRequest.Code,
			Message: err.Error(),
		}})
		return false
	}
	return true
}

func sessionContext(c *gin.Context) SessionContext {
	return SessionContext{
		UserAgent: c.Request.UserAgent(),
		IP:        c.ClientIP(),
	}
}

func respondError(c *gin.Context, err error) {
	e := AsError(err)
	if e.Status >= http.StatusInternalServerError {
		_ = c.Error(err)
	}
	c.JSON(e.Status, gin.H{"error": e})
}

func abortWithError(c *gin.Context, e *Error) {
	c.AbortWithStatusJSON(e.Status, gin.H{"error": e})
}
