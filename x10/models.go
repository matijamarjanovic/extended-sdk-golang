// Type aliases for x10/models types used as service arguments; import only x10 to use the SDK.
package x10

import (
	"github.com/matijamarjanovic/extended-sdk-golang/x10/client"
	"github.com/matijamarjanovic/extended-sdk-golang/x10/models"
)

type MarketModel = models.MarketModel
type L2ConfigModel = models.L2ConfigModel
type StarknetDomain = models.StarknetDomain
type StarkPerpetualAccount = client.StarkPerpetualAccount
type OrderSide = models.OrderSide
type OrderType = models.OrderType
type TimeInForce = models.TimeInForce
type SelfTradeProtectionLevel = models.SelfTradeProtectionLevel
type TpSlType = models.TpSlType
type TpSlTriggerParam = models.TpSlTriggerParam
type PositionSide = models.PositionSide
type TradeType = models.TradeType
type AssetOperationType = models.AssetOperationType
type AssetOperationStatus = models.AssetOperationStatus
type CandleType = models.CandleType
type CandleInterval = models.CandleInterval
type TriggerPriceType = models.TriggerPriceType
type ExecutionPriceType = models.ExecutionPriceType
type TriggerDirection = models.TriggerDirection

type OrderbookUpdateModel = models.OrderbookUpdateModel
type StreamPublicTradeModel = models.StreamPublicTradeModel
type FundingRateModel = models.FundingRateModel
type CandleModel = models.CandleModel
type AccountStreamDataModel = models.AccountStreamDataModel
type MarkPriceModel = models.MarkPriceModel
type IndexPriceModel = models.IndexPriceModel

const (
	PositionSideAll   = models.PositionSideAll
	PositionSideLong  = models.PositionSideLong
	PositionSideShort = models.PositionSideShort
)

const (
	OrderTypeAll         = models.OrderTypeAll
	OrderTypeLimit       = models.OrderTypeLimit
	OrderTypeMarket      = models.OrderTypeMarket
	OrderTypeConditional = models.OrderTypeConditional
	OrderTypeTpsl        = models.OrderTypeTpsl
)

const (
	OrderSideAll  = models.OrderSideAll
	OrderSideBuy  = models.OrderSideBuy
	OrderSideSell = models.OrderSideSell
)

const (
	TimeInForceGTT = models.TimeInForceGTT
	TimeInForceFOK = models.TimeInForceFOK
	TimeInForceIOC = models.TimeInForceIOC
)

const (
	SelfTradeProtectionDisabled = models.SelfTradeProtectionDisabled
	SelfTradeProtectionAccount  = models.SelfTradeProtectionAccount
	SelfTradeProtectionClient   = models.SelfTradeProtectionClient
)

const (
	TradeTypeAll         = models.TradeTypeAll
	TradeTypeTrade       = models.TradeTypeTrade
	TradeTypeLiquidation = models.TradeTypeLiquidation
	TradeTypeDeleverage  = models.TradeTypeDeleverage
)

const (
	TpSlTypeOrder    = models.TpSlTypeOrder
	TpSlTypePosition = models.TpSlTypePosition
)

const (
	CandleTypeTrades      = models.CandleTypeTrades
	CandleTypeMarkPrices  = models.CandleTypeMarkPrices
	CandleTypeIndexPrices = models.CandleTypeIndexPrices
)

const (
	CandleIntervalPT1M  = models.CandleIntervalPT1M
	CandleIntervalPT5M  = models.CandleIntervalPT5M
	CandleIntervalPT15M = models.CandleIntervalPT15M
	CandleIntervalPT30M = models.CandleIntervalPT30M
	CandleIntervalPT1H  = models.CandleIntervalPT1H
	CandleIntervalPT2H  = models.CandleIntervalPT2H
	CandleIntervalPT4H  = models.CandleIntervalPT4H
	CandleIntervalP1D   = models.CandleIntervalP1D
)

const (
	AssetOperationTypeClaim      = models.AssetOperationTypeClaim
	AssetOperationTypeDeposit    = models.AssetOperationTypeDeposit
	AssetOperationTypeWithdrawal = models.AssetOperationTypeWithdrawal
	AssetOperationTypeTransfer   = models.AssetOperationTypeTransfer
)

const (
	AssetOperationStatusCreated    = models.AssetOperationStatusCreated
	AssetOperationStatusInProgress = models.AssetOperationStatusInProgress
	AssetOperationStatusRejected   = models.AssetOperationStatusRejected
	AssetOperationStatusCompleted  = models.AssetOperationStatusCompleted
)

const (
	TriggerPriceTypeUnknown = models.TriggerPriceTypeUnknown
	TriggerPriceTypeLast    = models.TriggerPriceTypeLast
	TriggerPriceTypeMid     = models.TriggerPriceTypeMid
	TriggerPriceTypeMark    = models.TriggerPriceTypeMark
	TriggerPriceTypeIndex   = models.TriggerPriceTypeIndex
)

const (
	ExecutionPriceTypeUnknown = models.ExecutionPriceTypeUnknown
	ExecutionPriceTypeLimit   = models.ExecutionPriceTypeLimit
	ExecutionPriceTypeMarket  = models.ExecutionPriceTypeMarket
)

const (
	TriggerDirectionUnknown = models.TriggerDirectionUnknown
	TriggerDirectionUp      = models.TriggerDirectionUp
	TriggerDirectionDown    = models.TriggerDirectionDown
)
