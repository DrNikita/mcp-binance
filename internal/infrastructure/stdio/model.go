package stdio

import "mcpbinance/internal/entity"

type GetTradePairsHistoryInput struct {
	DateFromTimestamp int    `json:"date_from" jsonschema:"start date as timestamp in milliseconds to get trades"`
	DateToTimestamp   int    `json:"date_to" jsonschema:"end date as timestamp in milliseconds to get trades; zero means current time"`
	Symbol            string `json:"symbol" jsonschema:"symbol to fetch trades; default running stream is BTCUSDT, but can be different if editional streams were launched"`
}

type GetTradePairsHistoryOutput struct {
	Trades []entity.Trade `json:"trades" jsonschema:"info about price changes"`
}

type RunStockMonitoringInput struct {
	Symbols     []string `json:"symbols" jsonschema:"which trade pairs to monitor; avalible options: BTCUSDT,ETHUSDT,BNBUSDT,SOLUSDT,XRPUSDT,ADAUSDT,DOGEUSDT,AVAXUSDT,DOTUSDT,LINKUSDT,LTCUSDT,NEARUSDT,ATOMUSDT,FILUSDT,ETCUSDT,UNIUSDT,MATICUSDT,SUIUSDT"`
	StreamTypes []string `json:"streamTypes" jsonschema:"functions to apply for monitoring from binance.com; avalible options: aggTrade,markPriceUpdate,kline,continuous_kline,24hrTicker,24hrMiniTicker,bookTicker,forceOrder,depthUpdate,compositeIndex,contractInfo,assetIndexUpdate"`
}
