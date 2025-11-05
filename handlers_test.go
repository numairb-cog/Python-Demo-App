package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestIndexHandler(t *testing.T) {
	router := gin.New()
	router.GET("/", indexHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Go Gin Demo App")
	assert.Contains(t, w.Body.String(), "/wave/abc")
	assert.Contains(t, w.Body.String(), "/error/always")
	assert.Contains(t, w.Body.String(), "/query/pgsql")
}

func TestWaveHandler(t *testing.T) {
	router := gin.New()
	router.GET("/wave/:id", waveHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/wave/123", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "<!DOCTYPE html>")
	assert.Contains(t, w.Body.String(), "<title>Wave</title>")
}

func TestWaveHandlerWithNonNumericID(t *testing.T) {
	router := gin.New()
	router.GET("/wave/:id", waveHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/wave/abc", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "<!DOCTYPE html>")
}

func TestErrorHandlerSometimes(t *testing.T) {
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/error/:when", errorHandler)

	successCount := 0
	errorCount := 0

	for i := 0; i < 20; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/error/sometimes", nil)
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			successCount++
			assert.Contains(t, w.Body.String(), "OK This Time")
		} else if w.Code == http.StatusInternalServerError {
			errorCount++
		}
	}

	assert.True(t, successCount > 0, "Should have at least one successful response")
}

func TestErrorHandlerAlways(t *testing.T) {
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/error/:when", errorHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/error/always", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestQueryHandlerInvalidDBType(t *testing.T) {
	router := gin.New()
	router.GET("/query/:dbtype", queryHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/query/mysql", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Only 'pgsql' is supported")
}

func TestQueryHandlerPgSQL(t *testing.T) {
	config := LoadConfig()
	InitDB(config)

	router := gin.New()
	router.GET("/query/:dbtype", queryHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/query/pgsql", nil)
	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		assert.Contains(t, w.Body.String(), "Ran DB Query")
	} else {
		assert.Contains(t, w.Body.String(), "Database error")
	}
}

func TestHTTPHandlerMissingURL(t *testing.T) {
	router := gin.New()
	router.GET("/http", httpHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/http", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing required argument: url")
}

func TestHTTPHandlerInvalidURL(t *testing.T) {
	router := gin.New()
	router.GET("/http", httpHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/http?url=notaurl", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "url must be a URL with protocol")
}

func TestHTTPHandlerValidURL(t *testing.T) {
	router := gin.New()
	router.GET("/http", httpHandler)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Test response"))
	}))
	defer testServer.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/http?url="+testServer.URL, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Response from")
	assert.Contains(t, w.Body.String(), "Content length")
}
