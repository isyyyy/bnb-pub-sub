package database

import (
	"github.com/isyyyy/bnb-pub-sub/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type Inserter struct {
	Trade    *mongo.Collection
	AggTrade *mongo.Collection
	Kline    *mongo.Collection
	Candles  *mongo.Collection
}

func NewInserter(database *mongo.Database) *Inserter {
	trade := database.Collection("trade")
	aggtrade := database.Collection("aggtrade")
	kline := database.Collection("kline")
	candles := database.Collection("candles")
	return &Inserter{
		Trade:    trade,
		AggTrade: aggtrade,
		Kline:    kline,
		Candles:  candles,
	}
}

func CreateCandle(mapToken map[string]*models.RealTimeToken, data models.AggTrade, startTime time.Time) {
	atTime := time.Unix(0, int64(data.EventTime)*int64(time.Millisecond))
	diff := atTime.Sub(startTime)
	if diff < time.Minute && diff >= 0 {
		mapToken[data.Symbol].UpdateValue(data)
		log.Println(mapToken[data.Symbol].Data)
	} else {
		//mapToken[data.Symbol].InsertValue()
		value := mapToken[data.Symbol].Data
		log.Println("Candle ", value)
		startTime = startTime.Add(time.Minute)
		mapToken[data.Symbol].ResetValue()

	}
}
