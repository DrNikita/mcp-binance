package stdio

import (
	"context"
	"mcpbinance/internal/entity"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const parallel = 5

type Store interface {
	TradesCreatorWorkerPool(context.Context, int, <-chan []byte) chan error
	GetTradesInfo(ctx context.Context, periodSeconds int) ([]entity.Trade, error)
}

type StdioTrarnsport struct {
	dbServer Store
}

func NewStdioTrarnsport(dbServer Store) *StdioTrarnsport {
	return &StdioTrarnsport{dbServer}
}

func (st *StdioTrarnsport) RegisterTools(mcpServer *mcp.Server) {
	mcp.AddTool(mcpServer, &mcp.Tool{Name: "getTradesHistory", Description: ""}, st.GetTradePairsHistory)
}

func (st *StdioTrarnsport) GetTradePairsHistory(ctx context.Context, req *mcp.CallToolRequest, input GetTradePairsHistoryInput) (*mcp.CallToolResult, GetTradePairsHistoryOutput, error) {
	trades, err := st.dbServer.GetTradesInfo(ctx, input.Seconds)
	if err != nil {
		return nil, GetTradePairsHistoryOutput{}, err
	}

	tradesInfo := GetTradePairsHistoryOutput{trades}
	return nil, tradesInfo, nil
}
