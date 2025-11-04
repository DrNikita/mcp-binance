package application

import (
	"context"
	"log"
	"mcpbinance/internal/config"
	"mcpbinance/internal/infrastructure/stdio"
	"mcpbinance/internal/websocket"
	"mcpbinance/internal/websocket/enum"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	db "mcpbinance/internal/db/mongo"
)

type Application struct {
	cfg *config.EnvConfig
}

func NewApplication(cfg *config.EnvConfig) *Application {
	return &Application{cfg}
}

func (a *Application) Run(ctx context.Context) error {
	clientParams := websocket.ClientParams{
		URL:          a.cfg.WsURL,
		MsgCap:       a.cfg.WsMsgCap,
		ReadLimit:    a.cfg.WsReadLimit,
		ReconnPeriod: time.Duration(a.cfg.WsReconnPeriodHr) * time.Hour,
		RetryBackoff: time.Duration(a.cfg.WsRetryBackoffSec) * time.Second,
		PingWait:     time.Duration(a.cfg.WsPingWaitMin) * time.Minute,
		PongPeriod:   time.Duration(a.cfg.WsPongPeriodMin) * time.Minute,
	}

	wsClient, err := websocket.NewWebsocketClient(clientParams)
	if err != nil {
		return err
	}

	msgs := wsClient.Receive()
	client, err := mongo.Connect(options.Client().ApplyURI(a.cfg.DbURL))
	if err != nil {
		return err
	}
	// fmt.Println("mongo client connected")
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
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

	mcpServer := mcp.NewServer(&mcp.Implementation{Name: a.cfg.ServerName, Version: a.cfg.ServerVersion}, nil)
	monitorService := websocket.NewStockMonitorService(wsClient, &wg)

	stdioTransport := stdio.NewStdioTrarnsport(dbServer, monitorService)
	stdioTransport.RegisterTools(mcpServer)

	symbols, streamTypes, err := enum.CreateStreamParams(a.cfg.StartupSymbols, a.cfg.StartupStreamTypes)
	if err != nil {
		return err
	}
	monitorService.RunSymbolsMonitoring(ctx, symbols, streamTypes)

	if err := mcpServer.Run(ctx, &mcp.StdioTransport{}); err != nil {
		return err
	}

	wg.Wait()
	return nil
}
