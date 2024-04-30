package main

import (
	"1-cat-social/config"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load ENV from OS env or from .env variable
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db := config.InitDb()

	fmt.Println(db.Ping())

	r := gin.Default()

	// Grouping the routes and give prefix for the API
	api := r.Group("api/v1")
	// Route ping targetted ping Handler
	// (using gin) Handler is a function that has gin context param
	api.GET("ping", pingHandler)

	// Handle no route
	r.NoRoute(NoRouteHandler)
	NewRoute(db, api)

	// Start the server
	r.Run("0.0.0.0:8080")
}

// Handler for ping request from routes
func pingHandler(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		struct {
			Data    any    `json:"data"`
			Message string `json:"message"`
			Success bool   `json:"success"`
		}{
			Success: true,
			Message: "Server is online",
			Data:    true,
		},
	)
}
