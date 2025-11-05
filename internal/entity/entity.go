package entity

type Trade struct {
	EventType        string `json:"e" bson:"e"`
	EventTime        int64  `json:"E" bson:"E"`
	Symbol           string `json:"s" bson:"s"`
	AggregateTradeID int64  `json:"a" bson:"a"`
	Price            string `json:"p" bson:"p"`
	Quantity         string `json:"q" bson:"q"`
	FirstTradeID     int64  `json:"f" bson:"f"`
	LastTradeID      int64  `json:"l" bson:"l"`
	TradeTime        int64  `json:"T" bson:"T"`
	MarketMaker      bool   `json:"m" bson:"m"`
}
