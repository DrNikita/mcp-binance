package mongo

import (
	"context"
	"encoding/json"
	"sync"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoServer struct {
	repo *mongoRepository
}

func NewMongoServer(client *mongo.Client) *MongoServer {
	repository := newMongoRepository(client)
	return &MongoServer{repository}
}

func (ms *MongoServer) TradesCreatorWorkerPool(ctx context.Context, gNumber int, ch chan []byte) chan error {
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
		close(ch)
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

func (ms *MongoServer) GetTradesInfo() {}
