package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Zitadel  ZitadelConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret      string
	ExpiryHours int
}

type ZitadelConfig struct {
	URL          string
	ClientID     string
	ClientSecret string
}

func Load() (*Config, error) {
	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	jwtExpiryHours, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRY_HOURS: %w", err)
	}

	return &Config{
		Server: ServerConfig{
			Port: serverPort,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DATABASE_HOST", getEnv("DB_HOST", "localhost")),
			Port:     dbPort,
			User:     getEnv("DB_USER", "tsunagu"),
			Password: getEnv("DB_PASSWORD", "tsunagu_password"),
			DBName:   getEnv("DB_NAME", "tsunagu_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			ExpiryHours: jwtExpiryHours,
		},
		Zitadel: ZitadelConfig{
			URL:          getEnv("ZITADEL_URL", ""),
			ClientID:     getEnv("ZITADEL_CLIENT_ID", ""),
			ClientSecret: getEnv("ZITADEL_CLIENT_SECRET", ""),
		},
	}, nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
