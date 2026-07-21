package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/MRNaveed-stack/LinkPulse/internal/logger"
)

// Connect establishes a connection to the PostgreSQL database with retry logic.
func Connect() (*pgxpool.Pool, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	// Read environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Set default values if not provided
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}

	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s ",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	// Parse connection string
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		logger.Log.Error(
			"Failed to connect to database",
			"error", err,
		)
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	maxRetries := 10
	backoff := 1 * time.Second

	var pool *pgxpool.Pool
	for attempt := 1; attempt <= maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		pool, err = pgxpool.NewWithConfig(ctx, config)
		if err == nil {
			err = pool.Ping(ctx)
		}
		cancel()

		if err == nil {
			logger.Log.Info("Database connected")
			return pool, nil
		}

		if pool != nil {
			pool.Close()
		}

		logger.Log.Info(fmt.Sprintf("Database connection failed (attempt %d/%d), retrying in %v...", attempt, maxRetries, backoff), "error", err)
		time.Sleep(backoff)
		if backoff < 8*time.Second {
			backoff *= 2
		}
	}

	logger.Log.Error("Failed to connect to database after retries", "error", err)
	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

