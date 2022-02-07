package route

import (
	"github.com/gin-gonic/gin"
	"github.com/isyyyy/bnb-pub-sub/internal/handlers"
)

const (
	baseurl = "/api/v1"
)

func NewRoute(r *gin.Engine, handle *handlers.Handler) {
	aggTrade := baseurl + "/aggtrade"
	r.GET(aggTrade, handle.AggTradeHandler.GetAll)
	r.POST(aggTrade, handle.AggTradeHandler.SearchAggTrade)

	candle  := baseurl + "/candle"
	r.POST(candle,handle.CandleHandler.SearchCandles)

}
