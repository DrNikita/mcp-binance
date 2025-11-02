package stdio

import (
	"context"
	"mcpbinance/internal/entity"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Store interface {
	TradesCreatorWorkerPool(context.Context, int, <-chan []byte) chan error
	GetTradesInfo(ctx context.Context, symbol string, periodSeconds int) ([]entity.Trade, error)
}

type StockMonitorer interface {
	RunSymbolsMonitoring(ctx context.Context, symbols, streamTypes []string)
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
	st.stockMonitorService.RunSymbolsMonitoring(ctx, input.Symbols, input.StreamTypes)
	return nil, struct{}{}, nil
}
