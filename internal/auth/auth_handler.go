package auth

import (
	"1-cat-social/internal/user"
	"1-cat-social/pkg/response"
	"1-cat-social/pkg/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Auth handler structure for auth
type authHandler struct {
	uc IAuthUsecase
}

// Constructor for auth handler struct
func NewAuthHandler(uc IAuthUsecase) *authHandler {
	return &authHandler{
		uc: uc,
	}
}

// Router is required to wrap all user request by spesific path URL
func (h *authHandler) Router(r *gin.RouterGroup) {
	// Grouping to give URL prefix
	// ex : localhost/user
	group := r.Group("user")

	// Utillize group to use global setting on group parent (if exists)
	group.POST("login", h.login)
	group.POST("register", h.register)
}

func (h *authHandler) login(ctx *gin.Context) {
	var request user.LoginDTO

	// Parse request body to DTO
	// If error return error response
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.GenerateResponse(ctx, 400)
		ctx.Abort()
		return
	}

	// Validate request
	validate := validator.New(validator.WithRequiredStructEnabled())

	// Generate error validation if not any field is not valid
	if err := validate.Struct(request); err != nil {
		validatorMessage := validation.GenerateStructValidationError(err)

		response.GenerateResponse(ctx, http.StatusBadRequest, response.WithMessage("Any input is not valid"), response.WithData(validatorMessage))
		ctx.Abort()
		return
	}

	// Process login on usecase
	result, err := h.uc.Login(request)
	if err != nil {
		response.GenerateResponse(ctx, err.Code, response.WithMessage(err.Message))
		ctx.Abort()
		return
	}

	response.GenerateResponse(ctx, 200, response.WithMessage("User logged successfully"), response.WithData(result))
}

func (h *authHandler) register(ctx *gin.Context) {
	var request user.UserDTO

	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.GenerateResponse(ctx, 400)
		ctx.Abort()
		return
	}

	// Validate request
	validate := validator.New(validator.WithRequiredStructEnabled())

	// Generate error validation if not any field is not valid
	if err := validate.Struct(request); err != nil {
		validatorMessage := validation.GenerateStructValidationError(err)

		response.GenerateResponse(ctx, http.StatusBadRequest, response.WithMessage("Any input is not valid"), response.WithData(validatorMessage))
		ctx.Abort()
		return
	}

	// Process register via usecase
	result, err := h.uc.Register(request)
	if err != nil {
		response.GenerateResponse(ctx, err.Code, response.WithMessage(err.Message))
		ctx.Abort()
		return
	}

	response.GenerateResponse(ctx, 201, response.WithMessage("User registered successfully"), response.WithData(result))
}
