package stdio

import "mcpbinance/internal/entity"

// TODO: add symbol param for search
type GetTradePairsHistoryInput struct {
	Seconds int    `json:"seconds" jsonschema:"the number of seconds from the current date for which information is obtained"`
	Symbol  string `json:"symbol" jsonschema:"symbol to fetch price changes history"`
}

type GetTradePairsHistoryOutput struct {
	Trades []entity.Trade `json:"trades" jsonschema:"info about price changes"`
}

type RunStockMonitoringInput struct {
	Symbols     []string `json:"symbols" jsonschema:"which trade pairs to monitor; avalible options: BTCUSDT,ETHUSDT,BNBUSDT,SOLUSDT,XRPUSDT,ADAUSDT,DOGEUSDT,AVAXUSDT,DOTUSDT,LINKUSDT,LTCUSDT,NEARUSDT,ATOMUSDT,FILUSDT,ETCUSDT,UNIUSDT,MATICUSDT,SUIUSDT"`
	StreamTypes []string `json:"streamTypes" jsonschema:"functions to apply for monitoring from binance.com; avalible options: aggTrade,markPriceUpdate,kline,continuous_kline,24hrTicker,24hrMiniTicker,bookTicker,forceOrder,depthUpdate,compositeIndex,contractInfo,assetIndexUpdate"`
}
