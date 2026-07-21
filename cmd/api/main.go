package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/MRNaveed-stack/LinkPulse/internal/auth"
	"github.com/MRNaveed-stack/LinkPulse/internal/database"
	"github.com/MRNaveed-stack/LinkPulse/internal/handler"
	"github.com/MRNaveed-stack/LinkPulse/internal/logger"
	"github.com/MRNaveed-stack/LinkPulse/internal/middleware"
	"github.com/MRNaveed-stack/LinkPulse/internal/repository"
	"github.com/MRNaveed-stack/LinkPulse/internal/router"
	"github.com/MRNaveed-stack/LinkPulse/internal/service"
)

func main() {
	// Load .env file
	if err := loadEnv(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}
	logger.Init()

	// Verify Google OAuth config
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if googleClientID == "" || googleClientSecret == "" {
		log.Fatalf("Fatal: GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET environment variables must be configured")
	}

	// Connect to database
	dbPool, err := database.Connect()
	if err != nil {
		logger.Log.Error(
			"Failed to connect to database",
			"error", err,
		)
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Run database migrations
	migrationsDir := database.GetMigrationsDir()
	if err := database.RunMigrations(dbPool, migrationsDir); err != nil {
		logger.Log.Error("Failed to run database migrations", "error", err)
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(dbPool)
	prRepo := repository.NewPasswordResetTokenRepository(dbPool)
	linkRepo := repository.NewLinkRepository(dbPool)
	clickRepo := repository.NewClickRepository(dbPool)
	profileRepo := repository.NewProfileRepository(dbPool)

	// Initialize services
	authService := auth.NewAuthService(userRepo, prRepo, profileRepo)
	linkService := service.NewLinkService(linkRepo, clickRepo)
	profileService := service.NewProfileService(userRepo, profileRepo, linkRepo)
	analyticsService := service.NewAnalyticsService(linkRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	linkHandler := handler.NewLinkHandler(linkService, userRepo)
	profileHandler := handler.NewProfileHandler(profileService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	// Set Gin mode
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.DebugMode
	}
	gin.SetMode(ginMode)

	// Create router
	ginRouter := gin.New()
	ginRouter.Use(middleware.CORSMiddleware())
	ginRouter.Use(middleware.Recovery())
	ginRouter.Use(middleware.RequestLogger())

	// Setup routes
	router.SetupRoutes(ginRouter, authHandler, linkHandler, profileHandler, analyticsHandler)

	// Get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	logger.Log.Info("Server started", "port", port)
	if err := ginRouter.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func loadEnv() error {
	// Try current directory
	if err := godotenv.Load(); err == nil {
		return nil
	}

	// Try parent directory (when running from cmd/api)
	if err := godotenv.Load("../.env"); err == nil {
		return nil
	}

	// Try to find .env by walking up
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	dir := cwd
	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return godotenv.Load(envPath)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return os.ErrNotExist
}
