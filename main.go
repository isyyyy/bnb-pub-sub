package main

import (
	"context"
	"github.com/isyyyy/bnb-pub-sub/internal/app"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	time.Local = loc
	if err != nil {
		log.Println(err)
	}
	app := app.NewApp(ctx)
	log.Println("Start")

	//go app.Produce()
	//go app.Consume(ctx)
	//go app.Job.StartJobSchedule(ctx)
	//app.Start(ctx)

	go app.Producer.Produce()
	go app.Consumer.Consume(ctx)
	app.Server()
	log.Println("Close Program")

}
