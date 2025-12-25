package markets

import (
	"github.com/extended-protocol/extended-sdk-golang/src"
	"github.com/extended-protocol/extended-sdk-golang/src/models"
)

// Service provides market-related API operations.
// It holds a reference to the main client to access shared infrastructure.
type Service struct {
	client *sdk.Client // Reference to main client
}

// Methods to be implemented:
// - GetMarkets (move from api_client.go)
// - GetMarketsDict (new)
// - GetMarketStatistics (new)
// - GetCandlesHistory (new)
// - GetFundingRatesHistory (new)
// - GetOrderbookSnapshot (new)
//
// Split into multiple files (e.g., markets_candles.go) as code grows

