package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/isyyyy/bnb-pub-sub/internal/models"
	"github.com/isyyyy/bnb-pub-sub/internal/services"
	"net/http"
)


type CandleHandler struct {
	service services.CandleService
}

func NewCandleHandler (candleService *services.CandleService) *CandleHandler {
	return &CandleHandler{service: *candleService}
}
func (h *CandleHandler) SearchCandles(c *gin.Context)  {
	var req models.AggRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, nil)
		return
	}

	ctx := c.Request.Context()
	result, err := h.service.SearchCandles(ctx, req)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, result)
}