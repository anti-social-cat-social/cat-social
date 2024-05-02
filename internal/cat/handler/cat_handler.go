package cat

import (
	dto "1-cat-social/internal/cat/dto"
	uc "1-cat-social/internal/cat/usecase"
	validate "1-cat-social/internal/cat/validate"
	"1-cat-social/internal/middleware"
	"1-cat-social/pkg/logger"
	"1-cat-social/pkg/response"
	"1-cat-social/pkg/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CatHandler struct {
	uc uc.ICatUsecase
	mc uc.IMatchUsecase
}

func NewCatHandler(uc uc.ICatUsecase, mc uc.IMatchUsecase) *CatHandler {
	return &CatHandler{
		uc: uc,
		mc: mc,
	}
}

func (h *CatHandler) Router(r *gin.RouterGroup) {
	endpoint := r.Group("/cat")
	endpoint.Use(middleware.UseJwtAuth)

	endpoint.GET("", h.GetAll)
	endpoint.PUT("/:id", h.Update)
	endpoint.POST("/match", h.Match)
}

func (h *CatHandler) GetAll(c *gin.Context) {
	var queryParam dto.CatRequestQueryParams
	if err := c.ShouldBindQuery(&queryParam); err != nil {
		logger.Info(err.Error())
		response.GenerateResponse(c, http.StatusBadRequest, response.WithMessage(err.Error()))
		c.Abort()
		return
	}

	userID := c.MustGet("userID").(string)

	cats, err := h.uc.GetAll(&queryParam, userID)
	if err != nil {
		if err.Code == http.StatusInternalServerError {
			logger.Error(err)
		} else {
			logger.Info(err.Err)
		}
		response.GenerateResponse(c, err.Code, response.WithMessage(err.Err))
		c.Abort()
		return
	}

	catResponse := dto.FormatCatsResponse(cats)

	response.GenerateResponse(c, http.StatusOK, response.WithMessage("success"), response.WithData(catResponse))
}

func (h *CatHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var request dto.CatUpdateRequestBody
	if err := c.ShouldBindJSON(&request); err != nil {
		response.GenerateResponse(c, 400, response.WithMessage(err.Error()))
		c.Abort()
		return
	}

	errr := validate.ValidateUpdateCatForm(request)
	if errr != nil {
		response.GenerateResponse(c, 400, response.WithMessage(errr.Error()))
		c.Abort()
		return
	}

	cat, err := h.uc.Update(id, request)
	if err != nil {
		response.GenerateResponse(c, err.Code, response.WithMessage(err.Message))
		c.Abort()
		return
	}

	modifiedCat := dto.CatUpdateResponseBody{
		ID:          cat.ID,
		Name:        cat.Name,
		Race:        cat.Race,
		Sex:         cat.Sex,
		AgeInMonth:  cat.AgeInMonth,
		Description: cat.Description,
		ImageUrls:   cat.ImageUrls,
	}

	response.GenerateResponse(c, http.StatusOK, response.WithMessage("Success"), response.WithData(modifiedCat))
}

func (h *CatHandler) Match(c *gin.Context) {
	var request dto.CatMatchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		validation := validator.FormatValidation(err)
		logger.Info(validation)
		response.GenerateResponse(c, http.StatusBadRequest, response.WithMessage(validation))
		return
	}

	userID := c.MustGet("userID").(string)
	err := h.mc.Match(&request, userID)
	if err != nil {
		if err.Code == http.StatusInternalServerError {
			logger.Error(err)
		} else {
			logger.Info(err.Err)
		}
		response.GenerateResponse(c, err.Code, response.WithMessage(err.Err))
		return
	}

	response.GenerateResponse(c, http.StatusCreated, response.WithMessage("success"))
}
