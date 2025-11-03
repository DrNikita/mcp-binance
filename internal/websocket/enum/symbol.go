package enum

import (
	"errors"
	"strings"
)

type Symbol int

const (
	BTCUSDT Symbol = iota
	ETCUSDT
)

func NewSymbol(symbol string) (Symbol, error) {
	upperSymbol := strings.ToUpper(symbol)
	switch upperSymbol {
	case "BTCUSDT":
		return BTCUSDT, nil
	case "ETCUSDT":
		return ETCUSDT, nil
	}

	return 0, errors.New("no such a symbol")
}

var symbols = map[Symbol]string{
	BTCUSDT: "BTCUSDT",
	ETCUSDT: "ETCUSDT",
}

func (s Symbol) String() string {
	return symbols[s]
}
