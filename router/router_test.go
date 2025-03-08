package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return OK for root endpoint", func(t *testing.T) {
		r := SetupRouter("", "")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "Hi, GET /incr to check session.", w.Body.String())
	})

	t.Run("should return OK for healthz endpoint", func(t *testing.T) {
		r := SetupRouter("", "")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "\"status\":\"OK\"")
	})

	t.Run("should return OK for readyz endpoint", func(t *testing.T) {
		r := SetupRouter("", "")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/readyz", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "\"status\":\"OK\"")
	})
}
