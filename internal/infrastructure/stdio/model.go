package stdio

import "mcpbinance/internal/entity"

// TODO: add symbol param for search
type GetTradePairsHistoryInput struct {
	Seconds int `json:"seconds" jsonschema:"the number of seconds from the current date for which information is obtained"`
}

type GetTradePairsHistoryOutput struct {
	Trades []entity.Trade `json:"trades" jsonschema:"info about price changes"`
}
