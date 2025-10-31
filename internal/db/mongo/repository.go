package mongo

import (
	"context"

	"mcpbinance/internal/entity"

	"go.mongodb.org/mongo-driver/v2/bson"
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

func (r *mongoRepository) findTradesInfo(ctx context.Context, filter bson.M) ([]entity.Trade, error) {
	collection := r.client.Database("mcp-binance").Collection("trade")
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	trades := make([]entity.Trade, 0)
	for cursor.Next(ctx) {
		var trade entity.Trade
		if err := bson.Unmarshal(cursor.Current, &trade); err != nil {
			return nil, err
		}
		trades = append(trades, trade)
	}

	return trades, nil
}
