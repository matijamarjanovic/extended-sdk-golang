package models

// APIResponse represents a generic API response with status and data
// T can be any type: a single model, a slice of models, or a struct
type APIResponse[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data"`
}

// OrderSubmissionData represents the data returned after order submission
type OrderSubmissionData struct {
	OrderID    uint   `json:"id"`
	ExternalID string `json:"externalId"`
}

// Type aliases for common response types (for convenience and clarity)
type MarketResponse = APIResponse[[]MarketModel]
type FeeResponse = APIResponse[[]TradingFeeModel]
type OrderResponse = APIResponse[OrderSubmissionData]
type BalanceResponse = APIResponse[BalanceModel]
type PositionsResponse = APIResponse[[]PositionModel]
type PositionsHistoryResponse = APIResponse[[]PositionHistoryModel]
type OpenOrdersResponse = APIResponse[[]OpenOrderModel]
type OrdersHistoryResponse = APIResponse[[]OpenOrderModel]
type TradesResponse = APIResponse[[]AccountTradeModel]
type AccountResponse = APIResponse[AccountModel]
type ClientResponse = APIResponse[ClientModel]
type LeverageResponse = APIResponse[[]AccountLeverage]
type AssetOperationsResponse = APIResponse[[]AssetOperationModel]
type BridgesConfigResponse = APIResponse[BridgesConfig]
type QuoteResponse = APIResponse[Quote]
type MarketStatsResponse = APIResponse[MarketStatsModel]
type CandlesResponse = APIResponse[[]CandleModel]
type FundingRatesResponse = APIResponse[[]FundingRateModel]
type OrderbookResponse = APIResponse[OrderbookUpdateModel]

// EmptyResponse represents an empty API response (for operations that don't return data)
type EmptyResponse struct {
	Status string `json:"status"`
}

