package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/MRNaveed-stack/LinkPulse/internal/logger"
	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/MRNaveed-stack/LinkPulse/internal/repository"
	"github.com/MRNaveed-stack/LinkPulse/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LinkHandler struct {
	linkService service.LinkService
	userRepo    repository.UserRepository
}

func NewLinkHandler(linkService service.LinkService, userRepo repository.UserRepository) *LinkHandler {
	return &LinkHandler{
		linkService: linkService,
		userRepo:    userRepo,
	}
}

func (h *LinkHandler) CreateLink(c *gin.Context) {
	var req models.CreateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn(
			"validation failed",
			"endpoint", c.FullPath(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	err := h.linkService.CreateLink(userID, req)
	if err != nil {
		log.Printf("CreateLink: error creating link for user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create link"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Link created successfully",
	})
}

func (h *LinkHandler) GetUserLinks(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	links, err := h.linkService.GetUserLinks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve links"})
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{"links": links},
	)
}

func (h *LinkHandler) Redirect(c *gin.Context) {
	username := c.Param("username")
	slug := c.Param("slug")
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	referrer := c.Request.Referer()

	// Get user by username
	user, err := h.userRepo.GetByUsername(c.Request.Context(), username)
	if err != nil {
		logger.Log.Warn("failed to get user by username", "username", username, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	url, err := h.linkService.HandleRedirect(user.ID, slug, ip, userAgent, referrer)
	if err != nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "link not found",
			},
		)
		return
	}
	c.Redirect(
		http.StatusFound,
		url,
	)
}

func (h *LinkHandler) UpdateLink(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	linkIDstr := c.Param("id")
	linkID, err := uuid.Parse(linkIDstr)
	if err != nil {
		logger.Log.Warn(
			"validation failed",
			"endpoint", c.FullPath(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid link ID format"})
		return
	}
	var req models.UpdateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn(
			"validation failed",
			"endpoint", c.FullPath(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request body: " + err.Error()})
		return
	}

	if req.Title == "" || req.Slug == "" || req.DestinationURL == "" {
		err := errors.New("All fields are required")
		logger.Log.Warn(
			"validation failed",
			"endpoint", c.FullPath(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	updatedLink, err := h.linkService.UpdateLink(userID.(uuid.UUID), linkID, req)
	if err != nil {
		// Handle specific errors
		if err.Error() == "unauthorized: you don't own this link" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "you don't have permission to update this link",
			})
			return
		}

		if err.Error() == "link not found" ||
			err.Error() == "link not found: sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "link not found",
			})
			return
		}

		// Generic server error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update link: " + err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "link updated successfully",
		"link":    updatedLink,
	})
}

func (h *LinkHandler) DeleteLink(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	linkIDstr := c.Param("id")
	linkID, err := uuid.Parse(linkIDstr)
	if err != nil {
		logger.Log.Warn(
			"validation failed",
			"endpoint", c.FullPath(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid link ID format"})
		return
	}

	err = h.linkService.DeleteLink(userID.(uuid.UUID), linkID)
	if err != nil {
		if err.Error() == "unauthorized: you don't own this link" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "you don't have permission to delete this link",
			})
			return
		}
		if err.Error() == "link not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "link not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "link deleted successfully"})
}

func (h *LinkHandler) UpdateLinkStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	linkIDstr := c.Param("id")
	linkID, err := uuid.Parse(linkIDstr)
	if err != nil {
		logger.Log.Warn(
			"validation failed",
			"endpoint", c.FullPath(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid link ID format"})
		return
	}
	var req models.UpdateLinkStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn(
			"validation failed",
			"endpoint", c.FullPath(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request body: " + err.Error()})
		return
	}

	err = h.linkService.UpdateLinkStatus(userID.(uuid.UUID), linkID, req.IsActive)
	if err != nil {
		if err.Error() == "unauthorized: you don't own this link" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "you don't have permission to update this link",
			})
			return
		}
		if err.Error() == "link not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "link not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	statusText := "disabled"
	if req.IsActive {
		statusText = "enabled"
	}
	c.JSON(http.StatusOK, gin.H{
		"message":   fmt.Sprintf("link %s successfully", statusText),
		"is_active": req.IsActive,
	})
}
