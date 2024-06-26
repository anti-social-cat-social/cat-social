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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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

func (h *CatHandler) Router(r *gin.RouterGroup, db *sqlx.DB) {
	endpoint := r.Group("/cat")
	endpoint.Use(middleware.UseJwtAuth)

	endpoint.GET("", h.GetAll)
	endpoint.POST("", h.Store)
	endpoint.PUT("/:id", h.Update)
	endpoint.DELETE("/:id", h.Delete)
	endpoint.GET("/match", h.GetMatches)
	endpoint.POST("/match", middleware.DBTransactionMiddleware(db), h.Match)
	endpoint.POST("/match/approve", middleware.DBTransactionMiddleware(db), h.ApproveMatch)
	endpoint.POST("/match/reject", middleware.DBTransactionMiddleware(db), h.RejectMatch)
	endpoint.DELETE("/match/:id", h.DeleteMatch)
}

func (h *CatHandler) GetAll(c *gin.Context) {
	var queryParam dto.CatRequestQueryParams
	if err := c.ShouldBindQuery(&queryParam); err != nil {
		validation := validator.FormatValidation(err)
		logger.Info(validation)
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

func (h *CatHandler) Store(c *gin.Context) {
	var request dto.CatUpdateRequestBody

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Info(err.Error())
		response.GenerateResponse(c, http.StatusBadRequest, response.WithMessage(err.Error()))
		c.Abort()
		return
	}

	errValidate := validate.ValidateUpdateCatForm(request)
	if errValidate != nil {
		response.GenerateResponse(c, 400, response.WithMessage(errValidate.Error()))
		c.Abort()
		return
	}

	userID := c.MustGet("userID").(string)

	cat, err := h.uc.Create(request, userID)
	if err != nil {
		logger.Info(err.Message)
		response.GenerateResponse(c, http.StatusBadRequest, response.WithMessage(err.Message))
		c.Abort()
		return
	}

	catResponse := dto.CatStoreResponse{
		ID:        cat.ID,
		CreatedAt: cat.CreatedAt.Format(time.RFC3339),
	}

	response.GenerateResponse(c, http.StatusCreated, response.WithMessage("success"), response.WithData(catResponse))
}

func (h *CatHandler) Update(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("userID").(string)

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

	cat, err := h.uc.Update(id, request, userID)
	if err != nil {
		response.GenerateResponse(c, err.Code, response.WithMessage(err.Err))
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

func (h *CatHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	userId := c.MustGet("userID").(string)

	if err := h.uc.Delete(id, userId); err != nil {
		response.GenerateResponse(c, err.Code, response.WithMessage(err.Message))
		c.Abort()
		return
	}

	response.GenerateResponse(c, http.StatusOK, response.WithMessage("success"))
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
	txHandle := c.MustGet("db_trx").(*sqlx.Tx)

	err := h.mc.WithTrx(txHandle).Match(&request, userID)
	h.mc.RemoveTrx()
	if err != nil {
		err.Message = err.Message + " >>> error match"
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

func (h *CatHandler) ApproveMatch(c *gin.Context) {
	var request dto.MatchApproveRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		validation := validator.FormatValidation(err)
		logger.Info(validation)
		response.GenerateResponse(c, http.StatusBadRequest, response.WithMessage(validation))
		return
	}

	userID := c.MustGet("userID").(string)
	txHandle := c.MustGet("db_trx").(*sqlx.Tx)

	err := h.mc.WithTrx(txHandle).Approve(&request, userID)
	h.mc.RemoveTrx()
	if err != nil {
		err.Message = err.Message + " >>> error approve"
		if err.Code == http.StatusInternalServerError {
			logger.Error(err)
		} else {
			logger.Info(err.Err)
		}
		response.GenerateResponse(c, err.Code, response.WithMessage(err.Err))
		return
	}

	response.GenerateResponse(c, http.StatusOK, response.WithMessage("success"))
}

func (h *CatHandler) RejectMatch(c *gin.Context) {
	var request dto.MatchApproveRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		validation := validator.FormatValidation(err)
		logger.Info(validation)
		response.GenerateResponse(c, http.StatusBadRequest, response.WithMessage(validation))
		return
	}

	userID := c.MustGet("userID").(string)
	txHandle := c.MustGet("db_trx").(*sqlx.Tx)

	err := h.mc.WithTrx(txHandle).Reject(&request, userID)
	h.mc.RemoveTrx()
	if err != nil {
		err.Message = err.Message + " >>> error reject"
		if err.Code == http.StatusInternalServerError {
			logger.Error(err)
		} else {
			logger.Info(err.Err)
		}
		response.GenerateResponse(c, err.Code, response.WithMessage(err.Err))
		return
	}

	response.GenerateResponse(c, http.StatusOK, response.WithMessage("success"))
}

func (h *CatHandler) GetMatches(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	matches, err := h.mc.GetMatches(userID)
	if err != nil {
		if err.Code == http.StatusInternalServerError {
			logger.Error(err)
		} else {
			logger.Info(err.Err)
		}
		response.GenerateResponse(c, err.Code, response.WithMessage(err.Err))
		return
	}

	response.GenerateResponse(c, http.StatusOK, response.WithMessage("success"), response.WithData(matches))
}

func (h *CatHandler) DeleteMatch(c *gin.Context) {
	matchID := c.Param("id")
	userID := c.MustGet("userID").(string)

	err := h.mc.DeleteMatch(matchID, userID)
	if err != nil {
		response.GenerateResponse(c, err.Code, response.WithMessage(err.Err))
		return
	}

	response.GenerateResponse(c, http.StatusOK, response.WithMessage("success"))
}
