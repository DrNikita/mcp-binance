package mongo

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	"mcpbinance/internal/entity"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoServer struct {
	repo *mongoRepository
}

func NewMongoServer(client *mongo.Client) *MongoServer {
	repository := newMongoRepository(client)
	return &MongoServer{repository}
}

func (ms *MongoServer) TradesCreatorWorkerPool(ctx context.Context, gNumber int, ch <-chan []byte) chan error {
	errChan := make(chan error)
	wg := sync.WaitGroup{}

	wg.Add(gNumber)
	for range gNumber {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					break
				case tradeInfo := <-ch:
					err := ms.createTradeInfo(ctx, tradeInfo)
					if err != nil {
						errChan <- err
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	return errChan
}

func (ms *MongoServer) createTradeInfo(ctx context.Context, tradeInfo []byte) error {
	tradeInfoMap := make(map[string]any)
	err := json.Unmarshal(tradeInfo, &tradeInfoMap)
	if err != nil {
		return err
	}

	if err := ms.repo.insertTradeInfo(ctx, tradeInfoMap); err != nil {
		return err
	}

	return nil
}

func (ms *MongoServer) GetTradesInfo(ctx context.Context, symbol string, periodSeconds int) ([]entity.Trade, error) {
	upperSymbol := strings.ToUpper(symbol)
	filter := createFilterParams(upperSymbol, periodSeconds)
	trades, err := ms.repo.findTradesInfo(ctx, filter)
	if err != nil {
		return nil, err
	}

	log.Printf("filter: %v", filter)

	return trades, nil
}

func createFilterParams(symbol string, seconds int) bson.M {
	now := time.Now().UnixMilli()
	dateFrom := now - int64(seconds*1000)
	return bson.M{"s": symbol, "E": bson.M{"$gte": dateFrom}}
}
