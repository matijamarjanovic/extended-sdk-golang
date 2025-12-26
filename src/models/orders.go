package models

import "github.com/shopspring/decimal"

type OrderType string

const (
	OrderTypeLimit      OrderType = "LIMIT"
	OrderTypeMarket     OrderType = "MARKET"
	OrderTypeConditional OrderType = "CONDITIONAL"
	OrderTypeTpsl       OrderType = "TPSL"
)

type OrderSide string

const (
	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"
)

type TimeInForce string

const (
	TimeInForceGTT TimeInForce = "GTT" // Good till time
	TimeInForceFOK TimeInForce = "FOK" // Fill or kill
	TimeInForceIOC TimeInForce = "IOC" // Immediate or cancel
)

type SelfTradeProtectionLevel string

const (
	SelfTradeProtectionDisabled SelfTradeProtectionLevel = "DISABLED"
	SelfTradeProtectionAccount  SelfTradeProtectionLevel = "ACCOUNT"
	SelfTradeProtectionClient   SelfTradeProtectionLevel = "CLIENT"
)

type OrderStatus string

const (
	OrderStatusUnknown         OrderStatus = "UNKNOWN"
	OrderStatusNew             OrderStatus = "NEW"
	OrderStatusUntriggered     OrderStatus = "UNTRIGGERED"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusCancelled       OrderStatus = "CANCELLED"
	OrderStatusExpired         OrderStatus = "EXPIRED"
	OrderStatusRejected        OrderStatus = "REJECTED"
)

type OrderStatusReason string

const (
	OrderStatusReasonUnknown            OrderStatusReason = "UNKNOWN"
	OrderStatusReasonNone               OrderStatusReason = "NONE"
	OrderStatusReasonUnknownMarket      OrderStatusReason = "UNKNOWN_MARKET"
	OrderStatusReasonDisabledMarket     OrderStatusReason = "DISABLED_MARKET"
	OrderStatusReasonNotEnoughFunds    OrderStatusReason = "NOT_ENOUGH_FUNDS"
	OrderStatusReasonNoLiquidity        OrderStatusReason = "NO_LIQUIDITY"
	OrderStatusReasonInvalidFee         OrderStatusReason = "INVALID_FEE"
	OrderStatusReasonInvalidQty         OrderStatusReason = "INVALID_QTY"
	OrderStatusReasonInvalidPrice       OrderStatusReason = "INVALID_PRICE"
	OrderStatusReasonInvalidValue       OrderStatusReason = "INVALID_VALUE"
	OrderStatusReasonUnknownAccount     OrderStatusReason = "UNKNOWN_ACCOUNT"
	OrderStatusReasonSelfTradeProtection OrderStatusReason = "SELF_TRADE_PROTECTION"
	OrderStatusReasonPostOnlyFailed     OrderStatusReason = "POST_ONLY_FAILED"
	OrderStatusReasonReduceOnlyFailed   OrderStatusReason = "REDUCE_ONLY_FAILED"
	OrderStatusReasonInvalidExpireTime  OrderStatusReason = "INVALID_EXPIRE_TIME"
	OrderStatusReasonPositionTpslConflict OrderStatusReason = "POSITION_TPSL_CONFLICT"
	OrderStatusReasonInvalidLeverage    OrderStatusReason = "INVALID_LEVERAGE"
	OrderStatusReasonPrevOrderNotFound  OrderStatusReason = "PREV_ORDER_NOT_FOUND"
	OrderStatusReasonPrevOrderTriggered OrderStatusReason = "PREV_ORDER_TRIGGERED"
	OrderStatusReasonTpslOtherSideFilled OrderStatusReason = "TPSL_OTHER_SIDE_FILLED"
	OrderStatusReasonPrevOrderConflict  OrderStatusReason = "PREV_ORDER_CONFLICT"
	OrderStatusReasonOrderReplaced       OrderStatusReason = "ORDER_REPLACED"
	OrderStatusReasonPostOnlyMode        OrderStatusReason = "POST_ONLY_MODE"
	OrderStatusReasonReduceOnlyMode     OrderStatusReason = "REDUCE_ONLY_MODE"
	OrderStatusReasonTradingOffMode     OrderStatusReason = "TRADING_OFF_MODE"
)

type TriggerPriceType string

const (
	TriggerPriceTypeUnknown TriggerPriceType = "UNKNOWN"
	TriggerPriceTypeLast   TriggerPriceType = "LAST"
	TriggerPriceTypeMid    TriggerPriceType = "MID"
	TriggerPriceTypeMark   TriggerPriceType = "MARK"
	TriggerPriceTypeIndex  TriggerPriceType = "INDEX"
)

type TriggerDirection string

const (
	TriggerDirectionUnknown TriggerDirection = "UNKNOWN"
	TriggerDirectionUp     TriggerDirection = "UP"
	TriggerDirectionDown   TriggerDirection = "DOWN"
)

type ExecutionPriceType string

