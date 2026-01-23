package models

import "github.com/shopspring/decimal"

// StreamDataType represents the type of stream message
type StreamDataType string

const (
	StreamDataTypeUnknown    StreamDataType = "UNKNOWN"
	StreamDataTypeBalance    StreamDataType = "BALANCE"
	StreamDataTypeDelta      StreamDataType = "DELTA"
	StreamDataTypeDeposit    StreamDataType = "DEPOSIT"
	StreamDataTypeOrder      StreamDataType = "ORDER"
	StreamDataTypePosition   StreamDataType = "POSITION"
	StreamDataTypeSnapshot   StreamDataType = "SNAPSHOT"
	StreamDataTypeTrade      StreamDataType = "TRADE"
	StreamDataTypeTransfer   StreamDataType = "TRANSFER"
	StreamDataTypeWithdrawal StreamDataType = "WITHDRAWAL"
)

// WrappedStreamResponse wraps stream messages with metadata
type WrappedStreamResponse struct {
	Type  *StreamDataType `json:"type,omitempty"`
	Data  interface{}     `json:"data,omitempty"`
	Error *string         `json:"error,omitempty"`
	TS    int64           `json:"ts"`
	Seq   int64           `json:"seq"`
}

// StreamPublicTradeModel represents a public trade from the stream with short JSON tags
// (streaming API uses "i", "m", "S", "tT", "T", "p", "q" instead of full names)
type StreamPublicTradeModel struct {
	ID        int64           `json:"i"`
	Market    string          `json:"m"`
	Side      OrderSide       `json:"S"`
	TradeType TradeType       `json:"tT"`
	Timestamp int64           `json:"T"`
	Price     decimal.Decimal `json:"p"`
	Qty       decimal.Decimal `json:"q"`
}

// AccountStreamDataModel represents account update data from the stream
type AccountStreamDataModel struct {
	Orders    []OpenOrderModel    `json:"orders,omitempty"`
	Positions []PositionModel     `json:"positions,omitempty"`
	Trades    []AccountTradeModel `json:"trades,omitempty"`
	Balance   *BalanceModel       `json:"balance,omitempty"`
}

// MarkPriceModel represents a mark price update from the stream
type MarkPriceModel struct {
	Market    string          `json:"m"`
	Price     decimal.Decimal `json:"p"`
	Timestamp int64           `json:"ts"`
}

// IndexPriceModel represents an index price update from the stream
type IndexPriceModel struct {
	Market    string          `json:"m"`
	Price     decimal.Decimal `json:"p"`
	Timestamp int64           `json:"ts"`
}
