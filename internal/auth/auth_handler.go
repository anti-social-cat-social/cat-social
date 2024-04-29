package auth

import (
	"1-cat-social/internal/user"
	"1-cat-social/pkg/response"

	"github.com/gin-gonic/gin"
)

type loginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authHandler struct {
	uc IAuthUsecase
}

func NewAuthHandler(uc IAuthUsecase) *authHandler {
	return &authHandler{
		uc: uc,
	}
}

func (h *authHandler) Router(r *gin.RouterGroup) {
	group := r.Group("user")

	group.POST("login", h.login)
	group.POST("register", h.register)
}

func (h *authHandler) login(ctx *gin.Context) {
	var request loginDTO

	// Parse request body to DTO
	// If error return error response
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.GenerateResponse(ctx, 500)
		ctx.Abort()
		return
	}

	// Process login on usecase
	result, err := h.uc.Login()
	if err != nil {
		response.GenerateResponse(ctx, err.Code, response.WithMessage(err.Message))
		ctx.Abort()
		return
	}

	response.GenerateResponse(ctx, 200, response.WithMessage("Success Login"), response.WithData(result))
}

func (h *authHandler) register(ctx *gin.Context) {
	var request user.UserDTO

	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.GenerateResponse(ctx, 500)
		ctx.Abort()
		return
	}

	response.GenerateResponse(ctx, 200, response.WithMessage("Success Register"))
}
