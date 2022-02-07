package models

import "strconv"

type Kline struct {
	EventType string      `json:"eventType" bson:"eventType"`
	EventTime float64     `json:"eventTime" bson:"eventTime"`
	Symbol    string      `json:"symbol" bson:"symbol"`
	Candle    Candlestick `json:"candle" bson:"candle"`
	Timestamp string      `json:"timestamp" bson:"timestamp"`
}

type Candlestick struct {
	KlineStartTime float64 `json:"klineStartTime" bson:"klineStartTime"`
	KlineCloseTime float64 `json:"klineCloseTime" bson:"klineCloseTime"`
	Symbol         string  `json:"symbol" bson:"symbol"`
	Interval       string  `json:"interval" bson:"interval"`
	FirstTradeID   float64 `json:"firstTradeID" bson:"firstTradeID"`
	LastTradeID    float64 `json:"lastTradeID" bson:"lastTradeID"`
	OpenPrice      float64 `json:"openPrice" bson:"openPrice"`
	ClosePrice     float64 `json:"closePrice" bson:"closePrice"`
	HighPrice      float64 `json:"highPrice" bson:"highPrice"`
	LowPrice       float64 `json:"lowPrice" bson:"lowPrice"`
	BaseVolume     float64 `json:"baseVolume" bson:"baseVolume"`
	NumberTrades   float64 `json:"numberTrades" bson:"numberTrades"`
	IsClosed       bool    `json:"isClosed" bson:"isClosed"`
	QuoteVolume    float64 `json:"quoteVolume" bson:"quoteVolume"`
	BuyBaseVolume  float64 `json:"buyBaseVolume" bson:"buyBaseVolume"`
	BuyQuoteVolume float64 `json:"buyQuoteVolume" bson:"buyQuoteVolume"`
	Ignore         string  `json:"ignore" bson:"ignore"`
}

func ConvertToKline(model map[string]interface{}) Kline {
	eventType := model["e"].(string)
	eventTime := model["E"].(float64)
	ts := ConvertToDateTime(eventTime)
	symbol := model["s"].(string)
	kline := model["k"].(map[string]interface{})
	kStartTime := kline["t"].(float64)
	kCloseTime := kline["T"].(float64)
	kSymbol := kline["s"].(string)
	kInterval := kline["i"].(string)
	kFirst := kline["f"].(float64)
	kLast := kline["L"].(float64)
	kOpen, _ := strconv.ParseFloat(kline["o"].(string), 64)
	kClose, _ := strconv.ParseFloat(kline["c"].(string), 64)
	kHigh, _ := strconv.ParseFloat(kline["h"].(string), 64)
	kLow, _ := strconv.ParseFloat(kline["l"].(string), 64)
	kBaseVol, _ := strconv.ParseFloat(kline["v"].(string), 64)
	kNum := kline["n"].(float64)
	kIsClose := kline["x"].(bool)
	kQuoteVol, _ := strconv.ParseFloat(kline["q"].(string), 64)
	kBBaseVol, _ := strconv.ParseFloat(kline["V"].(string), 64)
	kBQuoteVol, _ := strconv.ParseFloat(kline["Q"].(string), 64)
	kIgnore := kline["B"].(string)
	return Kline{
		EventType: eventType,
		EventTime: eventTime,
		Symbol:    symbol,
		Timestamp: ts,
		Candle: Candlestick{
			KlineStartTime: kStartTime,
			KlineCloseTime: kCloseTime,
			Symbol:         kSymbol,
			Interval:       kInterval,
			FirstTradeID:   kFirst,
			LastTradeID:    kLast,
			OpenPrice:      kOpen,
			ClosePrice:     kClose,
			HighPrice:      kHigh,
			LowPrice:       kLow,
			BaseVolume:     kBaseVol,
			NumberTrades:   kNum,
			IsClosed:       kIsClose,
			QuoteVolume:    kQuoteVol,
			BuyBaseVolume:  kBBaseVol,
			BuyQuoteVolume: kBQuoteVol,
			Ignore:         kIgnore,
		},
	}

}
