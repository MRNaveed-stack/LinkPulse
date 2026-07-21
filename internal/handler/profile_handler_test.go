package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MRNaveed-stack/LinkPulse/internal/handler"
	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProfileHandler_GetPublicProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockProfileService{}
	h := handler.NewProfileHandler(mockSvc)

	mockSvc.GetPublicProfileFunc = func(username string) (*models.PublicProfileResponse, error) {
		return &models.PublicProfileResponse{Username: username, DisplayName: "Disp"}, nil
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/profile/myuser", nil)

	_, r := gin.CreateTestContext(w)
	r.GET("/profile/:username", h.GetPublicProfile)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProfileHandler_UpdateProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockProfileService{}
	h := handler.NewProfileHandler(mockSvc)

	userID := uuid.New()
	mockSvc.UpdateProfileFunc = func(uid uuid.UUID, req *models.UpdateProfileRequest) (*models.Profile, error) {
		return &models.Profile{UserID: uid, DisplayName: req.DisplayName}, nil
	}

	w := httptest.NewRecorder()
	reqBody, _ := json.Marshal(models.UpdateProfileRequest{DisplayName: "Updated Name"})
	req := httptest.NewRequest("PUT", "/profile", bytes.NewBuffer(reqBody))

	_, r := gin.CreateTestContext(w)
	r.PUT("/profile", func(c *gin.Context) {
		c.Set("user_id", userID)
	}, h.UpdateProfile)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
