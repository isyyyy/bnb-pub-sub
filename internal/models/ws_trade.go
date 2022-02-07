package models

import (
	"strconv"
)

// Trade Websocket Market streams
type Trade struct {
	EventType     string  `json:"eventType" bson:"eventType"`
	EventTime     float64 `json:"eventTime" bson:"eventTime"`
	Symbol        string  `json:"symbol" bson:"symbol"`
	TradeID       float64 `json:"tradeID" bson:"tradeID"`
	Price         float64 `json:"price" bson:"price"`
	Quantity      float64 `json:"quantity" bson:"quantity"`
	BuyerOrderID  float64 `json:"buyOrderID" bson:"buyerOrderID"`
	SellerOrderID float64 `json:"sellerOrderID" bson:"sellerOrderID"`
	TradeTime     float64 `json:"tradeTime" bson:"tradeTime"`
	IsTheMarket   bool    `json:"isTheMarket" bson:"isTheMarket"`
	Ignore        bool    `json:"ignore" bson:"ignore"`
	Timestamp     string  `json:"timestamp" json:"timestamp"`
}

func ConvertToTrade(model map[string]interface{}) Trade {

	eventType := model["e"].(string)
	eventTime := model["E"].(float64)
	ts := ConvertToDateTime(eventTime)
	symbol := model["s"].(string)
	tradeID := model["t"].(float64)
	price, _ := strconv.ParseFloat(model["p"].(string), 64)
	quantity, _ := strconv.ParseFloat(model["q"].(string), 64)
	buyer := model["b"].(float64)
	seller := model["a"].(float64)
	tradeTime := model["T"].(float64)
	market := model["m"].(bool)
	ignore := model["M"].(bool)

	return Trade{
		EventType:     eventType,
		EventTime:     eventTime,
		Symbol:        symbol,
		TradeID:       tradeID,
		Price:         price,
		Quantity:      quantity,
		BuyerOrderID:  buyer,
		SellerOrderID: seller,
		TradeTime:     tradeTime,
		IsTheMarket:   market,
		Ignore:        ignore,
		Timestamp:     ts,
	}

}
