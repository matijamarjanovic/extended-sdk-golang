package services

import (
	"time"

	"github.com/extended-protocol/extended-sdk-golang/src/models"
	"github.com/shopspring/decimal"
)

// PlaceOrderConfig holds all parameters for placing an order (required and optional)
type PlaceOrderConfig struct {
	// Required fields
	Market                   models.MarketModel
	SyntheticAmount          decimal.Decimal
	Price                    decimal.Decimal
	Side                     models.OrderSide
	Type                     models.OrderType
	TimeInForce              models.TimeInForce
	SelfTradeProtectionLevel models.SelfTradeProtectionLevel

	// Optional fields with defaults
	Nonce                    *int // nil means auto-generate
	PostOnly                bool
	ReduceOnly              bool
	ExpireTime              *time.Time // nil means default (1 hour from now)
	PreviousOrderExternalID *string
	OrderExternalID         *string
	BuilderFee              *decimal.Decimal
	BuilderID               *int
	TpSlType                *models.TpSlType
	TakeProfit              *models.TpSlTriggerParam
	StopLoss                *models.TpSlTriggerParam
}

type PlaceOrderOption func(*PlaceOrderConfig)

func WithPostOnly(postOnly bool) PlaceOrderOption {
	return func(c *PlaceOrderConfig) {
		c.PostOnly = postOnly
	}
}

func WithReduceOnly(reduceOnly bool) PlaceOrderOption {
	return func(c *PlaceOrderConfig) {
		c.ReduceOnly = reduceOnly
	}
}

// WithExpireTime sets the expiration time for the order
func WithExpireTime(expireTime time.Time) PlaceOrderOption {
	return func(c *PlaceOrderConfig) {
		c.ExpireTime = &expireTime
	}
}

// WithPreviousOrderExternalID sets the previous order external ID (for order replacement)
func WithPreviousOrderExternalID(id string) PlaceOrderOption {
	return func(c *PlaceOrderConfig) {
		c.PreviousOrderExternalID = &id
	}
}

// WithOrderExternalID sets a custom external ID for the order
func WithOrderExternalID(id string) PlaceOrderOption {
	return func(c *PlaceOrderConfig) {
		c.OrderExternalID = &id
	}
}

// WithBuilderFee sets the builder fee for the order
func WithBuilderFee(fee decimal.Decimal) PlaceOrderOption {
	return func(c *PlaceOrderConfig) {
		c.BuilderFee = &fee
	}
}

// WithBuilderID sets the builder ID for the order
func WithBuilderID(id int) PlaceOrderOption {
	return func(c *PlaceOrderConfig) {
		c.BuilderID = &id
	}
}

// WithTpSlType sets the TPSL type for the order
func WithTpSlType(tpSlType models.TpSlType) PlaceOrderOption {
	return func(c *PlaceOrderConfig) {
		c.TpSlType = &tpSlType
	}
}

// WithTakeProfit sets the take profit trigger parameters
func WithTakeProfit(trigger models.TpSlTriggerParam) PlaceOrderOption {
	return func(c *PlaceOrderConfig) {
		c.TakeProfit = &trigger
	}
}

// WithStopLoss sets the stop loss trigger parameters
func WithStopLoss(trigger models.TpSlTriggerParam) PlaceOrderOption {
	return func(c *PlaceOrderConfig) {
		c.StopLoss = &trigger
	}
}

// WithNonce sets a custom nonce for the order. If not provided, nonce will be auto-generated.
func WithNonce(nonce int) PlaceOrderOption {
	return func(c *PlaceOrderConfig) {
		c.Nonce = &nonce
	}
}

// buildPlaceOrderConfig builds a PlaceOrderConfig from required parameters and options.
func buildPlaceOrderConfig(
	market models.MarketModel,
	syntheticAmount decimal.Decimal,
	price decimal.Decimal,
	side models.OrderSide,
	orderType models.OrderType,
	timeInForce models.TimeInForce,
	selfTradeProtectionLevel models.SelfTradeProtectionLevel,
	opts ...PlaceOrderOption,
) *PlaceOrderConfig {
	// Initialize config with required parameters
	config := &PlaceOrderConfig{
		Market:                   market,
		SyntheticAmount:          syntheticAmount,
		Price:                    price,
		Side:                     side,
		Type:                     orderType,
		TimeInForce:              timeInForce,
		SelfTradeProtectionLevel: selfTradeProtectionLevel,
		// Defaults
		PostOnly:   false,
		ReduceOnly: false,
		Nonce:      nil, // Auto-generate if not provided via option
	}

	// Apply options
	for _, opt := range opts {
		opt(config)
	}

	return config
}

