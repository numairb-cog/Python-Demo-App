package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	config := LoadConfig()

	if err := InitDB(config); err != nil {
		log.Printf("Warning: Database initialization failed: %v", err)
		log.Println("Application will continue but database routes may not work")
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(gin.Recovery())

	router.GET("/", indexHandler)
	router.GET("/wave/:id", waveHandler)
	router.GET("/error/:when", errorHandler)
	router.GET("/query/:dbtype", queryHandler)
	router.GET("/http", httpHandler)

	addr := fmt.Sprintf("0.0.0.0:%s", config.Port)
	log.Printf("Starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
