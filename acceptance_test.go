package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

const testPort = "9001"

func setupTestServer() *http.Server {
	gin.SetMode(gin.TestMode)

	config := LoadConfig()
	InitDB(config)

	router := gin.Default()
	router.GET("/", indexHandler)
	router.GET("/wave/:id", waveHandler)
	router.GET("/error/:when", errorHandler)
	router.GET("/query/:dbtype", queryHandler)
	router.GET("/http", httpHandler)

	srv := &http.Server{
		Addr:    ":" + testPort,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	return srv
}

func TestAcceptanceIndexPage(t *testing.T) {
	srv := setupTestServer()
	defer srv.Close()

	resp, err := http.Get("http://localhost:" + testPort + "/")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Go Gin Demo App") {
		t.Error("Index page should contain 'Go Gin Demo App'")
	}

	if !strings.Contains(bodyStr, "/wave/abc") {
		t.Error("Index page should contain link to /wave/abc")
	}

	if !strings.Contains(bodyStr, "/error/always") {
		t.Error("Index page should contain link to /error/always")
	}
}

func TestAcceptanceWaveEndpoint(t *testing.T) {
	srv := setupTestServer()
	defer srv.Close()

	resp, err := http.Get("http://localhost:" + testPort + "/wave/123")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "<!DOCTYPE html>") {
		t.Error("Wave page should contain HTML doctype")
	}

	if !strings.Contains(bodyStr, "<title>Wave</title>") {
		t.Error("Wave page should contain Wave title")
	}
}

func TestAcceptanceErrorEndpoint(t *testing.T) {
	srv := setupTestServer()
	defer srv.Close()

	resp, err := http.Get("http://localhost:" + testPort + "/error/sometimes")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", resp.StatusCode)
	}
}

func TestAcceptanceQueryInvalidType(t *testing.T) {
	srv := setupTestServer()
	defer srv.Close()

	resp, err := http.Get("http://localhost:" + testPort + "/query/mysql")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Only 'pgsql' is supported") {
		t.Error("Should return error message about pgsql only")
	}
}

func TestAcceptanceHTTPMissingURL(t *testing.T) {
	srv := setupTestServer()
	defer srv.Close()

	resp, err := http.Get("http://localhost:" + testPort + "/http")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "Missing required argument: url") {
		t.Error("Should return error message about missing URL")
	}
}

func TestAcceptanceHTTPInvalidURL(t *testing.T) {
	srv := setupTestServer()
	defer srv.Close()

	resp, err := http.Get("http://localhost:" + testPort + "/http?url=notaurl")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "url must be a URL with protocol") {
		t.Error("Should return error message about URL format")
	}
}
