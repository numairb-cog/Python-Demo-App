package main

import (
	"fmt"
	"html/template"
	"io"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	contentTypeHTML = "text/html; charset=utf-8"
)

func indexHandler(c *gin.Context) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading template: %v", err)
		return
	}

	data := map[string]interface{}{
		"WaveURL":        "/wave/abc",
		"ErrorAlways":    "/error/always",
		"ErrorSometimes": "/error/sometimes",
		"QueryPgSQL":     "/query/pgsql",
		"HTTPExitCall":   "/http",
	}

	c.Header("Content-Type", contentTypeHTML)
	if err := tmpl.Execute(c.Writer, data); err != nil {
		c.String(http.StatusInternalServerError, "Error rendering template: %v", err)
	}
}

func waveHandler(c *gin.Context) {
	whatever := c.Param("id")

	a := 0.2
	b := 3.0
	cv := 1.0
	x := float64(time.Now().Unix()) * math.Pi / 180.0

	if val, err := strconv.Atoi(whatever); err == nil {
		cv = float64(val)
	}

	delay := a * (math.Sin(b*x+cv) + 1.0)
	time.Sleep(time.Duration(delay * float64(time.Second)))

	html := fmt.Sprintf("<!DOCTYPE html><title>Wave</title><h1>%f</h1>", delay)
	c.Header("Content-Type", contentTypeHTML)
	c.String(http.StatusOK, html)
}

func errorHandler(c *gin.Context) {
	when := c.Param("when")

	shouldError := when == "always" || rand.Intn(10) == 0

	if shouldError {
		errorType := rand.Intn(4)
		switch errorType {
		case 0:
			panic("KeyError: 'typo_eror'")
		case 1:
			panic("IndexError: list index out of range")
		case 2:
			panic("AssertionError: True is False")
		default:
			panic("ValueError: invalid literal for int() with base 10: 'abc'")
		}
	}

	html := "<!DOCTYPE html><title>No Exception This Time</title><h1>OK This Time</h1>"
	c.Header("Content-Type", contentTypeHTML)
	c.String(http.StatusOK, html)
}

func queryHandler(c *gin.Context) {
	dbtype := c.Param("dbtype")

	if dbtype != "pgsql" {
		c.String(http.StatusBadRequest, "Only 'pgsql' is supported")
		return
	}

	queryTypes := []string{"slow", "error", "normal", "normal", "normal", "normal", "normal", "normal", "normal", "normal"}
	queryType := queryTypes[rand.Intn(len(queryTypes))]
	sleepTime := float64(rand.Intn(7)+1) / 10.0

	database := GetDB()
	if database == nil {
		c.String(http.StatusInternalServerError, "Database not initialized")
		return
	}

	var query string
	switch queryType {
	case "slow":
		query = fmt.Sprintf("SELECT pg_sleep(%f)", sleepTime)
	case "error":
		query = "SELECT sql - error"
	default:
		query = "SELECT 123"
	}

	_, err := database.Exec(query)
	if err != nil && queryType != "error" {
		c.String(http.StatusInternalServerError, "Database error: %v", err)
		return
	}

	html := fmt.Sprintf("<!DOCTYPE html><title>Query %s</title><h1>Ran DB Query %s</h1>", dbtype, queryType)
	c.Header("Content-Type", contentTypeHTML)
	c.String(http.StatusOK, html)
}

func httpHandler(c *gin.Context) {
	url := c.Query("url")

	if url == "" {
		html := "<!DOCTYPE html><title>Missing Required Argument</title><h1>Missing required argument: url</h1>"
		c.Header("Content-Type", contentTypeHTML)
		c.String(http.StatusBadRequest, html)
		return
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		html := "<!DOCTYPE html><title>Missing Required Argument</title><h1>Missing required argument: url must be a URL with protocol, like http://...</h1>"
		c.Header("Content-Type", contentTypeHTML)
		c.String(http.StatusBadRequest, html)
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		body, _ := io.ReadAll(resp.Body)
		contentLength = strconv.Itoa(len(body))
	}

	html := fmt.Sprintf("<!DOCTYPE html><title>HTTP Exit Call</title><h1>Response from %s</h1><p>Content length %s</p>", url, contentLength)
	c.Header("Content-Type", contentTypeHTML)
	c.String(http.StatusOK, html)
}
