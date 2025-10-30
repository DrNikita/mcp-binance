package main

import (
	"context"
	"fmt"
	"log"
	"mcpbinance/internal/websocket"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	db "mcpbinance/internal/db/mongo"
)

func main() {
	ctx := context.TODO()

	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("mongo client connected")
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	dbServer := db.NewMongoServer(client)

	tradesInfoCh := make(chan []byte)
	errCh := dbServer.TradesCreatorWorkerPool(ctx, 5, tradesInfoCh)
	go func() {
		for err := range errCh {
			fmt.Println(err.Error())
		}
	}()

	clientParams := websocket.ClientParams{
		URL:          "wss://fstream.binance.com/ws",
		MsgCap:       100,
		ReadLimit:    3000000,
		ReconnPeriod: 24 * time.Hour,
		RetryBackoff: 5 * time.Second,
		PingWait:     3 * time.Minute,
		PongPeriod:   10 * time.Minute,
	}

	wsClient, err := websocket.NewWebsocketClient(clientParams)
	if err != nil {
		log.Fatal(err)
	}

	msgs := wsClient.Receive()
	go func() {
		for {
			tradesInfoCh <- <-msgs
		}
	}()

	if err := wsClient.Run(ctx, []string{"etcusdt"}, []string{"aggTrade"}); err != nil {
		log.Fatal(err)
	}
}
