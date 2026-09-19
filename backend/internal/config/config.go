package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type AppConfig struct {
	Env  string `mapstructure:"env"`
	Port int    `mapstructure:"port"`
	CORSOrigins []string `mapstructure:"cors_origins"`
}

type DatabaseConfig struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	Name           string `mapstructure:"name"`
	SSLMode        string `mapstructure:"sslmode"`
	ConnectTimeout int    `mapstructure:"connect_timeout"`
}

type JWTConfig struct {
	Secret    string        `mapstructure:"secret"`
	ExpiresIn time.Duration `mapstructure:"expires_in"`
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode, c.ConnectTimeout,
	)
}

func Load(configPath, envPath string, configRequired, envRequired bool) (*Config, error) {
	if err := godotenv.Load(envPath); err != nil {
		if envRequired || !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("load env file %q: %w", envPath, err)
		}
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	v.SetDefault("app.env", "development")
	v.SetDefault("app.port", 8080)
	v.SetDefault("app.cors_origins", []string{"http://localhost:5173"})
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.connect_timeout", 5)
	v.SetDefault("jwt.expires_in", "24h")

	bindings := map[string]string{
		"app.env":                  "APP_ENV",
		"app.port":                 "APP_PORT",
		"app.cors_origins":         "APP_CORS_ORIGINS",
		"database.host":            "POSTGRES_HOST",
		"database.port":            "POSTGRES_PORT",
		"database.user":            "POSTGRES_USER",
		"database.password":        "POSTGRES_PASSWORD",
		"database.name":            "POSTGRES_DB",
		"database.sslmode":         "POSTGRES_SSLMODE",
		"database.connect_timeout": "POSTGRES_CONNECT_TIMEOUT",
		"jwt.secret":               "JWT_SECRET",
		"jwt.expires_in":           "JWT_EXPIRES_IN",
	}

	for key, env := range bindings {
		if err := v.BindEnv(key, env); err != nil {
			return nil, fmt.Errorf("bind env %s: %w", env, err)
		}
	}

	if err := v.ReadInConfig(); err != nil {
		if configRequired || !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("read config %q: %w", configPath, err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Database.User == "" {
		return errors.New("config: POSTGRES_USER wajib diisi di .env")
	}
	if c.Database.Password == "" {
		return errors.New("config: POSTGRES_PASSWORD wajib diisi di .env")
	}
	if c.Database.Name == "" {
		return errors.New("config: POSTGRES_DB wajib diisi di .env")
	}
	if len(c.JWT.Secret) < 32 {
		return errors.New("config: JWT_SECRET wajib diisi di .env, minimal 32 karakter")
	}
	if c.Database.ConnectTimeout <= 0 {
		return errors.New("config: database.connect_timeout harus positif, dalam detik")
	}
	if c.JWT.ExpiresIn <= 0 {
		return errors.New("config: jwt.expires_in harus positif, misalnya 24h")
	}
	return nil
}
