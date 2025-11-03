package enum

import "errors"

type StreamType int

const (
	AggTrade StreamType = iota
	MarkPriceUpdate
	Kline
	ContinuousKline
	Hr24Ticker
	Hr24MiniTicker
	BookTicker
	ForceOrder
	DepthUpdate
	CompositeIndex
	ContractInfo
	AssetIndexUpdate
)

var streamTypes = map[StreamType]string{
	AggTrade:         "aggTrade",
	MarkPriceUpdate:  "markPriceUpdate",
	Kline:            "kline",
	ContinuousKline:  "continuous_kline",
	Hr24Ticker:       "24hrTicker",
	Hr24MiniTicker:   "24hrMiniTicker",
	BookTicker:       "bookTicker",
	ForceOrder:       "forceOrder",
	DepthUpdate:      "depthUpdate",
	CompositeIndex:   "compositeIndex",
	ContractInfo:     "contractInfo",
	AssetIndexUpdate: "assetIndexUpdate",
}

func NewStreamType(streamType string) (StreamType, error) {
	switch streamType {
	case "aggTrade":
		return AggTrade, nil
	case "markPriceUpdate":
		return MarkPriceUpdate, nil
	case "kline":
		return Kline, nil
	case "continuous_kline":
		return ContinuousKline, nil
	case "24hrTicker":
		return Hr24Ticker, nil
	case "24hrMiniTicker":
		return Hr24MiniTicker, nil
	case "bookTicker":
		return BookTicker, nil
	case "forceOrder":
		return ForceOrder, nil
	case "depthUpdate":
		return DepthUpdate, nil
	case "compositeIndex":
		return CompositeIndex, nil
	case "contractInfo":
		return ContractInfo, nil
	case "assetIndexUpdate":
		return AssetIndexUpdate, nil
	}

	return 0, errors.New("no such a stream type")
}

func (st StreamType) String() string {
	return streamTypes[st]
}
