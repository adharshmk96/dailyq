package journal

import "github.com/gin-gonic/gin"

// Module wires the journal handler and its auth middleware together.
type Module struct {
	Service    *Service
	Handler    *Handler
	Middleware gin.HandlerFunc
}

// RegisterRoutes mounts the journal endpoints under the given router group.
// Everything here is scoped to the authenticated user.
func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/journal", m.Middleware)

	group.GET("/overview", m.Handler.Overview)
	group.PATCH("/tasks/:id/done", m.Handler.SetTaskDone)
}
