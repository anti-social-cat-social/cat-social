package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load ENV from OS env or from .env variable
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(os.Getenv("DB_PARAMS"))

	r := gin.Default()

	// Grouping the routes and give prefix for the API
	api := r.Group("api/v1")
	// Route ping targetted ping Handler
	// (using gin) Handler is a function that has gin context param
	api.GET("ping", pingHandler)

	// Initialize other routes
	NewRoute(api)

	// Start the server
	r.Run("0.0.0.0:8080")
}

// Handler for ping request from routes
func pingHandler(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
			Data    any    `json:"data"`
		}{
			Success: true,
			Message: "Server is online",
			Data:    true,
		},
	)
}
