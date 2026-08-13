package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"dailyq-api/internal/config"
	"dailyq-api/internal/database"
	"dailyq-api/internal/logger"
	"dailyq-api/internal/server"
)

var (
	cfgFile string
	cfg     *config.Config
	log     *slog.Logger
)

var rootCmd = &cobra.Command{
	Use:   "dailyq-api",
	Short: "DailyQ API server",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if cfg, err = config.Load(cfgFile); err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		log = logger.New(cfg.Log)
		logger.SetDefault(log)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runServe()
	},
	SilenceUsage: true,
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the HTTP API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runServe()
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Apply database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return withDB(func(db *gorm.DB) error {
			return database.Migrate(db, log)
		})
	},
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with development data",
	RunE: func(cmd *cobra.Command, args []string) error {
		return withDB(func(db *gorm.DB) error {
			if err := database.Migrate(db, log); err != nil {
				return err
			}
			return database.Seed(db, log)
		})
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to config file (default: ./config.yaml)")
	rootCmd.AddCommand(serveCmd, migrateCmd, seedCmd)
}

// Execute runs the CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func runServe() error {
	return withDB(func(db *gorm.DB) error {
		if err := database.Migrate(db, log); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
		return server.New(cfg, db, log).Run()
	})
}

// withDB opens the database, runs fn, and closes the connection.
func withDB(fn func(db *gorm.DB) error) error {
	db, err := database.Open(cfg.Database, log)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer func() {
		if cerr := database.Close(db); cerr != nil {
			log.Error("closing database", "error", cerr)
		}
	}()

	return fn(db)
}
