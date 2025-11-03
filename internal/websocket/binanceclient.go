package websocket

import (
	"fmt"
	"mcpbinance/internal/websocket/enum"
	"strings"
	"sync/atomic"
)

type BinanceClient struct {
	msgID int64
}

func (c *BinanceClient) nextID() int64 {
	return atomic.AddInt64(&c.msgID, 1)
}

func (c *BinanceClient) buildSubs(symbols []enum.Symbol, streamTypes []enum.StreamType) []string {
	var streams []string

	for _, symbol := range symbols {
		for _, stmT := range streamTypes {
			streams = append(
				streams, fmt.Sprintf("%s@%s", strings.ToLower(symbol.String()), stmT.String()),
			)
		}
	}

	return streams
}

func (c *BinanceClient) makeSubMsg(streams []string) map[string]any {
	return map[string]any{
		"method": "SUBSCRIBE",
		"params": streams,
		"id":     c.nextID(),
	}
}
