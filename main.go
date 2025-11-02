package main

import (
	"context"
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
				log.Printf("trade processing error: %v", err)
			}
		},
	)

	mcpServer := mcp.NewServer(&mcp.Implementation{Name: "SCAM", Version: "v1.0.0"}, nil)
	monitorService := websocket.NewStockMonitorService(wsClient, &wg)

	stdioTransport := stdio.NewStdioTrarnsport(dbServer, monitorService)
	stdioTransport.RegisterTools(mcpServer)

	monitorService.RunSymbolsMonitoring(ctx, []string{"etcusdt"}, []string{"aggTrade"})

	if err := mcpServer.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
