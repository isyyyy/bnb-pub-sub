package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

type HealthCheck struct {
	db      *mongo.Database
	name    string
	timeout time.Duration
}

func NewHealthChecker(db *mongo.Database) *HealthCheck {
	return &HealthCheck{
		db:      db,
		name:    "mongo",
		timeout: 4 * time.Second,
	}
}

func (h *HealthCheck) Check(ctx context.Context) (map[string]interface{}, error) {
	cancel := func() {}
	if h.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, h.timeout)
	}
	defer cancel()

	res := make(map[string]interface{})
	info := make(map[string]interface{})
	checkerChan := make(chan error)
	go func() {
		checkerChan <- h.db.RunCommand(ctx, bsonx.Doc{{"ping", bsonx.Int32(1)}}).Decode(&info)
	}()
	select {
	case err := <-checkerChan:
		return res, err
	case <-ctx.Done():
		return res, fmt.Errorf("timeout")
	}
}
func (h *HealthCheck) Build(ctx context.Context, data map[string]interface{}, err error) map[string]interface{} {
	if err == nil {
		return data
	}
	if data == nil {
		data = make(map[string]interface{}, 0)
	}
	data["error"] = err
	return data
}
