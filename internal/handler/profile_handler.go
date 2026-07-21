package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/MRNaveed-stack/LinkPulse/internal/logger"
	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/MRNaveed-stack/LinkPulse/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfileHandler struct {
	service service.ProfileService
}

func NewProfileHandler(service service.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		service: service,
	}
}

func (h *ProfileHandler) GetPublicProfile(c *gin.Context) {
	username := c.Param("username")
	profile, err := h.service.GetPublicProfile(username)
	if err != nil {
		log.Printf("GetPublicProfile: error fetching public profile for username %q: %v", username, err)
		c.JSON(http.StatusNotFound,
			gin.H{"error": "profile not found"},
		)
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn(
			"validation failed",
			"endpoint", c.FullPath(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if req.DisplayName == "" && req.Bio == "" && req.AvatarURL == "" && req.Username == "" && req.Email == "" {
		err := errors.New("at least one field must be provided")
		logger.Log.Warn(
			"validation failed",
			"endpoint", c.FullPath(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one field must be provided"})
		return
	}

	profile, err := h.service.UpdateProfile(userID, &req)
	if err != nil {
		log.Printf("UpdateProfile: error updating profile for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "profile updated successfully",
		"profile": profile,
	})
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,min=8"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn(
			"validation failed",
			"endpoint", c.FullPath(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body or weak password"})
		return
	}

	err := h.service.ChangePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		log.Printf("ChangePassword failed for user %s: %v", userID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

func (h *ProfileHandler) DeleteAccount(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	err := h.service.DeleteAccount(c.Request.Context(), userID)
	if err != nil {
		log.Printf("DeleteAccount failed for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "account deleted successfully"})
}
