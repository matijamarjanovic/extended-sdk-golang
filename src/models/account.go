package models

import "github.com/shopspring/decimal"

// The account balance
type BalanceModel struct {
	CollateralName           string          `json:"collateralName"`
	Balance                  decimal.Decimal `json:"balance"`
	Equity                   decimal.Decimal `json:"equity"`
	AvailableForTrade        decimal.Decimal `json:"availableForTrade"`
	AvailableForWithdrawal    decimal.Decimal `json:"availableForWithdrawal"`
	UnrealisedPnl            decimal.Decimal `json:"unrealisedPnl"`
	InitialMargin            decimal.Decimal `json:"initialMargin"`
	MarginRatio              decimal.Decimal `json:"marginRatio"`
	UpdatedTime              int64           `json:"updatedTime"`
}

type PositionSide string

const (
	PositionSideLong  PositionSide = "LONG"
	PositionSideShort PositionSide = "SHORT"
)

type PositionStatus string

const (
	PositionStatusOpened PositionStatus = "OPENED"
	PositionStatusClosed PositionStatus = "CLOSED"
)

type PositionModel struct {
	ID              int             `json:"id"`
	AccountID       int             `json:"accountId"`
	Market          string          `json:"market"`
	Status          PositionStatus  `json:"status"`
	Side            PositionSide    `json:"side"`
	Leverage        decimal.Decimal `json:"leverage"`
	Size            decimal.Decimal `json:"size"`
	Value           decimal.Decimal `json:"value"`
	OpenPrice       decimal.Decimal `json:"openPrice"`
	MarkPrice       decimal.Decimal `json:"markPrice"`
	LiquidationPrice *decimal.Decimal `json:"liquidationPrice,omitempty"`
	UnrealisedPnl   decimal.Decimal `json:"unrealisedPnl"`
	RealisedPnl     decimal.Decimal `json:"realisedPnl"`
	TpPrice         *decimal.Decimal `json:"tpPrice,omitempty"`
	SlPrice         *decimal.Decimal `json:"slPrice,omitempty"`
	Adl             *int             `json:"adl,omitempty"`
	CreatedAt       int64           `json:"createdAt"`
	UpdatedAt       int64           `json:"updatedAt"`
}

type ExitType string

const (
	ExitTypeTrade       ExitType = "TRADE"
	ExitTypeLiquidation ExitType = "LIQUIDATION"
	ExitTypeAdl         ExitType = "ADL"
)

type RealisedPnlBreakdownModel struct {
	TradePnl    decimal.Decimal `json:"tradePnl"`
	FundingFees decimal.Decimal `json:"fundingFees"`
	OpenFees    decimal.Decimal `json:"openFees"`
	CloseFees   decimal.Decimal `json:"closeFees"`
}

type PositionHistoryModel struct {
	ID                  int                     `json:"id"`
	AccountID           int                     `json:"accountId"`
	Market              string                  `json:"market"`
	Side                PositionSide            `json:"side"`
	Size                decimal.Decimal         `json:"size"`
	MaxPositionSize     decimal.Decimal         `json:"maxPositionSize"`
	Leverage            decimal.Decimal         `json:"leverage"`
	OpenPrice           decimal.Decimal        `json:"openPrice"`
	ExitPrice           *decimal.Decimal       `json:"exitPrice,omitempty"`
	RealisedPnl         decimal.Decimal         `json:"realisedPnl"`
	RealisedPnlBreakdown RealisedPnlBreakdownModel `json:"realisedPnlBreakdown"`
	CreatedTime         int64                   `json:"createdTime"`
	ExitType            *ExitType               `json:"exitType,omitempty"`
	ClosedTime          *int64                  `json:"closedTime,omitempty"`
}

type AccountModel struct {
	ID                   int      `json:"id"`
	Description          string   `json:"description"`
	AccountIndex         int      `json:"accountIndex"`
	Status               string   `json:"status"`
	L2Key                string   `json:"l2Key"`
	L2Vault              int      `json:"l2Vault"`
	BridgeStarknetAddress *string  `json:"bridgeStarknetAddress,omitempty"`
	APIKeys              []string `json:"apiKeys,omitempty"`
}

type ClientModel struct {
	ID                    int     `json:"id"`
	EvmWalletAddress      *string `json:"evmWalletAddress,omitempty"`
	StarknetWalletAddress *string `json:"starknetWalletAddress,omitempty"`
	ReferralLinkCode      *string `json:"referralLinkCode,omitempty"`
}

// Leverage for a specific market
type AccountLeverage struct {
	Market   string          `json:"market"`
	Leverage  decimal.Decimal `json:"leverage"`
}

type TradeType string

const (
	TradeTypeTrade       TradeType = "TRADE"
	TradeTypeLiquidation TradeType = "LIQUIDATION"
	TradeTypeDeleverage  TradeType = "DELEVERAGE"
)

type AccountTradeModel struct {
	ID          int             `json:"id"`
	AccountID   int             `json:"accountId"`
	Market      string          `json:"market"`
	OrderID     int             `json:"orderId"`
	Side        string          `json:"side"` // OrderSide
	Price       decimal.Decimal `json:"price"`
	Qty         decimal.Decimal `json:"qty"`
	Value       decimal.Decimal `json:"value"`
	Fee         decimal.Decimal `json:"fee"`
	IsTaker     bool            `json:"isTaker"`
	TradeType   TradeType       `json:"tradeType"`
	CreatedTime int64           `json:"createdTime"`
}

type AssetOperationType string

const (
	AssetOperationTypeClaim          AssetOperationType = "CLAIM"
	AssetOperationTypeDeposit       AssetOperationType = "DEPOSIT"
	AssetOperationTypeFastWithdrawal AssetOperationType = "FAST_WITHDRAWAL"
	AssetOperationTypeSlowWithdrawal AssetOperationType = "SLOW_WITHDRAWAL"
	AssetOperationTypeTransfer       AssetOperationType = "TRANSFER"
)

type AssetOperationStatus string

const (
	AssetOperationStatusUnknown        AssetOperationStatus = "UNKNOWN"
	AssetOperationStatusCreated        AssetOperationStatus = "CREATED"
	AssetOperationStatusInProgress     AssetOperationStatus = "IN_PROGRESS"
	AssetOperationStatusRejected       AssetOperationStatus = "REJECTED"
	AssetOperationStatusReadyForClaim  AssetOperationStatus = "READY_FOR_CLAIM"
	AssetOperationStatusCompleted      AssetOperationStatus = "COMPLETED"
)

type AssetOperationModel struct {
	ID                   string               `json:"id"`
	Type                 AssetOperationType   `json:"type"`
	Status               AssetOperationStatus `json:"status"`
	Amount               decimal.Decimal     `json:"amount"`
	Fee                  decimal.Decimal     `json:"fee"`
	Asset                int                 `json:"asset"`
	Time                 int64               `json:"time"`
	AccountID            int                 `json:"accountId"`
	CounterpartyAccountID *int                `json:"counterpartyAccountId,omitempty"`
	TransactionHash      *string             `json:"transactionHash,omitempty"`
}

type ChainConfig struct {
	Chain           string `json:"chain"`
	ContractAddress string `json:"contractAddress"`
}

type BridgesConfig struct {
	Chains []ChainConfig `json:"chains"`
}

type Quote struct {
	ID  string          `json:"id"`
	Fee decimal.Decimal `json:"fee"`
}

