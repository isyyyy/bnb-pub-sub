package handlers

import (
	"github.com/isyyyy/bnb-pub-sub/internal/services"
)

type Handler struct {
	AggTradeHandler *AggTradeHandler
	CandleHandler *CandleHandler
}

func NewHandlers(service *services.Service) *Handler {
	aggTradeHandler := NewAggTradeHandler(service.AggTradeService)
	candleHandler := NewCandleHandler(service.CandleService)
	return &Handler{AggTradeHandler: aggTradeHandler, CandleHandler: candleHandler}
}

