package models

import "github.com/shopspring/decimal"

// EndpointConfig represents the API endpoint configuration
type EndpointConfig struct {
	APIBaseURL string `json:"apiBaseURL"`
}

// StarknetDomain represents the Starknet domain for signing
type StarknetDomain struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	ChainID  string `json:"chainId"`
	Revision string `json:"revision"`
}

// TradingFeeModel represents trading fees for a market
type TradingFeeModel struct {
	Market         string          `json:"market"`
	MakerFeeRate   decimal.Decimal `json:"makerFeeRate"`
	TakerFeeRate   decimal.Decimal `json:"takerFeeRate"`
	BuilderFeeRate decimal.Decimal `json:"builderFeeRate"`
}

// DefaultFees represents the default trading fees
var DefaultFees = TradingFeeModel{
	Market:         "BTC-USD",
	MakerFeeRate:   decimal.NewFromFloat(0.0002), // 2/10000 = 0.0002
	TakerFeeRate:   decimal.NewFromFloat(0.0005), // 5/10000 = 0.0005
	BuilderFeeRate: decimal.NewFromFloat(0),      // 0
}

