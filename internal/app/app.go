package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/isyyyy/bnb-pub-sub/internal/config"
	"github.com/isyyyy/bnb-pub-sub/internal/handlers"
	"github.com/isyyyy/bnb-pub-sub/internal/job"
	"github.com/isyyyy/bnb-pub-sub/internal/models"
	"github.com/isyyyy/bnb-pub-sub/internal/mq"
	"github.com/isyyyy/bnb-pub-sub/internal/route"
	"github.com/isyyyy/bnb-pub-sub/internal/services"
	"github.com/isyyyy/bnb-pub-sub/internal/ws"
	"github.com/isyyyy/bnb-pub-sub/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
)

type App struct {
	Port      string
	Subscribe *mq.MQConnection
	Publish   *mq.MQConnection
	WebSocket *ws.WebSocket
	Writer    *database.Inserter
	Handle    *handlers.Handler
	Job       *job.JobScheduler
	Producer  *mq.Producer
	Consumer  *mq.Consumer
}

var MapToken = make(map[string]*models.RealTimeToken)

func NewApp(ctx context.Context) *App {

	root := config.InitConfig()
	log.Println("Load config ", root)

	db, err := database.NewMongoConnection(ctx, root.Mongo)
	if err != nil {
		log.Printf("cann't connect to Mongo: %s", err)
	}
	log.Println(db)
	var initData = true
	names, err := db.ListCollectionNames(ctx, bson.M{})
	for _, name := range names {
		if name == "mycandle" {
			initData = false
		}
	}
	log.Println(initData)

	inserter := database.NewInserter(db)
	subscribe, err := mq.NewRabbitMQConnection(root.RabbitMQ)
	if err != nil {
		log.Printf("cann't connect to RabbitMQ: %s", err)
	}

	publish, err := mq.NewRabbitMQConnection(root.RabbitMQ)
	if err != nil {
		log.Printf("cann't connect to RabbitMQ: %s", err)
	}

	ws, err := ws.NewWebSocket(root.Bnb)
	if err != nil {
		log.Printf("cann't connect to BNB Websocket: %s", err)
	}

	producer := mq.NewProducer(publish, ws)
	consumer := mq.NewConsumer(subscribe, inserter)
	service := services.NewServices(db)
	handler := handlers.NewHandlers(service)

	job := job.NewJob(service, root.Bnb, initData)

	MapToken = ws.MapToken
	return &App{
		Port:      root.Server.Port,
		Subscribe: subscribe,
		Publish:   publish,
		WebSocket: ws,
		Writer:    inserter,
		Handle:    handler,
		Job:       job,
		Producer:  producer,
		Consumer:  consumer,
	}

}

func (a *App) Server() {

	r := gin.Default()
	route.NewRoute(r, a.Handle)
	realTime := "api/v1/realtime"
	r.POST(realTime, func(c *gin.Context) {
		var req models.AggRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.IndentedJSON(http.StatusBadRequest, nil)
			return
		}
		res := MapToken[req.Symbol]
		c.IndentedJSON(http.StatusOK, res.Data)
	})
	server := ""
	if a.Port != "" {
		server = ":" + a.Port
	}
	r.Run(server)

}
