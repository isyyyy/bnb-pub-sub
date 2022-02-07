package services

import (
	"context"
	"github.com/isyyyy/bnb-pub-sub/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CandleService struct {
	Collection *mongo.Collection
}


func NewCandleService (data *mongo.Database, collectionName string) *CandleService{
	collection := data.Collection(collectionName)
	return &CandleService{Collection: collection}
}

func (s *CandleService) SearchCandles(ctx context.Context, req models.AggRequest) (*[]models.MyCandle, error) {
	option := options.Find().SetSort(bson.D{{"startTime", -1}}).SetLimit(req.Length)

	cursor, err := s.Collection.Find(ctx, bson.M{"symbol": req.Symbol}, option)
	if err != nil {
		return nil, err
	}

	var records []models.MyCandle
	err = cursor.All(ctx, &records)
	if err != nil {
		return nil, err
	}
	return &records, nil
}
