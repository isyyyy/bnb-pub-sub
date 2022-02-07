package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/isyyyy/bnb-pub-sub/internal/models"
	"github.com/isyyyy/bnb-pub-sub/internal/services"
	"net/http"
)

type AggTradeHandler struct {
	service services.AggTradeService
}

func NewAggTradeHandler(aggTradeService *services.AggTradeService) *AggTradeHandler {
	return &AggTradeHandler{service: *aggTradeService}
}

func (a *AggTradeHandler) GetAll(c *gin.Context) {
	result, err := a.service.GetAll(c.Request.Context())
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, nil)
		return
	}
	c.IndentedJSON(http.StatusOK, result)

}

func (a *AggTradeHandler) SearchAggTrade(c *gin.Context) {

	var req models.AggRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, nil)
		return

	}
	ctx := c.Request.Context()
	result, err := a.service.SearchAggTrade(ctx, req)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusOK, result)

}


