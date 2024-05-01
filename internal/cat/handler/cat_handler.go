package cat

import (
	dto "1-cat-social/internal/cat/dto"
	uc "1-cat-social/internal/cat/usecase"
	validate "1-cat-social/internal/cat/validate"
	"1-cat-social/pkg/logger"
	"1-cat-social/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CatHandler struct {
	uc uc.ICatUsecase
}

func NewCatHandler(uc uc.ICatUsecase) *CatHandler {
	return &CatHandler{
		uc: uc,
	}
}

func (h *CatHandler) Router(r *gin.RouterGroup) {
	endpoint := r.Group("/cats")

	endpoint.GET("", h.GetAll)
	endpoint.PUT("/:id", h.Update)
}

func (h *CatHandler) GetAll(c *gin.Context) {
	var queryParam dto.CatRequestQueryParams
	if err := c.ShouldBindQuery(&queryParam); err != nil {
		logger.Info(err.Error())
		response.GenerateResponse(c, http.StatusBadRequest, response.WithMessage(err.Error()))
		c.Abort()
		return
	}

	cats, err := h.uc.GetAll(&queryParam)
	if err != nil {
		logger.Error(err)
		response.GenerateResponse(c, err.Code, response.WithMessage(err.Message), response.WithData(err.Err))
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
