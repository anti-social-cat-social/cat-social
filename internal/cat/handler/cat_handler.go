package cat

import (
	dto "1-cat-social/internal/cat/dto"
	uc "1-cat-social/internal/cat/usecase"
	validate "1-cat-social/internal/cat/validate"
	"1-cat-social/pkg/logger"
	"1-cat-social/pkg/response"
	"1-cat-social/pkg/validator"
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
	endpoint := r.Group("/cat")

	endpoint.GET("", h.GetAll)
	endpoint.PUT("/:id", h.Update)
	endpoint.POST("/match", h.Match)
}

func (h *CatHandler) GetAll(c *gin.Context) {
	var queryParam dto.CatRequestQueryParams
	if err := c.ShouldBindQuery(&queryParam); err != nil {
		logger.Info(err.Error())
		c.JSON(http.StatusBadRequest, response.GenerateResponse(err.Error(), nil))
		c.Abort()
		return
	}

	cats, err := h.uc.GetAll(&queryParam)
	if err != nil {
		logger.Error(err)
		c.JSON(err.Code, response.GenerateResponse(err.Err, nil))
		c.Abort()
		return
	}

	catResponse := dto.FormatCatsResponse(cats)

	c.JSON(http.StatusOK, response.GenerateResponse("success", catResponse))
}

func (h *CatHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var request dto.CatUpdateRequestBody
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, response.GenerateResponse(err.Error(), nil))
		c.Abort()
		return
	}

	errr := validate.ValidateUpdateCatForm(request)
	if errr != nil {
		c.JSON(400, response.GenerateResponse(errr.Error(), nil))
		c.Abort()
		return
	}

	cat, err := h.uc.Update(id, request)
	if err != nil {
		c.JSON(err.Code, response.GenerateResponse(err.Err, nil))
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

	c.JSON(200, response.GenerateResponse("success", modifiedCat))
}

func (h *CatHandler) Match(c *gin.Context) {
	var request dto.CatMatchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		validation := validator.FormatValidation(err)
		logger.Info(validation)
		c.JSON(http.StatusBadRequest, response.GenerateResponse(validation, nil))
		return
	}

	c.JSON(http.StatusCreated, response.GenerateResponse("success", nil))
}
