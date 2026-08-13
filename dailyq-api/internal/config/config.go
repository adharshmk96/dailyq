package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all runtime configuration for the API.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Mode            string        `mapstructure:"mode"` // gin mode: debug | release | test
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	BaseURL         string        `mapstructure:"base_url"`
}

type DatabaseConfig struct {
	DSN     string `mapstructure:"dsn"`
	LogMode bool   `mapstructure:"log_mode"`
}

type AuthConfig struct {
	JWTSecret         string        `mapstructure:"jwt_secret"`
	JWTIssuer         string        `mapstructure:"jwt_issuer"`
	AccessTokenTTL    time.Duration `mapstructure:"access_token_ttl"`
	SessionTTL        time.Duration `mapstructure:"session_ttl"`
	ResetTokenTTL     time.Duration `mapstructure:"reset_token_ttl"`
	BcryptCost        int           `mapstructure:"bcrypt_cost"`
	PasswordMinLength int           `mapstructure:"password_min_length"`
	ResetPasswordPath string        `mapstructure:"reset_password_path"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`  // debug | info | warn | error
	Format string `mapstructure:"format"` // json | text
}

// Load reads configuration from (in order of precedence) env vars, config file,
// then built-in defaults. cfgFile may be empty, in which case the default
// search paths are used.
func Load(cfgFile string) (*Config, error) {
	v := viper.New()

	setDefaults(v)

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("/etc/dailyq")
	}

	v.SetEnvPrefix("DAILYQ")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		// A missing config file is fine — defaults and env vars still apply.
		if _, notFound := err.(viper.ConfigFileNotFoundError); !notFound {
			return nil, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "debug")
	v.SetDefault("server.read_timeout", 15*time.Second)
	v.SetDefault("server.write_timeout", 15*time.Second)
	v.SetDefault("server.shutdown_timeout", 10*time.Second)
	v.SetDefault("server.base_url", "http://localhost:8080")

	v.SetDefault("database.dsn", "dailyq.db")
	v.SetDefault("database.log_mode", false)

	v.SetDefault("auth.jwt_secret", "change-me-in-production")
	v.SetDefault("auth.jwt_issuer", "dailyq-api")
	v.SetDefault("auth.access_token_ttl", 24*time.Hour)
	v.SetDefault("auth.session_ttl", 720*time.Hour)
	v.SetDefault("auth.reset_token_ttl", 30*time.Minute)
	v.SetDefault("auth.bcrypt_cost", 12)
	v.SetDefault("auth.password_min_length", 8)
	v.SetDefault("auth.reset_password_path", "/reset-password")

	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "text")
}