const (
	ExecutionPriceTypeUnknown ExecutionPriceType = "UNKNOWN"
	ExecutionPriceTypeLimit   ExecutionPriceType = "LIMIT"
	ExecutionPriceTypeMarket  ExecutionPriceType = "MARKET"
)

type TpSlType string

const (
	TpSlTypeOrder    TpSlType = "ORDER"
	TpSlTypePosition TpSlType = "POSITION"
)

// Cryptographic signature
type Signature struct {
	R string `json:"r"`
	S string `json:"s"`
}

// Settlement information for an order
type Settlement struct {
	Signature          Signature `json:"signature"`
	StarkKey           string    `json:"starkKey"`
	CollateralPosition string    `json:"collateralPosition"`
}

// Conditional trigger settings
type ConditionalTrigger struct {
	TriggerPrice       string             `json:"triggerPrice"`
	TriggerPriceType   TriggerPriceType   `json:"triggerPriceType"`
	Direction          TriggerDirection   `json:"direction"`
	ExecutionPriceType ExecutionPriceType `json:"executionPriceType"`
}

// Take profit or stop loss trigger settings
type TpSlTrigger struct {
	TriggerPrice     string             `json:"triggerPrice"`
	TriggerPriceType TriggerPriceType   `json:"triggerPriceType"`
	Price            string             `json:"price"`
	PriceType        ExecutionPriceType `json:"priceType"`
	Settlement       Settlement         `json:"settlement"`
}

type PerpetualOrderModel struct {
	ID                       string                   `json:"id"`
	Market                   string                   `json:"market"`
	Type                     OrderType                `json:"type"`
	Side                     OrderSide                `json:"side"`
	Qty                      string                   `json:"qty"`
	Price                    string                   `json:"price"`
	TimeInForce              TimeInForce              `json:"timeInForce"`
	ExpiryEpochMillis        int64                    `json:"expiryEpochMillis"`
	Fee                      string                   `json:"fee"`
	Nonce                    string                   `json:"nonce"`
	Settlement               Settlement               `json:"settlement"`
	ReduceOnly               bool                     `json:"reduceOnly"`
	PostOnly                 bool                     `json:"postOnly"`
	SelfTradeProtectionLevel SelfTradeProtectionLevel `json:"selfTradeProtectionLevel"`
	Trigger                  *ConditionalTrigger       `json:"trigger,omitempty"`
	TpSlType                 *TpSlType                `json:"tpSlType,omitempty"`
	TakeProfit               *TpSlTrigger             `json:"takeProfit,omitempty"`
	StopLoss                 *TpSlTrigger             `json:"stopLoss,omitempty"`
	BuilderFee               *string                  `json:"builderFee,omitempty"`
	BuilderID                *int                     `json:"builderId,omitempty"`
	CancelID                 *string                  `json:"cancelId,omitempty"`
}

// TPSL trigger for an open order
type OpenOrderTpslTriggerModel struct {
	TriggerPrice     decimal.Decimal `json:"triggerPrice"`
	TriggerPriceType TriggerPriceType `json:"triggerPriceType"`
	Price            decimal.Decimal `json:"price"`
	PriceType        ExecutionPriceType `json:"priceType"`
	Status           *OrderStatus    `json:"status,omitempty"`
}

// An open order
type OpenOrderModel struct {
	ID                       int                    `json:"id"`
	AccountID                int                    `json:"accountId"`
	ExternalID               string                 `json:"externalId"`
	Market                   string                 `json:"market"`
	Type                     OrderType              `json:"type"`
	Side                     OrderSide              `json:"side"`
	Status                   OrderStatus            `json:"status"`
	StatusReason             *OrderStatusReason     `json:"statusReason,omitempty"`
	Price                    decimal.Decimal        `json:"price"`
	AveragePrice             *decimal.Decimal       `json:"averagePrice,omitempty"`
	Qty                      decimal.Decimal       `json:"qty"`
	FilledQty                *decimal.Decimal       `json:"filledQty,omitempty"`
	ReduceOnly               bool                   `json:"reduceOnly"`
	PostOnly                 bool                   `json:"postOnly"`
	PayedFee                 *decimal.Decimal      `json:"payedFee,omitempty"`
	CreatedTime              int64                  `json:"createdTime"`
	UpdatedTime              int64                  `json:"updatedTime"`
	ExpiryTime               *int64                 `json:"expiryTime,omitempty"`
	TimeInForce              TimeInForce            `json:"timeInForce"`
	TpSlType                 *TpSlType             `json:"tpSlType,omitempty"`
	TakeProfit               *OpenOrderTpslTriggerModel `json:"takeProfit,omitempty"`
	StopLoss                 *OpenOrderTpslTriggerModel `json:"stopLoss,omitempty"`
}

// TpSlTriggerParam represents parameters for a take profit or stop loss trigger
// This matches the Python SDK's OrderTpslTriggerParam structure.
type TpSlTriggerParam struct {
	TriggerPrice     decimal.Decimal
	TriggerPriceType TriggerPriceType
	Price            decimal.Decimal
	PriceType        ExecutionPriceType
}

