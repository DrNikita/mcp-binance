package main

import (
	"context"
	"fmt"
	"log"
	"mcpbinance/internal/infrastructure/stdio"
	"mcpbinance/internal/websocket"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	db "mcpbinance/internal/db/mongo"
)

func main() {
	ctx := context.TODO()

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
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("mongo client connected")
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	wg := sync.WaitGroup{}

	dbServer := db.NewMongoServer(client)
	parallel := 5
	errCh := dbServer.TradesCreatorWorkerPool(ctx, parallel, msgs)
	wg.Go(
		func() {
			defer wg.Done()
			for err := range errCh {
				fmt.Println(err.Error())
			}
		},
	)

	wg.Go(
		func() {
			defer wg.Done()
			if err := wsClient.Run(ctx, []string{"btcusdt"}, []string{"aggTrade"}); err != nil {
				log.Fatal(err)
			}
		},
	)

	mcpServer := mcp.NewServer(&mcp.Implementation{Name: "SCAM", Version: "v1.0.0"}, nil)
	stdioTransport := stdio.NewStdioTrarnsport(dbServer)
	stdioTransport.RegisterTools(mcpServer)

	if err := mcpServer.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
