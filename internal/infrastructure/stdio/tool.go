package stdio

import (
	"context"
	"mcpbinance/internal/entity"
	"mcpbinance/internal/websocket/enum"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Store interface {
	TradesCreatorWorkerPool(context.Context, int, <-chan []byte) chan error
	GetTradesInfo(ctx context.Context, symbol string, dateFrom, dateTo int) ([]entity.Trade, error)
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
	trades, err := st.dbServer.GetTradesInfo(ctx, input.Symbol, input.DateFromTimestamp, input.DateToTimestamp)
	if err != nil {
		return nil, GetTradePairsHistoryOutput{}, err
	}

	tradesInfo := GetTradePairsHistoryOutput{trades}
	return nil, tradesInfo, nil
}

func (st *StdioTrarnsport) RunSymbolsMonitoring(ctx context.Context, req *mcp.CallToolRequest, input RunStockMonitoringInput) (*mcp.CallToolResult, struct{}, error) {
	symbols, streamTypes, err := enum.CreateStreamParams(input.Symbols, input.StreamTypes)
	if err != nil {
		return nil, struct{}{}, err
	}

	st.stockMonitorService.RunSymbolsMonitoring(ctx, symbols, streamTypes)
	return nil, struct{}{}, nil
}
