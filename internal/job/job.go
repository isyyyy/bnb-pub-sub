package job

import (
	"context"
	"github.com/go-co-op/gocron"
	"github.com/isyyyy/bnb-pub-sub/internal/config"
	"github.com/isyyyy/bnb-pub-sub/internal/services"
	"log"
	"strings"
	"sync"
	"time"
)

type JobScheduler struct {
	Scheduler *gocron.Scheduler
	Service   *services.Service
	BnbConfig config.BnbConfig
	IsEmpty   bool
}

func Test(a string, b string) {
	log.Printf("Symbol : %s - Kline : %s \n", a, b)
}

func NewJob(service *services.Service, bnbConfig config.BnbConfig, isEmpty bool) *JobScheduler {
	scheduler := gocron.NewScheduler(time.UTC)

	return &JobScheduler{Scheduler: scheduler, Service: service, BnbConfig: bnbConfig, IsEmpty: isEmpty}
}

func (j *JobScheduler) InitData(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(len(j.BnbConfig.ListSymbol) * len(j.BnbConfig.Kline))
	log.Println(len(j.BnbConfig.ListSymbol) * len(j.BnbConfig.Kline))
	for _, symbol := range j.BnbConfig.ListSymbol {
		for _, kline := range j.BnbConfig.Kline {
			go func(ctx1 context.Context, wg *sync.WaitGroup, s string, k string) {
				defer wg.Done()
				log.Printf("Symbol : %s - Kline : %s \n", s, k)
				j.Service.AggTradeService.InitMyCandles(ctx1, s, k)
				log.Printf("Done - Symbol : %s - Kline : %s \n", s, k)
			}(ctx, &wg, strings.ToUpper(symbol), kline)
		}
	}
	wg.Wait()
	log.Println("Sync Done")
}

func (j *JobScheduler) StartJobSchedule(ctx context.Context) {
	if j.IsEmpty {
		j.InitData(ctx)
	}



	log.Println("Job will run now")
	//j.Scheduler.StartAsync()
}
