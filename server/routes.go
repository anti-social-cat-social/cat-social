package server

import (
	"1-cat-social/internal/auth"
	catHandler "1-cat-social/internal/cat/handler"
	cr "1-cat-social/internal/cat/repository"
	catUseCase "1-cat-social/internal/cat/usecase"
	"1-cat-social/internal/user"
	"1-cat-social/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func NewRoute(engine *gin.Engine, db *sqlx.DB) {
	// Handle for not found routes
	engine.NoRoute(NoRouteHandler)
	router := engine.Group("v1")

	router.GET("ping", pingHandler)

	initializeAuthHandler(db, router)
	initializeCatHandler(router, db)
}

func initializeCatHandler(router *gin.RouterGroup, db *sqlx.DB) {
	catRepository := cr.NewCatRepository(db)
	matchRepository := cr.NewMatchRepository(db)
	catUsecase := catUseCase.NewCatUsecase(catRepository)
	matchUsecase := catUseCase.NewMatchUsecase(catRepository, matchRepository)
	catHandler := catHandler.NewCatHandler(catUsecase, matchUsecase)
	catHandler.Router(router)
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

func NoRouteHandler(ctx *gin.Context) {
	response.GenerateResponse(ctx, http.StatusNotFound, response.WithMessage("Page not found"))
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
