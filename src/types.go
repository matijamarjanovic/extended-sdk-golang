package sdk

import (
	"github.com/extended-protocol/extended-sdk-golang/src/client"
	"github.com/extended-protocol/extended-sdk-golang/src/models"
)

// Type aliases for commonly used model types
// These allow using shorter names in the root package without the models. prefix

type EndpointConfig = models.EndpointConfig
type StarknetDomain = models.StarknetDomain
type TradingFeeModel = models.TradingFeeModel

type MarketModel = models.MarketModel
type L2ConfigModel = models.L2ConfigModel

type OrderType = models.OrderType
type OrderSide = models.OrderSide
type TimeInForce = models.TimeInForce
type SelfTradeProtectionLevel = models.SelfTradeProtectionLevel
type PerpetualOrderModel = models.PerpetualOrderModel

type MarketResponse = models.MarketResponse
type FeeResponse = models.FeeResponse
type OrderResponse = models.OrderResponse

// Constants for commonly used enum values
const (
	OrderTypeLimit      = models.OrderTypeLimit
	OrderTypeMarket     = models.OrderTypeMarket
	OrderTypeConditional = models.OrderTypeConditional
	OrderTypeTpsl       = models.OrderTypeTpsl

	OrderSideBuy  = models.OrderSideBuy
	OrderSideSell = models.OrderSideSell

	TimeInForceGTT = models.TimeInForceGTT
	TimeInForceFOK = models.TimeInForceFOK
	TimeInForceIOC = models.TimeInForceIOC

	SelfTradeProtectionDisabled = models.SelfTradeProtectionDisabled
	SelfTradeProtectionAccount  = models.SelfTradeProtectionAccount
	SelfTradeProtectionClient   = models.SelfTradeProtectionClient
)

var DefaultFees = models.DefaultFees

// Type aliases for client package types
type StarkPerpetualAccount = client.StarkPerpetualAccount

// Function aliases for client package functions
var NewStarkPerpetualAccount = client.NewStarkPerpetualAccount
var GetOrderHash = client.GetOrderHash
var SignMessage = client.SignMessage

