package database

import (
	"context"
	"github.com/isyyyy/bnb-pub-sub/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoConnection(ctx context.Context, conf config.MongoConfig) (*mongo.Database, error) {
	option := options.Client().ApplyURI(conf.URL)
	client, err := mongo.Connect(ctx, option)
	if err != nil {
		return nil, err
	}
	db := client.Database(conf.Database)
	return db, nil
}
