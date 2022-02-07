package mq

import (
	"context"
	"encoding/json"
	"github.com/isyyyy/bnb-pub-sub/internal/models"
	"github.com/isyyyy/bnb-pub-sub/pkg/database"
	"log"
	"time"
)

type Consumer struct {
	RabbitMQ *MQConnection
	Inserter *database.Inserter
	StartTime time.Time
}

func NewConsumer (mq *MQConnection, inserter  *database.Inserter ) *Consumer {
	return &Consumer{
		RabbitMQ: mq,
		Inserter: inserter,
		StartTime: time.Now(),
	}
}



func (c *Consumer) Consume(ctx context.Context)  {

	msgs, err := c.RabbitMQ.Channel.Consume(c.RabbitMQ.QueueName, "", true, false, false, false, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println(MapToken)
	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			//log.Println(string(msg.Body))
			var model map[string]interface{}
			er1 := json.Unmarshal(msg.Body, &model)
			if err != nil {
				log.Println(er1)
			}

			event := model["e"]
			switch event {
			case "trade":
				trade := models.ConvertToTrade(model)
				//log.Println(trade)
				c.Inserter.Trade.InsertOne(ctx, trade)
			case "aggTrade":
				aggTrade := models.ConvertToAggTrade(model)
				atTime := time.Unix(0, int64(aggTrade.EventTime)*int64(time.Millisecond))
				min := atTime.Minute() - c.StartTime.Minute()
				if min == 0 {
					MapToken[aggTrade.Symbol].UpdateValue(aggTrade)
					//log.Printf("%s - Token : %s - Open price : %f - Close price: %f - Low price : %f - High price : %f - Price : %f\n",atTime.String(), aggTrade.Symbol, MapToken[aggTrade.Symbol].Data.OpenPrice, MapToken[aggTrade.Symbol].Data.ClosePrice, MapToken[aggTrade.Symbol].Data.LowPrice, MapToken[aggTrade.Symbol].Data.HighPrice, aggTrade.Price)
				} else {
					//log.Println("--------------------------------------------------------------------------------------------")
					for _, token := range MapToken {
						if len(token.Records) != 0 {
							token.Data.ClosePrice = token.Records[len(token.Records)-1].Price
						}
						token.Data.StartTime = c.StartTime.Format(time.RFC3339)
						token.Data.CloseTime = c.StartTime.Add(time.Minute).Format(time.RFC3339)
						value := token.Data
						c.Inserter.Candles.InsertOne(ctx, value)
						//log.Printf("Candle %s at : %s \n", value.Symbol, c.StartTime.Format(time.RFC3339))
						//log.Printf("Open price : %f - Close price: %f - Low price : %f - High price : %f\n", value.OpenPrice, value.ClosePrice, value.LowPrice, value.HighPrice)
						//log.Println("--------------------------------------------------------------------------------------------")
					}
					c.StartTime = c.StartTime.Add(time.Minute)
					MapToken[aggTrade.Symbol].ResetValue()
					MapToken[aggTrade.Symbol].UpdateValue(aggTrade)

				}
				c.Inserter.AggTrade.InsertOne(ctx, aggTrade)
			case "kline":
				kline := models.ConvertToKline(model)
				//log.Println(kline)
				c.Inserter.Kline.InsertOne(ctx, kline)
			}
		}
	}()
	<-forever

}