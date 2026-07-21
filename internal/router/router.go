package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MRNaveed-stack/LinkPulse/internal/config"
	"github.com/MRNaveed-stack/LinkPulse/internal/handler"
	"github.com/MRNaveed-stack/LinkPulse/internal/middleware"
)

// SetupRoutes wires all endpoints.
func SetupRoutes(
	r *gin.Engine,
	authHandler *handler.AuthHandler,
	linkHandler *handler.LinkHandler,
	profileHandler *handler.ProfileHandler,
	analyticsHandler *handler.AnalyticsHandler,
) {
	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Auth
	auth := r.Group("")
	{
		auth.POST(
			"/auth/register",
			middleware.RateLimit(config.RegisterRate, config.RegisterBurst),
			authHandler.Register,
		)

		auth.POST(
			"/auth/login",
			middleware.RateLimit(config.LoginRate, config.LoginBurst),
			authHandler.Login,
		)

		auth.POST(
			"/auth/forgot-password",
			middleware.RateLimit(config.ForgotPasswordRate, config.ForgotPasswordBurst),
			authHandler.ForgotPassword,
		)

		auth.POST("/auth/reset-password", authHandler.ResetPassword)
		auth.GET("/auth/google/login", authHandler.GoogleLogin)
		auth.GET("/auth/google/callback", authHandler.GoogleCallback)
		auth.POST("/auth/google/token", authHandler.GoogleTokenExchange)
		auth.POST("/auth/refresh", authHandler.Refresh)
	}

	// Protected
	protected := r.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/me", authHandler.GetMe)

		protected.POST("/links", linkHandler.CreateLink)
		protected.GET("/links", linkHandler.GetUserLinks)
		protected.PUT("/links/:id", linkHandler.UpdateLink)
		protected.DELETE("/links/:id", linkHandler.DeleteLink)
		protected.PATCH("/links/:id/status", linkHandler.UpdateLinkStatus)

		protected.PUT("/profile", profileHandler.UpdateProfile)
		protected.PUT("/profile/password", profileHandler.ChangePassword)
		protected.DELETE("/profile", profileHandler.DeleteAccount)
	}

	r.GET(
		"/u/:username/:slug",
		middleware.RateLimit(config.RedirectRate, config.RedirectBurst),
		linkHandler.Redirect,
	)

	r.GET("/u/:username", profileHandler.GetPublicProfile)

	analytics := protected.Group("/analytics")
	{
		analytics.GET("/overview", analyticsHandler.GetOverview)
		analytics.GET("/links", analyticsHandler.GetLinkAnalytics)
		analytics.GET("/daily", analyticsHandler.GetDailyAnalytics)
		analytics.GET("/recent", analyticsHandler.GetRecentActivity)
		analytics.GET("/referrers", analyticsHandler.GetReferrerAnalytics)
	}
}
