package enum

import (
	"errors"
	"strings"
)

type Symbol int

const (
	BTCUSDT = iota
	ETHUSDT
	BNBUSDT
	SOLUSDT
	XRPUSDT
	ADAUSDT
	DOGEUSDT
	AVAXUSDT
	DOTUSDT
	LINKUSDT
	LTCUSDT
	NEARUSDT
	ATOMUSDT
	FILUSDT
	ETCUSDT
	UNIUSDT
	MATICUSDT
	SUIUSDT
)

func NewSymbol(symbol string) (Symbol, error) {
	upperSymbol := strings.ToUpper(symbol)
	switch upperSymbol {
	case "BTCUSDT":
		return BTCUSDT, nil
	case "ETHUSDT":
		return ETHUSDT, nil
	case "BNBUSDT":
		return BNBUSDT, nil
	case "SOLUSDT":
		return SOLUSDT, nil
	case "XRPUSDT":
		return XRPUSDT, nil
	case "ADAUSDT":
		return ADAUSDT, nil
	case "DOGEUSDT":
		return DOGEUSDT, nil
	case "AVAXUSDT":
		return AVAXUSDT, nil
	case "DOTUSDT":
		return DOTUSDT, nil
	case "LINKUSDT":
		return LINKUSDT, nil
	case "LTCUSDT":
		return LTCUSDT, nil
	case "NEARUSDT":
		return NEARUSDT, nil
	case "ATOMUSDT":
		return ATOMUSDT, nil
	case "FILUSDT":
		return FILUSDT, nil
	case "ETCUSDT":
		return ETCUSDT, nil
	case "UNIUSDT":
		return UNIUSDT, nil
	case "MATICUSDT":
		return MATICUSDT, nil
	case "SUIUSDT":
		return SUIUSDT, nil
	}

	return 0, errors.New("no such a symbol")
}

var symbols = map[Symbol]string{
	BTCUSDT:   "BTCUSDT",
	ETHUSDT:   "ETHUSDT",
	BNBUSDT:   "BNBUSDT",
	SOLUSDT:   "SOLUSDT",
	XRPUSDT:   "XRPUSDT",
	ADAUSDT:   "ADAUSDT",
	DOGEUSDT:  "DOGEUSDT",
	AVAXUSDT:  "AVAXUSDT",
	DOTUSDT:   "DOTUSDT",
	LINKUSDT:  "LINKUSDT",
	LTCUSDT:   "LTCUSDT",
	NEARUSDT:  "NEARUSDT",
	ATOMUSDT:  "ATOMUSDT",
	FILUSDT:   "FILUSDT",
	ETCUSDT:   "ETCUSDT",
	UNIUSDT:   "UNIUSDT",
	MATICUSDT: "MATICUSDT",
	SUIUSDT:   "SUIUSDT",
}

func (s Symbol) String() string {
	return symbols[s]
}
