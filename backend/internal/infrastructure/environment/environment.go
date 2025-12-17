package environment

import (
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	// PostgreSQL (Admin)
	PostgresDBUser     string
	PostgresDBPassword string
	PostgresDBName     string
	PostgresDBPort     string
	PostgresDBHost     string

	// PostgreSQL (App user - RLS)
	PostgresAppUser     string
	PostgresAppPassword string

	// JWT
	JWTSecret string

	// Server
	PublicAPIPort string

	// Mail
	SMTPHost string
	SMTPPort string
}

func NewEnvironment() *Environment {
	_ = godotenv.Load()

	return &Environment{
		PostgresDBUser:      getEnv("POSTGRES_DB_USER", "postgres"),
		PostgresDBPassword:  getEnv("POSTGRES_DB_PASSWORD", "postgres"),
		PostgresDBName:      getEnv("POSTGRES_DB_NAME", "good_todo_go"),
		PostgresDBPort:      getEnv("POSTGRES_DB_PORT", "5432"),
		PostgresDBHost:      getEnv("POSTGRES_DB_HOST", "localhost"),
		PostgresAppUser:     getEnv("POSTGRES_APP_USER", "app_user"),
		PostgresAppPassword: getEnv("POSTGRES_APP_PASSWORD", "app_password"),
		JWTSecret:           getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		PublicAPIPort:       getEnv("PUBLIC_API_PORT", "8000"),
		SMTPHost:            getEnv("SMTP_HOST", "localhost"),
		SMTPPort:            getEnv("SMTP_PORT", "1025"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetAdminDSN returns the DSN for admin connection (without RLS)
func (e *Environment) GetAdminDSN() string {
	return "postgres://" + e.PostgresDBUser + ":" + e.PostgresDBPassword +
		"@" + e.PostgresDBHost + ":" + e.PostgresDBPort +
		"/" + e.PostgresDBName + "?sslmode=disable"
}

// GetAppDSN returns the DSN for application connection (with RLS)
func (e *Environment) GetAppDSN() string {
	return "postgres://" + e.PostgresAppUser + ":" + e.PostgresAppPassword +
		"@" + e.PostgresDBHost + ":" + e.PostgresDBPort +
		"/" + e.PostgresDBName + "?sslmode=disable"
}
