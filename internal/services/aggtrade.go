package services

import (
	"context"
	"errors"
	"github.com/isyyyy/bnb-pub-sub/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sort"
	"sync"
	"time"
)

type AggTradeService struct {
	Collection       *mongo.Collection
	CandleCollection *mongo.Collection
	TimeFrame        map[string]time.Duration
	DefaultStartTime string
	DefaultEndTime   string
}

func NewAggTradeService(database *mongo.Database, tf map[string]time.Duration, collectionName string) *AggTradeService {
	collection := database.Collection(collectionName)
	candle := database.Collection("mycandle")
	startTime := time.Date(2021, 9, 21, 0, 0, 0, 0, time.Local)
	endTime := time.Now()
	return &AggTradeService{Collection: collection, CandleCollection: candle, TimeFrame: tf, DefaultStartTime: startTime.Format(time.RFC3339), DefaultEndTime: endTime.Format(time.RFC3339)}
}

func (a *AggTradeService) GetAll(ctx context.Context) (*[]interface{}, error) {
	query := bson.M{}
	cursor, err := a.Collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	var result []interface{}
	err = cursor.All(ctx, &result)
	return &result, nil
}

func (a *AggTradeService) InitMyCandles(ctx context.Context, symbol string, kline string) error {
	candles, err := a.searchAgg(ctx, symbol, kline)
	if err != nil {
		return err
	}

	var records []interface{}
	for _, c := range *candles {
		records = append(records, c)
	}

	return a.InsertMyCandles(ctx, records)
}

func (a AggTradeService) InsertMyCandles(ctx context.Context, records []interface{}) error {
	_, err := a.CandleCollection.InsertMany(ctx, records)
	if err != nil {
		return err
	}
	return nil
}

func (a *AggTradeService) searchAgg(ctx context.Context, symbol string, kline string) (*[]models.MyCandle, error) {

	tf, err := a.GetTimeFrame(a.DefaultStartTime, a.DefaultEndTime, kline)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var candles []models.MyCandle
	for _, t := range tf {
		wg.Add(1)
		go func(wg *sync.WaitGroup, t1 models.TimeFrame) {
			defer wg.Done()
			queryDate := bson.M{
				"timestamp": bson.M{
					"$gt": t1.Start.Format(time.RFC3339),
					"$lt": t1.End.Format(time.RFC3339),
				},
				"symbol": symbol,
			}
			cursor, err1 := a.Collection.Find(ctx, queryDate)
			if err1 != nil {
				log.Println(err1)
			}

			var records []models.AggTrade
			err2 := cursor.All(ctx, &records)
			if err2 != nil {
				log.Println(err2)
			}
			candle, err3 := CreateMyCandle(t1, records)
			if err3 == nil {
				candles = append(candles, candle)

			}
		}(&wg, t)

	}
	wg.Wait()
	sort.Slice(candles, func(i, j int) bool {
		return candles[i].StartTime < candles[j].StartTime
	})
	return &candles, nil
}
func (a *AggTradeService) SearchAggTrade(ctx context.Context, req models.AggRequest) (*[]models.MyCandle, error) {
	result, err := a.searchAgg(ctx, req.Symbol, req.Kline)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *AggTradeService) GetFirstLastRecord(ctx context.Context, symbol string) (string, string, error) {
	first := options.FindOne().SetSort(bson.D{{"timestamp", 1}})
	last := options.FindOne().SetSort(bson.D{{"timestamp", -1}})
	var firstRecord models.AggTrade
	err1 := a.Collection.FindOne(ctx, bson.M{"symbol": symbol}, first).Decode(&firstRecord)
	if err1 != nil {
		return "", "", err1
	}

	var lastRecord models.AggTrade
	err2 := a.Collection.FindOne(ctx, bson.M{}, last).Decode(&lastRecord)
	if err2 != nil {
		return "", "", err2
	}

	return firstRecord.Timestamp, lastRecord.Timestamp, nil
}

func (a *AggTradeService) GetTimeFrame(start, end string, frame string) ([]models.TimeFrame, error) {
	t1, err := time.Parse(time.RFC3339, start)
	if err != nil {
		return nil, err
	}

	t2, err := time.Parse(time.RFC3339, end)
	if err != nil {
		return nil, err
	}
	timeframe := a.TimeFrame[frame]

	var result []models.TimeFrame
	for {
		tmp := t1.Add(timeframe)
		if res := tmp.Sub(t2); res > 0 {
			tf := models.TimeFrame{
				Start:    t1,
				End:      t2,
				Interval: frame,
			}
			result = append(result, tf)
			break
		}
		tf := models.TimeFrame{
			Start:    t1,
			End:      tmp,
			Interval: frame,
		}
		result = append(result, tf)
		t1 = tmp
	}

	return result, nil
}

func CreateMyCandle(tf models.TimeFrame, records []models.AggTrade) (models.MyCandle, error) {

	if len(records) == 0 {
		return models.MyCandle{}, errors.New("no records")
	}
	openPrice := records[0].Price
	closePrice := records[len(records)-1].Price
	sort.Slice(records, func(i, j int) bool {
		return records[i].Price > records[j].Price
	})
	highPrice := records[0].Price
	lowPrice := records[len(records)-1].Price

	return models.MyCandle{
		Symbol:      records[0].Symbol,
		Interval:    tf.Interval,
		StartTime:   tf.Start.Format(time.RFC3339),
		CloseTime:   tf.End.Format(time.RFC3339),
		OpenPrice:   openPrice,
		ClosePrice:  closePrice,
		HighPrice:   highPrice,
		LowPrice:    lowPrice,
		BaseVolume:  0,
		QuoteVolume: 0,
	}, nil
}
