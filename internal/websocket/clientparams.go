package websocket

import (
	"time"
	// "github.com/sstanqq/hft-arbitrage-toolkit/market-data-collector/internal/domain/types"
)

type ClientParams struct {
	// Exchange     types.ExchangeType
	URL          string
	MsgCap       int
	ReadLimit    int64
	ReconnPeriod time.Duration
	RetryBackoff time.Duration
	PingWait     time.Duration
	PongPeriod   time.Duration
}

func NewWebsocketClient(params ClientParams) (*WebsocketClient, error) {
	baseConfig := &WebsocketConfig{
		// name:         string(params.Exchange),
		url:          params.URL,
		msgCap:       params.MsgCap,
		readLimit:    params.ReadLimit,
		reconnPeriod: params.ReconnPeriod,
		retryBackoff: params.RetryBackoff,
		pingWait:     params.PingWait,
		pongPeriod:   params.PongPeriod,
	}

	var hooks clientHooks
	hooks = &BinanceClient{}

	// switch params.Exchange {
	// case types.ExchangeBinance:
	// 	hooks = &BinanceClient{}
	// 	  case types.ExchangeMEXC:
	// 	    hooks = &MEXCClient{}
	// default:
	// 	return nil, fmt.Errorf("unsupported exchange: %s", params.Exchange)
	// }

	return &WebsocketClient{
		wc:    baseConfig,
		msgCh: make(chan []byte, params.MsgCap),
		subs:  make(map[string]struct{}),
		hooks: hooks,
	}, nil
}
