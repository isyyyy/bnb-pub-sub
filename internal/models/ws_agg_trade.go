package models

import (
	"strconv"
	"time"
)

// AggTrade Websocket Market streams
type AggTrade struct {
	EventType    string  `json:"eventType" bson:"eventType"`
	EventTime    float64 `json:"eventTime" bson:"eventTime"`
	Symbol       string  `json:"symbol" bson:"symbol"`
	AggTradeID   float64 `json:"aggTradeID" bson:"aggTradeID"`
	Price        float64 `json:"price" bson:"price"`
	Quantity     float64 `json:"quantity" bson:"quantity"`
	FirstTradeID float64 `json:"firstTradeID" bson:"firstTradeID"`
	LastTradeID  float64 `json:"lastTradeID" bson:"lastTradeID"`
	TradeTime    float64 `json:"tradeTime" bson:"tradeTime"`
	IsTheMarket  bool    `json:"isTheMarket" bson:"isTheMarket"`
	Ignore       bool    `json:"ignore" bson:"ignore"`
	Timestamp    string  `json:"timestamp" bson:"timestamp"`
}

func ConvertToAggTrade(model map[string]interface{}) AggTrade {

	eventType := model["e"].(string)
	eventTime := model["E"].(float64)
	ts := ConvertToDateTime(eventTime)
	symbol := model["s"].(string)
	aggTradeID := model["a"].(float64)
	price, _ := strconv.ParseFloat(model["p"].(string), 64)
	quantity, _ := strconv.ParseFloat(model["q"].(string), 64)
	first := model["f"].(float64)
	last := model["l"].(float64)
	tradeTime := model["T"].(float64)
	market := model["m"].(bool)
	ignore := model["M"].(bool)

	return AggTrade{
		EventType:    eventType,
		EventTime:    eventTime,
		Symbol:       symbol,
		AggTradeID:   aggTradeID,
		Price:        price,
		Quantity:     quantity,
		FirstTradeID: first,
		LastTradeID:  last,
		TradeTime:    tradeTime,
		IsTheMarket:  market,
		Ignore:       ignore,
		Timestamp:    ts,
	}

}

func ConvertToDateTime(eventTime float64) string {
	t := time.Unix(0, int64(eventTime)*int64(time.Millisecond))
	ts := t.Format(time.RFC3339)
	return ts
}
