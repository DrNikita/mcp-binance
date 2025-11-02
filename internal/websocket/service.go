package websocket

import (
	"context"
	"sync"
)

type StockMonitorService struct {
	wsClient *WebsocketClient
	wg       *sync.WaitGroup
}

func NewStockMonitorService(wsClient *WebsocketClient, wg *sync.WaitGroup) *StockMonitorService {
	return &StockMonitorService{wsClient, wg}
}

func (sm *StockMonitorService) RunSymbolsMonitoring(ctx context.Context, symbols, streamTypes []string) {
	sm.wg.Go(
		func() {
			if err := sm.wsClient.Run(ctx, symbols, streamTypes); err != nil {
				// TODO: log err
			}
		},
	)
}
