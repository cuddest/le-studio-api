package config

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config stores runtime application configuration values.
type Config struct {
	Port                 string
	Environment          string
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	DBSSLMode            string
	JWTSecret            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	CloudinaryCloudName  string
	CloudinaryAPIKey     string
	CloudinaryAPISecret  string
	AllowedOrigins       []string
}

// Load reads environment variables and returns config.
func Load() (Config, error) {
	_ = godotenv.Load()
	access, _ := time.ParseDuration(get("ACCESS_TOKEN_DURATION", "15m"))
	refresh, _ := time.ParseDuration(get("REFRESH_TOKEN_DURATION", "168h"))
	return Config{
		Port:                 get("PORT", "8080"),
		Environment:          get("ENV", "development"),
		DBHost:               get("DB_HOST", "localhost"),
		DBPort:               get("DB_PORT", "5432"),
		DBUser:               get("DB_USER", "postgres"),
		DBPassword:           get("DB_PASSWORD", "postgres"),
		DBName:               get("DB_NAME", "lestudio"),
		DBSSLMode:            get("DB_SSLMODE", "disable"),
		JWTSecret:            get("JWT_SECRET", "change-me"),
		AccessTokenDuration:  access,
		RefreshTokenDuration: refresh,
		CloudinaryCloudName:  get("CLOUDINARY_CLOUD_NAME", ""),
		CloudinaryAPIKey:     get("CLOUDINARY_API_KEY", ""),
		CloudinaryAPISecret:  get("CLOUDINARY_API_SECRET", ""),
		AllowedOrigins:       strings.Split(get("ALLOWED_ORIGINS", "*"), ","),
	}, nil
}

func get(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
