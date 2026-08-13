package auth

import (
	"github.com/gin-gonic/gin"
)

// Context keys set by the auth middleware.
const (
	ContextUserID    = "auth.user_id"
	ContextSessionID = "auth.session_id"
	ContextEmail     = "auth.email"
)

// Middleware rejects requests without a valid bearer token backed by an active
// session, and stashes the identity on the gin context.
func Middleware(svc *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			abortWithError(c, ErrUnauthorized)
			return
		}

		claims, session, err := svc.Authenticate(c.Request.Context(), token)
		if err != nil {
			abortWithError(c, AsError(err))
			return
		}

		c.Set(ContextUserID, session.UserID)
		c.Set(ContextSessionID, session.ID)
		c.Set(ContextEmail, claims.Email)

		c.Next()
	}
}

// UserIDFrom returns the authenticated user id set by Middleware.
func UserIDFrom(c *gin.Context) string {
	v, _ := c.Get(ContextUserID)
	id, _ := v.(string)
	return id
}

// SessionIDFrom returns the authenticated session id set by Middleware.
func SessionIDFrom(c *gin.Context) string {
	v, _ := c.Get(ContextSessionID)
	id, _ := v.(string)
	return id
}
