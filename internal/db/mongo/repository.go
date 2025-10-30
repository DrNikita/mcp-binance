package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type mongoRepository struct {
	client *mongo.Client
}

func newMongoRepository(client *mongo.Client) *mongoRepository {
	return &mongoRepository{client}
}

func (mr *mongoRepository) insertTradeInfo(ctx context.Context, tradeInfo map[string]any) error {
	collection := mr.client.Database("mcp-binance").Collection("trade")
	_, err := collection.InsertOne(ctx, tradeInfo)
	if err != nil {
		return err
	}
	return nil
}

func (r *mongoRepository) findTradeInfo(ctx context.Context) {
	// collection := r.client.Database("mcp-binance").Collection("trade")
}
