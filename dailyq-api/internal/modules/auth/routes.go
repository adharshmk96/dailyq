package auth

import "github.com/gin-gonic/gin"

// Module wires the auth handler, service and middleware together.
type Module struct {
	Service    *Service
	Handler    *Handler
	Middleware gin.HandlerFunc
}

// RegisterRoutes mounts the auth endpoints under the given router group.
func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	group := rg.Group("/auth")

	// Public
	group.POST("/register", m.Handler.Register)
	group.POST("/login", m.Handler.Login)
	group.POST("/forgot-password", m.Handler.ForgotPassword)
	group.POST("/reset-password", m.Handler.ResetPassword)

	// Authenticated
	protected := group.Group("", m.Middleware)
	protected.POST("/logout", m.Handler.Logout)
	protected.POST("/change-password", m.Handler.ChangePassword)
	protected.GET("/me", m.Handler.Me)
	protected.PATCH("/me", m.Handler.UpdateProfile)
}
