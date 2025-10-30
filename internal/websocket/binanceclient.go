package websocket

import (
	"fmt"
	"strings"
	"sync/atomic"
)

type BinanceClient struct {
	msgID int64
}

func (c *BinanceClient) nextID() int64 {
	return atomic.AddInt64(&c.msgID, 1)
}

func (c *BinanceClient) buildSubs(symbols, streamTypes []string) []string {
	var streams []string

	for _, sym := range symbols {
		for _, stmT := range streamTypes {
			streams = append(
				streams, fmt.Sprintf("%s@%s", strings.ToLower(sym), stmT),
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
