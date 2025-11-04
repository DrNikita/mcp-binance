package enum

import "fmt"

func CreateStreamParams(strSymbols, strStreamTypes []string) ([]Symbol, []StreamType, error) {
	symbols := make([]Symbol, 0)
	streamTypes := make([]StreamType, 0)
	for _, symbol := range strSymbols {
		eSymbol, err := NewSymbol(symbol)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %w", "failed to run symbols monitoring", err)
		}
		symbols = append(symbols, eSymbol)
	}
	for _, streamType := range strStreamTypes {
		eStreamType, err := NewStreamType(streamType)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %w", "failed to run symbols monitoring", err)
		}
		streamTypes = append(streamTypes, eStreamType)
	}

	return symbols, streamTypes, nil
}
