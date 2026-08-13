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
	group.GET("/entries", m.Handler.Entries)
	group.GET("/dates", m.Handler.Dates)

	group.GET("/tags", m.Handler.ListTags)
	group.POST("/tags", m.Handler.CreateTag)
	group.PATCH("/tags/:id", m.Handler.UpdateTag)
	group.DELETE("/tags/:id", m.Handler.DeleteTag)

	group.POST("/tasks", m.Handler.CreateTask)
	group.PATCH("/tasks/:id", m.Handler.UpdateTask)
	group.PATCH("/tasks/:id/done", m.Handler.SetTaskDone)
	group.DELETE("/tasks/:id", m.Handler.DeleteTask)

	group.POST("/notes", m.Handler.CreateNote)
	group.PATCH("/notes/:id", m.Handler.UpdateNote)
	group.DELETE("/notes/:id", m.Handler.DeleteNote)
}
