package cat

import (
	dto "1-cat-social/internal/cat/dto"
	uc "1-cat-social/internal/cat/usecase"
	validate "1-cat-social/internal/cat/validate"
	"1-cat-social/pkg/response"

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

	endpoint.PUT("/:id", h.Update)
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

	response.GenerateResponse(c, 200, response.WithMessage("Success"), response.WithData(modifiedCat))
}
