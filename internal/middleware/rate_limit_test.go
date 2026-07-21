package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MRNaveed-stack/LinkPulse/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestRateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w1 := httptest.NewRecorder()
	c1, r := gin.CreateTestContext(w1)

	// Rate limit: 1 request/sec, burst size 1
	r.Use(middleware.RateLimit(rate.Limit(1), 1))
	r.GET("/limit", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	req1 := httptest.NewRequest("GET", "/limit", nil)
	req1.RemoteAddr = "192.168.1.100:12345"
	c1.Request = req1
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Second request immediately after (should be blocked)
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/limit", nil)
	req2.RemoteAddr = "192.168.1.100:12345"
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusTooManyRequests, w2.Code)
}
