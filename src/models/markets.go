package models

import "github.com/shopspring/decimal"

type L2ConfigModel struct {
	Type                 string `json:"type"`
	CollateralID         string `json:"collateralId"`
	CollateralResolution int64  `json:"collateralResolution"`
	SyntheticID          string `json:"syntheticId"`
	SyntheticResolution  int64  `json:"syntheticResolution"`
}

type MarketModel struct {
	Name                     string        `json:"name"`
	AssetName                string        `json:"assetName"`
	AssetPrecision           int           `json:"assetPrecision"`
	CollateralAssetName      string        `json:"collateralAssetName"`
	CollateralAssetPrecision int           `json:"collateralAssetPrecision"`
	Active                   bool          `json:"active"`
	L2Config                 L2ConfigModel `json:"l2Config"`
	// Note: MarketStats and TradingConfig are typically included in API responses
	// but may not always be present, so they're separate models
}

type RiskFactorConfig struct {
	UpperBound decimal.Decimal `json:"upperBound"`
	RiskFactor decimal.Decimal `json:"riskFactor"`
}

type MarketStatsModel struct {
	DailyVolume        decimal.Decimal `json:"dailyVolume"`
	DailyVolumeBase    decimal.Decimal `json:"dailyVolumeBase"`
	DailyPriceChange   decimal.Decimal `json:"dailyPriceChange"`
	DailyLow           decimal.Decimal `json:"dailyLow"`
	DailyHigh          decimal.Decimal `json:"dailyHigh"`
	LastPrice          decimal.Decimal `json:"lastPrice"`
	AskPrice           decimal.Decimal `json:"askPrice"`
	BidPrice           decimal.Decimal `json:"bidPrice"`
	MarkPrice          decimal.Decimal `json:"markPrice"`
	IndexPrice         decimal.Decimal `json:"indexPrice"`
	FundingRate        decimal.Decimal `json:"fundingRate"`
	NextFundingRate    int64          `json:"nextFundingRate"`
	OpenInterest       decimal.Decimal `json:"openInterest"`
	OpenInterestBase   decimal.Decimal `json:"openInterestBase"`
}

type TradingConfigModel struct {
	MinOrderSize        decimal.Decimal    `json:"minOrderSize"`
	MinOrderSizeChange  decimal.Decimal    `json:"minOrderSizeChange"`
	MinPriceChange      decimal.Decimal    `json:"minPriceChange"`
	MaxMarketOrderValue decimal.Decimal    `json:"maxMarketOrderValue"`
	MaxLimitOrderValue  decimal.Decimal    `json:"maxLimitOrderValue"`
	MaxPositionValue    decimal.Decimal    `json:"maxPositionValue"`
	MaxLeverage         decimal.Decimal    `json:"maxLeverage"`
	MaxNumOrders        int                `json:"maxNumOrders"`
	LimitPriceCap       decimal.Decimal    `json:"limitPriceCap"`
	LimitPriceFloor     decimal.Decimal    `json:"limitPriceFloor"`
	RiskFactorConfig    []RiskFactorConfig `json:"riskFactorConfig"`
}

type CandleType string

const (
	CandleTypeTrades      CandleType = "trades"
	CandleTypeMarkPrices  CandleType = "mark-prices"
	CandleTypeIndexPrices CandleType = "index-prices"
)

type CandleInterval string

const (
	CandleIntervalPT1M  CandleInterval = "PT1M"
	CandleIntervalPT5M  CandleInterval = "PT5M"
	CandleIntervalPT15M CandleInterval = "PT15M"
	CandleIntervalPT30M CandleInterval = "PT30M"
	CandleIntervalPT1H  CandleInterval = "PT1H"
	CandleIntervalPT2H  CandleInterval = "PT2H"
	CandleIntervalPT4H  CandleInterval = "PT4H"
	CandleIntervalP1D   CandleInterval = "P1D"
)

type CandleModel struct {
	Open      decimal.Decimal `json:"open"`
	Low       decimal.Decimal `json:"low"`
	High      decimal.Decimal `json:"high"`
	Close     decimal.Decimal `json:"close"`
	Volume    *decimal.Decimal `json:"volume,omitempty"`
	Timestamp int64           `json:"timestamp"`
}

type FundingRateModel struct {
	Market       string          `json:"market"`
	FundingRate  decimal.Decimal `json:"fundingRate"`
	Timestamp    int64           `json:"timestamp"`
}

type PublicTradeModel struct {
	ID        int             `json:"id"`
	Market    string          `json:"market"`
	Side      string          `json:"side"` // OrderSide
	TradeType string          `json:"tradeType"` // TradeType
	Timestamp int64           `json:"timestamp"`
	Price     decimal.Decimal `json:"price"`
	Qty       decimal.Decimal `json:"qty"`
}

