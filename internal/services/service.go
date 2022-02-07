package services

import (
	"github.com/isyyyy/bnb-pub-sub/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	AggTradeService  *AggTradeService
	CandleService *CandleService
}

func NewServices(database *mongo.Database) *Service {
	aggTradeService := NewAggTradeService(database, models.KlineRequest, "aggtrade")
	candleService := NewCandleService(database,"candles")
	return &Service{AggTradeService: aggTradeService,CandleService: candleService}
}

