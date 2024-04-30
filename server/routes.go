package server

import (
	catHandler "1-cat-social/internal/cat/handler"
	catRepository "1-cat-social/internal/cat/repository"
	catUseCase "1-cat-social/internal/cat/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func NewRoute(engine *gin.Engine, db *sqlx.DB) {
	router := engine.Group("/v1")

	router.GET("ping", pingHandler)

	initializeCatHandler(router, db)

}

func initializeCatHandler(router *gin.RouterGroup, db *sqlx.DB) {

	catRepository := catRepository.NewCatRepository(db)
	catUsecase := catUseCase.NewCatUsecase(catRepository)
	catHandler := catHandler.NewCatHandler(catUsecase)

	catHandler.Router(router)
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
