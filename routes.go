package main

import (
	"1-cat-social/internal/auth"
	"1-cat-social/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Initialize all router from all handler
func NewRoute(db *sqlx.DB, router *gin.RouterGroup) {
	initializeAuthHandler(db, router)
}

func initializeAuthHandler(db *sqlx.DB, router *gin.RouterGroup) {
	// Initialize all ncessary dependecies
	userRepo := user.NewUserRepository(db)
	userUc := user.NewUserUsecase(userRepo)
	authUc := auth.NewAuthUsecase(userUc)
	authH := auth.NewAuthHandler(authUc)

	// Do not forget
	// Call auth router inside the handler
	authH.Router(router)
}
