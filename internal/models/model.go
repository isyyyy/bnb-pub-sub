package models

import "time"

var KlineRequest = map[string]time.Duration{
	"1m":  time.Minute,
	"2m":  time.Minute * 2,
	"3m":  time.Minute * 3,
	"5m":  time.Minute * 5,
	"10m": time.Minute * 10,
	"15m": time.Minute * 15,
	"30m": time.Minute * 30,
	"45m": time.Minute * 45,
	"1h":  time.Hour,
	"2h":  time.Hour * 2,
	"3h":  time.Hour * 3,
	"4h":  time.Hour * 4,
	"5h":  time.Hour * 5,
	"6h":  time.Hour * 6,
	"8h":  time.Hour * 8,
	"12h": time.Hour * 12,
	"1d":  time.Hour * 24,
	"1w":  time.Hour * 24 * 7,
}

type AggRequest struct {
	Symbol string `json:"symbol"`
	Kline  string `json:"kline"`
	Length int64 `json:"length"`
}

type TimeFrame struct {
	Start    time.Time
	End      time.Time
	Interval string
}

type MyCandle struct {
	Symbol      string  `json:"symbol" bson:"symbol"`
	Interval    string  `json:"interval" bson:"interval"`
	StartTime   string  `json:"startTime" bson:"startTime"`
	CloseTime   string  `json:"closeTime" bson:"closeTime"`
	OpenPrice   float64 `json:"openPrice" bson:"openPrice"`
	ClosePrice  float64 `json:"closePrice" bson:"closePrice"`
	HighPrice   float64 `json:"highPrice" bson:"highPrice"`
	LowPrice    float64 `json:"lowPrice" bson:"lowPrice"`
	BaseVolume  float64 `json:"baseVolume" bson:"baseVolume"`
	QuoteVolume float64 `json:"quoteVolume" bson:"quoteVolume"`
}

type RealTimeToken struct {
	Data MyCandle
	Records []AggTrade
}

func (r *RealTimeToken) UpdateValue(aggTrade AggTrade) {
	r.Records  = append(r.Records, aggTrade)
	r.Data.Symbol = aggTrade.Symbol
	r.Data.Interval = "1m"
	if len(r.Records) == 1 {
		price := r.Records[0].Price
		r.Data.OpenPrice = price
		r.Data.ClosePrice = price
		r.Data.HighPrice = price
		r.Data.LowPrice = price
	}
	price := aggTrade.Price
	if price > r.Data.HighPrice {
		r.Data.HighPrice = price
	}

	if price < r.Data.LowPrice {
		r.Data.LowPrice = price
	}
}


func (r *RealTimeToken) ResetValue() {
	r.Data = MyCandle{}
	r.Records = []AggTrade{}
}