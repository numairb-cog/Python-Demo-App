package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

const testPort = "9001"

var testBaseURL string

func TestMain(m *testing.M) {
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
	testBaseURL = "http://localhost:" + testPort

	code := m.Run()

	srv.Close()
	os.Exit(code)
}

func makeRequest(t *testing.T, path string) *http.Response {
	resp, err := http.Get(testBaseURL + path)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	return resp
}

func readBody(t *testing.T, resp *http.Response) string {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}
	return string(body)
}

func TestAcceptanceIndexPage(t *testing.T) {
	resp := makeRequest(t, "/")
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body := readBody(t, resp)
	checks := []string{"Go Gin Demo App", "/wave/abc", "/error/always"}
	for _, check := range checks {
		if !strings.Contains(body, check) {
			t.Errorf("Index page should contain '%s'", check)
		}
	}
}

func TestAcceptanceWaveEndpoint(t *testing.T) {
	resp := makeRequest(t, "/wave/123")
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body := readBody(t, resp)
	checks := []string{"<!DOCTYPE html>", "<title>Wave</title>"}
	for _, check := range checks {
		if !strings.Contains(body, check) {
			t.Errorf("Wave page should contain '%s'", check)
		}
	}
}

func TestAcceptanceErrorEndpoint(t *testing.T) {
	resp := makeRequest(t, "/error/sometimes")
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", resp.StatusCode)
	}
	readBody(t, resp)
}

func TestAcceptanceQueryInvalidType(t *testing.T) {
	resp := makeRequest(t, "/query/mysql")
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	body := readBody(t, resp)
	if !strings.Contains(body, "Only 'pgsql' is supported") {
		t.Error("Should return error message about pgsql only")
	}
}

func TestAcceptanceHTTPMissingURL(t *testing.T) {
	resp := makeRequest(t, "/http")
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	body := readBody(t, resp)
	if !strings.Contains(body, "Missing required argument: url") {
		t.Error("Should return error message about missing URL")
	}
}

func TestAcceptanceHTTPInvalidURL(t *testing.T) {
	resp := makeRequest(t, "/http?url=notaurl")
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}

	body := readBody(t, resp)
	if !strings.Contains(body, "url must be a URL with protocol") {
		t.Error("Should return error message about URL format")
	}
}
