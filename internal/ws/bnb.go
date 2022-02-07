package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/isyyyy/bnb-pub-sub/internal/config"
	"github.com/isyyyy/bnb-pub-sub/internal/models"
	"log"
	"strings"
)

const (
	aggTrade = "@aggTrade"
	trade    = "@trade"
	kline    = "@kline"
)

type WebSocket struct {
	URL    string
	Conn   *websocket.Conn
	Params []string
	MapToken map[string]*models.RealTimeToken
}

func NewWebSocket(bnbConfig config.BnbConfig) (*WebSocket, error) {
	conn, res, err := websocket.DefaultDialer.Dial(bnbConfig.URL, nil)
	if err != nil {
		return nil, err
	}
	log.Println(res)

	params,mapToken := BuildParams(bnbConfig)
	return &WebSocket{
		URL:    bnbConfig.URL,
		Conn:   conn,
		Params: params,
		MapToken: mapToken,
	}, nil

}

func BuildParams(bnbConfig config.BnbConfig) ([]string,map[string]*models.RealTimeToken) {
	var params []string
	var mapToken = make(map[string]*models.RealTimeToken)
	for _, symbol := range bnbConfig.ListSymbol {
		params = append(params, symbol+trade)
		params = append(params, symbol+aggTrade)
		for _, k := range bnbConfig.Kline {
			kl := fmt.Sprintf("%s_%s", symbol+kline, k)
			params = append(params, kl)
		}

		mapToken[strings.ToUpper(symbol)] = &models.RealTimeToken{
			Data:    models.MyCandle{
				Symbol: strings.ToUpper(symbol),
				Interval: "1m",
			},
			Records: []models.AggTrade{},
		}
	}
	return params,mapToken

}
