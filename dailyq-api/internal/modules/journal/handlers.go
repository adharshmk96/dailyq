package journal

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dailyq-api/internal/modules/auth"
)

// Handler adapts HTTP requests to the journal service.
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// GET /api/v1/journal/overview
func (h *Handler) Overview(c *gin.Context) {
	res, err := h.svc.Overview(c.Request.Context(), auth.UserIDFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// PATCH /api/v1/journal/tasks/:id/done
func (h *Handler) SetTaskDone(c *gin.Context) {
	var req ToggleTaskRequest
	if !bindJSON(c, &req) {
		return
	}

	res, err := h.svc.SetTaskDone(c.Request.Context(), auth.UserIDFrom(c), c.Param("id"), *req.Done)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
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

func respondError(c *gin.Context, err error) {
	e := AsError(err)
	if e.Status >= http.StatusInternalServerError {
		_ = c.Error(err)
	}
	c.JSON(e.Status, gin.H{"error": e})
}
