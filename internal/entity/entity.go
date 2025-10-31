package entity

type Trade struct {
	EventType        string `bson:"e"`
	EventTime        int64  `bson:"E"`
	Symbol           string `bson:"s"`
	AggregateTradeID int64  `bson:"a"`
	Price            string `bson:"p"`
	Quantity         string `bson:"q"`
	FirstTradeID     int64  `bson:"f"`
	LastTradeID      int64  `bson:"l"`
	TradeTime        int64  `bson:"T"`
	MarketMaker      bool   `bson:"m"`
}
