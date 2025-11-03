package mongo

import (
	"context"
	"errors"

	"mcpbinance/internal/entity"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mongoRepository struct {
	client *mongo.Client
}

func newMongoRepository(client *mongo.Client) *mongoRepository {
	repo := &mongoRepository{client}
	// Ensure unique index exists on AggregateTradeID
	repo.ensureUniqueIndex(context.Background())
	return repo
}

func (mr *mongoRepository) ensureUniqueIndex(ctx context.Context) {
	collection := mr.client.Database("mcp-binance").Collection("trade")
	
	// Create unique index on AggregateTradeID (field "a")
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "a", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("unique_aggregate_trade_id"),
	}
	
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		// Index might already exist, which is fine
		// Log but don't fail - duplicates will be caught on insert
	}
}

func (mr *mongoRepository) insertTradeInfo(ctx context.Context, tradeInfo map[string]any) error {
	collection := mr.client.Database("mcp-binance").Collection("trade")
	_, err := collection.InsertOne(ctx, tradeInfo)
	if err != nil {
		// Check if it's a duplicate key error
		var writeErr mongo.WriteException
		if errors.As(err, &writeErr) {
			for _, we := range writeErr.WriteErrors {
				// Error code 11000 is MongoDB duplicate key error
				if we.Code == 11000 {
					// Duplicate detected - ignore it
					return nil
				}
			}
		}
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
