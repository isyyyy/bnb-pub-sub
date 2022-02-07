package mq

import (
	"github.com/gorilla/websocket"
	"github.com/isyyyy/bnb-pub-sub/internal/models"
	"github.com/isyyyy/bnb-pub-sub/internal/ws"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

type Producer struct {
	RabbitMQ *MQConnection
	BnbWS    *ws.WebSocket
}

var MapToken = make(map[string]*models.RealTimeToken)

func NewProducer(mq *MQConnection, ws *ws.WebSocket) *Producer {
	MapToken = ws.MapToken
	return &Producer{
		RabbitMQ: mq,
		BnbWS:    ws,
	}
}

func (p *Producer) Produce() {
	req := map[string]interface{}{
		"method": "SUBSCRIBE",
		"params": p.BnbWS.Params,
		"id":     1,
	}

	err := p.BnbWS.Conn.WriteJSON(req)
	if err != nil {
		log.Println(err)
	}

	var conn *websocket.Conn
	var res *http.Response
	var err3 error
	conn = p.BnbWS.Conn
	isClose := false
	for {
		for {
			_, message, err1 := conn.ReadMessage()
			if err1 != nil {
				if websocket.IsUnexpectedCloseError(err1) {
					conn.Close()
					isClose = true
					break
				}
				log.Println(err1)
			}

			err2 := p.RabbitMQ.Channel.Publish("", p.RabbitMQ.QueueName, false, false, amqp.Publishing{
				ContentType: "application/json",
				Body:        message,
			})
			if err2 != nil {
				log.Println(err2)
			}
		}
		if isClose {
			conn, res, err3 = websocket.DefaultDialer.Dial(p.BnbWS.URL, nil)
			if err3 != nil {
				log.Println(err3)
				continue
			}
			log.Printf("Reconnected : %s ", res)
			err4 := conn.WriteJSON(req)
			if err4 != nil {
				log.Println(err4)
				continue
			}
			isClose = false
		}
	}
}
