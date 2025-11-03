package stdio

import (
	"context"
	"fmt"
	"mcpbinance/internal/entity"
	"mcpbinance/internal/websocket/enum"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Store interface {
	TradesCreatorWorkerPool(context.Context, int, <-chan []byte) chan error
	GetTradesInfo(ctx context.Context, symbol string, periodSeconds int) ([]entity.Trade, error)
}

type StockMonitorer interface {
	RunSymbolsMonitoring(ctx context.Context, symbols []enum.Symbol, streamTypes []enum.StreamType)
}

type StdioTrarnsport struct {
	dbServer            Store
	stockMonitorService StockMonitorer
}

func NewStdioTrarnsport(dbServer Store, stockMonitorService StockMonitorer) *StdioTrarnsport {
	return &StdioTrarnsport{dbServer, stockMonitorService}
}

func (st *StdioTrarnsport) RegisterTools(mcpServer *mcp.Server) {
	mcp.AddTool(mcpServer, &mcp.Tool{Name: "getTradesHistory", Description: "get price changes history of symbol"}, st.GetTradePairsHistory)
	mcp.AddTool(mcpServer, &mcp.Tool{Name: "runSymbolsMonitoring", Description: "run symbols price monitoring"}, st.RunSymbolsMonitoring)
}

func (st *StdioTrarnsport) GetTradePairsHistory(ctx context.Context, req *mcp.CallToolRequest, input GetTradePairsHistoryInput) (*mcp.CallToolResult, GetTradePairsHistoryOutput, error) {
	trades, err := st.dbServer.GetTradesInfo(ctx, input.Symbol, input.Seconds)
	if err != nil {
		return nil, GetTradePairsHistoryOutput{}, err
	}

	tradesInfo := GetTradePairsHistoryOutput{trades}
	return nil, tradesInfo, nil
}

func (st *StdioTrarnsport) RunSymbolsMonitoring(ctx context.Context, req *mcp.CallToolRequest, input RunStockMonitoringInput) (*mcp.CallToolResult, struct{}, error) {
	symbols, streamTypes, err := createStreamParams(input)
	if err != nil {
		return nil, struct{}{}, err
	}

	st.stockMonitorService.RunSymbolsMonitoring(ctx, symbols, streamTypes)
	return nil, struct{}{}, nil
}

func createStreamParams(monitoringData RunStockMonitoringInput) ([]enum.Symbol, []enum.StreamType, error) {
	symbols := make([]enum.Symbol, 0)
	streamTypes := make([]enum.StreamType, 0)
	for _, symbol := range monitoringData.Symbols {
		eSymbol, err := enum.NewSymbol(symbol)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %w", "failed to run symbols monitoring", err)
		}
		symbols = append(symbols, eSymbol)
	}
	for _, streamType := range monitoringData.StreamTypes {
		eStreamType, err := enum.NewStreamType(streamType)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %w", "failed to run symbols monitoring", err)
		}
		streamTypes = append(streamTypes, eStreamType)
	}

	return symbols, streamTypes, nil
}
