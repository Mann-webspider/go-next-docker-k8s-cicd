package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	e := echo.New()

	// Load .env for local dev (ignore error in K8s/Docker where ENV vars are pre-set)
	if err := godotenv.Load(); err != nil {
		e.Logger.Warn("No .env file found, relying on environment variables")
	}

	// Middlewares
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Adjust for production Next.js domain
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))

	// Database Configuration
	host := getEnvOrDefault("POSTGRES_HOST", "localhost")
	user := getEnvOrDefault("POSTGRES_USER", "postgres")
	password := getEnvOrDefault("POSTGRES_PASSWORD", "postgres")
	dbName := getEnvOrDefault("POSTGRES_DB", "events_db")
	port := getEnvOrDefault("POSTGRES_PORT", "5432")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
		host, user, password, dbName, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Auto-migrate schema
	if err := db.AutoMigrate(&Button{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Connection Pool Settings
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB instance: %v", err)
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// --- Routes ---

	// Kubernetes Liveness/Readiness Probe Endpoint
	e.GET("/health", func(c *echo.Context) error {
		if err := sqlDB.Ping(); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "unhealthy", "db": "disconnected"})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	})

	// Get total clicks
	e.GET("/clicks", func(c *echo.Context) error {
		var button Button
		db.FirstOrCreate(&button, Button{ID: 1})
		return c.JSON(http.StatusOK, map[string]interface{}{
			"count": button.ClickedCount,
		})
	})

	// Increment click endpoint
	e.POST("/clicked", func(c *echo.Context) error {
		count, err := IncrementClick(db)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to record click"})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Button clicked!",
			"count":   count,
		})
	})

	// Start Server
	serverPort := getEnvOrDefault("PORT", "8080")
	if err := e.Start(":" + serverPort); err != nil {
		e.Logger.Error("Shutting down the server: ", err)
	}
}

// Helper to fallback to default env vars
func getEnvOrDefault(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
