package markets

import (
	"context"
	"fmt"

	"github.com/extended-protocol/extended-sdk-golang/src/client"
	"github.com/extended-protocol/extended-sdk-golang/src/models"
)

// Service provides market-related API operations.
type Service struct {
	Base *client.BaseClient
}

// GetMarkets retrieves all available markets from the API
func (s *Service) GetMarkets(ctx context.Context, market []string) ([]models.MarketModel, error) {
	baseURL := s.Base.EndpointConfig().APIBaseURL + "/info/markets"

	if len(market) > 0 {
		baseURL += "?market=" + market[0]
		for i := 1; i < len(market); i++ {
			baseURL += "&market=" + market[i]
		}
	}

	var marketResponse models.MarketResponse
	if err := s.Base.DoRequest(ctx, "GET", baseURL, nil, &marketResponse); err != nil {
		return nil, err
	}

	if marketResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %s", marketResponse.Status)
	}

	return marketResponse.Data, nil
}

// Methods to be implemented:
// - GetMarketsDict (new)
// - GetMarketStatistics (new)
// - GetCandlesHistory (new)
// - GetFundingRatesHistory (new)
// - GetOrderbookSnapshot (new)
//
// Split into multiple files (e.g., markets_candles.go) as code grows

