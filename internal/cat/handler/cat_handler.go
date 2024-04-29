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
		c.JSON(400, response.GenerateResponse("error", nil))
		c.Abort()
		return
	}

	errr := validate.ValidateUpdateCatForm(request)
	if errr != nil {
		c.JSON(400, response.GenerateResponse("error", nil))
		c.Abort()
		return
	}

	cat, err := h.uc.Update(id, request)
	if err != nil {
		c.JSON(err.Code, response.GenerateResponse("error", nil))
		c.Abort()
		return
	}

	c.JSON(200, response.GenerateResponse("success", cat))
}
