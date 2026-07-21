package handler_test

import (
	"bytes"
	"context"
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

func TestLinkHandler_CreateLink(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockLinkService{}
	mockUserRepo := &MockUserRepository{}
	h := handler.NewLinkHandler(mockSvc, mockUserRepo)

	userID := uuid.New()
	mockSvc.CreateLinkFunc = func(uid uuid.UUID, req models.CreateLinkRequest) error {
		assert.Equal(t, userID, uid)
		return nil
	}

	w := httptest.NewRecorder()
	reqBody, _ := json.Marshal(models.CreateLinkRequest{
		Title:          "Git",
		Slug:           "github",
		DestinationURL: "https://github.com",
	})
	req := httptest.NewRequest("POST", "/links", bytes.NewBuffer(reqBody))

	_, r := gin.CreateTestContext(w)
	r.POST("/links", func(c *gin.Context) {
		c.Set("user_id", userID)
	}, h.CreateLink)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestLinkHandler_GetUserLinks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockLinkService{}
	mockUserRepo := &MockUserRepository{}
	h := handler.NewLinkHandler(mockSvc, mockUserRepo)

	userID := uuid.New()
	mockSvc.GetUserLinksFunc = func(uid uuid.UUID) ([]*models.Link, error) {
		return []*models.Link{{Title: "Test Link"}}, nil
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/links", nil)

	_, r := gin.CreateTestContext(w)
	r.GET("/links", func(c *gin.Context) {
		c.Set("user_id", userID)
	}, h.GetUserLinks)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLinkHandler_Redirect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockLinkService{}
	userID := uuid.New()
	mockUserRepo := &MockUserRepository{
		GetByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
			assert.Equal(t, "testuser", username)
			return &models.User{
				ID:       userID,
				Username: username,
			}, nil
		},
	}
	h := handler.NewLinkHandler(mockSvc, mockUserRepo)

	mockSvc.HandleRedirectFunc = func(uid uuid.UUID, slug, ip, ua, ref string) (string, error) {
		assert.Equal(t, userID, uid)
		assert.Equal(t, "slug", slug)
		return "https://destination.com", nil
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/u/testuser/slug", nil)

	_, r := gin.CreateTestContext(w)
	r.GET("/u/:username/:slug", h.Redirect)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "https://destination.com", w.Header().Get("Location"))
}

func TestLinkHandler_UpdateLink(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockLinkService{}
	mockUserRepo := &MockUserRepository{}
	h := handler.NewLinkHandler(mockSvc, mockUserRepo)

	userID := uuid.New()
	linkID := uuid.New()

	mockSvc.UpdateLinkFunc = func(uid, lid uuid.UUID, req models.UpdateLinkRequest) (*models.Link, error) {
		return &models.Link{ID: lid, UserID: uid, Title: req.Title}, nil
	}

	w := httptest.NewRecorder()
	reqBody, _ := json.Marshal(models.UpdateLinkRequest{
		Title:          "New Title",
		Slug:           "newslug",
		DestinationURL: "https://newdestination.com",
	})
	req := httptest.NewRequest("PUT", "/links/"+linkID.String(), bytes.NewBuffer(reqBody))

	_, r := gin.CreateTestContext(w)
	r.PUT("/links/:id", func(c *gin.Context) {
		c.Set("user_id", userID)
	}, h.UpdateLink)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLinkHandler_DeleteLink(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockLinkService{}
	mockUserRepo := &MockUserRepository{}
	h := handler.NewLinkHandler(mockSvc, mockUserRepo)

	userID := uuid.New()
	linkID := uuid.New()

	mockSvc.DeleteLinkFunc = func(uid, lid uuid.UUID) error {
		return nil
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/links/"+linkID.String(), nil)

	_, r := gin.CreateTestContext(w)
	r.DELETE("/links/:id", func(c *gin.Context) {
		c.Set("user_id", userID)
	}, h.DeleteLink)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLinkHandler_UpdateLinkStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockLinkService{}
	mockUserRepo := &MockUserRepository{}
	h := handler.NewLinkHandler(mockSvc, mockUserRepo)

	userID := uuid.New()
	linkID := uuid.New()

	mockSvc.UpdateLinkStatusFunc = func(uid, lid uuid.UUID, active bool) error {
		return nil
	}

	w := httptest.NewRecorder()
	reqBody, _ := json.Marshal(models.UpdateLinkStatusRequest{IsActive: false})
	req := httptest.NewRequest("PATCH", "/links/"+linkID.String()+"/status", bytes.NewBuffer(reqBody))

	_, r := gin.CreateTestContext(w)
	r.PATCH("/links/:id/status", func(c *gin.Context) {
		c.Set("user_id", userID)
	}, h.UpdateLinkStatus)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
