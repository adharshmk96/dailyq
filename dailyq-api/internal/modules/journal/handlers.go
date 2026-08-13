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

// GET /api/v1/journal/entries?date=YYYY-MM-DD
// Without a date the undated "general" bucket is returned.
func (h *Handler) Entries(c *gin.Context) {
	res, err := h.svc.Entries(c.Request.Context(), auth.UserIDFrom(c), dateParam(c))
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// GET /api/v1/journal/dates
func (h *Handler) Dates(c *gin.Context) {
	res, err := h.svc.Dates(c.Request.Context(), auth.UserIDFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// ---- tags ----

// GET /api/v1/journal/tags
func (h *Handler) ListTags(c *gin.Context) {
	res, err := h.svc.ListTags(c.Request.Context(), auth.UserIDFrom(c))
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"tags": res})
}

// POST /api/v1/journal/tags
func (h *Handler) CreateTag(c *gin.Context) {
	var req CreateTagRequest
	if !bindJSON(c, &req) {
		return
	}

	res, err := h.svc.CreateTag(c.Request.Context(), auth.UserIDFrom(c), req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

// PATCH /api/v1/journal/tags/:id
func (h *Handler) UpdateTag(c *gin.Context) {
	var req UpdateTagRequest
	if !bindJSON(c, &req) {
		return
	}

	res, err := h.svc.UpdateTag(c.Request.Context(), auth.UserIDFrom(c), c.Param("id"), req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// DELETE /api/v1/journal/tags/:id
func (h *Handler) DeleteTag(c *gin.Context) {
	if err := h.svc.DeleteTag(c.Request.Context(), auth.UserIDFrom(c), c.Param("id")); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ---- tasks ----

// POST /api/v1/journal/tasks
func (h *Handler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if !bindJSON(c, &req) {
		return
	}

	res, err := h.svc.CreateTask(c.Request.Context(), auth.UserIDFrom(c), req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

// PATCH /api/v1/journal/tasks/:id
func (h *Handler) UpdateTask(c *gin.Context) {
	var req UpdateTaskRequest
	if !bindJSON(c, &req) {
		return
	}

	res, err := h.svc.UpdateTask(c.Request.Context(), auth.UserIDFrom(c), c.Param("id"), req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// DELETE /api/v1/journal/tasks/:id
func (h *Handler) DeleteTask(c *gin.Context) {
	if err := h.svc.DeleteTask(c.Request.Context(), auth.UserIDFrom(c), c.Param("id")); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
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

// ---- notes ----

// POST /api/v1/journal/notes
func (h *Handler) CreateNote(c *gin.Context) {
	var req CreateNoteRequest
	if !bindJSON(c, &req) {
		return
	}

	res, err := h.svc.CreateNote(c.Request.Context(), auth.UserIDFrom(c), req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

// PATCH /api/v1/journal/notes/:id
func (h *Handler) UpdateNote(c *gin.Context) {
	var req UpdateNoteRequest
	if !bindJSON(c, &req) {
		return
	}

	res, err := h.svc.UpdateNote(c.Request.Context(), auth.UserIDFrom(c), c.Param("id"), req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// DELETE /api/v1/journal/notes/:id
func (h *Handler) DeleteNote(c *gin.Context) {
	if err := h.svc.DeleteNote(c.Request.Context(), auth.UserIDFrom(c), c.Param("id")); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// dateParam reads the optional ?date= scope; absent or empty means "general".
func dateParam(c *gin.Context) *string {
	date, ok := c.GetQuery("date")
	if !ok || date == "" {
		return nil
	}
	return &date
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
