package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MRNaveed-stack/LinkPulse/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(middleware.Recovery())
	r.GET("/panic", func(ctx *gin.Context) {
		panic("something went catastrophically wrong")
	})

	req := httptest.NewRequest("GET", "/panic", nil)
	c.Request = req
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
