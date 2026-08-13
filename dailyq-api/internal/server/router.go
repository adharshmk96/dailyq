package server

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"dailyq-api/internal/config"
	"dailyq-api/internal/modules/auth"
	"dailyq-api/internal/modules/journal"
)

// NewRouter builds the gin engine with all middleware and module routes.
func NewRouter(cfg *config.Config, db *gorm.DB, log *slog.Logger) *gin.Engine {
	gin.SetMode(ginMode(cfg.Server.Mode))

	engine := gin.New()
	engine.Use(requestLogger(log), gin.Recovery())

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	authSvc := auth.NewService(auth.NewRepository(db), cfg.Auth, cfg.Server.BaseURL, log)
	authModule := &auth.Module{
		Service:    authSvc,
		Handler:    auth.NewHandler(authSvc),
		Middleware: auth.Middleware(authSvc),
	}

	journalModule := &journal.Module{
		Service:    journal.NewService(journal.NewRepository(db), log),
		Middleware: authModule.Middleware,
	}
	journalModule.Handler = journal.NewHandler(journalModule.Service)

	v1 := engine.Group("/api/v1")
	authModule.RegisterRoutes(v1)
	journalModule.RegisterRoutes(v1)

	spa, err := spaHandler(log)
	if err != nil {
		log.Error("ui assets unavailable, serving api only", "error", err)
	}

	engine.NoRoute(func(c *gin.Context) {
		// API paths always answer with JSON; anything else is the UI's.
		if spa == nil || strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, notFoundBody())
			return
		}
		spa(c)
	})

	return engine
}

func ginMode(mode string) string {
	switch mode {
	case gin.ReleaseMode, gin.TestMode, gin.DebugMode:
		return mode
	default:
		return gin.ReleaseMode
	}
}

// requestLogger logs one line per request via slog.
func requestLogger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		attrs := []any{
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"duration", time.Since(start).String(),
			"ip", c.ClientIP(),
		}
		if len(c.Errors) > 0 {
			attrs = append(attrs, "errors", c.Errors.String())
			log.Error("request failed", attrs...)
			return
		}
		log.Info("request", attrs...)
	}
}
