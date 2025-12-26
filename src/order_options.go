package sdk

import (
	"time"

	"github.com/extended-protocol/extended-sdk-golang/src/models"
	"github.com/extended-protocol/extended-sdk-golang/src/services"
	"github.com/shopspring/decimal"
)

// The functions below are convenience wrappers around services package functions,
// allowing users to use sdk.WithPostOnly() instead of services.WithPostOnly() (only import sdk package).

// PlaceOrderOption is an alias for services.PlaceOrderOption to make it easier to use from the root package
type PlaceOrderOption = services.PlaceOrderOption

func WithPostOnly(postOnly bool) PlaceOrderOption {
	return services.WithPostOnly(postOnly)
}

func WithReduceOnly(reduceOnly bool) PlaceOrderOption {
	return services.WithReduceOnly(reduceOnly)
}

func WithExpireTime(expireTime time.Time) PlaceOrderOption {
	return services.WithExpireTime(expireTime)
}

func WithPreviousOrderExternalID(id string) PlaceOrderOption {
	return services.WithPreviousOrderExternalID(id)
}

func WithOrderExternalID(id string) PlaceOrderOption {
	return services.WithOrderExternalID(id)
}

func WithBuilderFee(fee decimal.Decimal) PlaceOrderOption {
	return services.WithBuilderFee(fee)
}

func WithBuilderID(id int) PlaceOrderOption {
	return services.WithBuilderID(id)
}

func WithTpSlType(tpSlType models.TpSlType) PlaceOrderOption {
	return services.WithTpSlType(tpSlType)
}

func WithTakeProfit(trigger models.TpSlTriggerParam) PlaceOrderOption {
	return services.WithTakeProfit(trigger)
}

func WithStopLoss(trigger models.TpSlTriggerParam) PlaceOrderOption {
	return services.WithStopLoss(trigger)
}

func WithNonce(nonce int) PlaceOrderOption {
	return services.WithNonce(nonce)
}
