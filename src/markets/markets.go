package markets

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

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

// GetMarketsDict retrieves all markets and returns them as a dictionary/map keyed by market name
func (s *Service) GetMarketsDict(ctx context.Context) (map[string]models.MarketModel, error) {
	markets, err := s.GetMarkets(ctx, []string{})
	if err != nil {
		return nil, err
	}

	marketsDict := make(map[string]models.MarketModel)
	for _, market := range markets {
		marketsDict[market.Name] = market
	}

	return marketsDict, nil
}

// GetMarketStatistics retrieves market statistics for a specific market
func (s *Service) GetMarketStatistics(ctx context.Context, marketName string) (*models.MarketStatsModel, error) {
	baseUrl, err := s.Base.GetURL(fmt.Sprintf("/info/markets/%s/stats", marketName), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	var statsResponse models.MarketStatsResponse
	if err := s.Base.DoRequest(ctx, "GET", baseUrl, nil, &statsResponse); err != nil {
		return nil, err
	}

	if statsResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", statsResponse.Status)
	}

	return &statsResponse.Data, nil
}

// GetCandlesHistory retrieves candle history for a specific market
func (s *Service) GetCandlesHistory(
	ctx context.Context,
	marketName string,
	candleType models.CandleType,
	interval models.CandleInterval,
	limit *int,
	endTime *time.Time,
) ([]models.CandleModel, error) {
	// Build URL with path parameters
	path := fmt.Sprintf("/info/candles/%s/%s", marketName, candleType)
	
	// Build query parameters
	query := make(url.Values)
	query.Set("interval", string(interval))
	if limit != nil {
		query.Set("limit", strconv.Itoa(*limit))
	}
	if endTime != nil {
		query.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	}

	baseUrl := s.Base.EndpointConfig().APIBaseURL + path
	if len(query) > 0 {
		baseUrl += "?" + query.Encode()
	}

	var candlesResponse models.CandlesResponse
	if err := s.Base.DoRequest(ctx, "GET", baseUrl, nil, &candlesResponse); err != nil {
		return nil, err
	}

	if candlesResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", candlesResponse.Status)
	}

	return candlesResponse.Data, nil
}

// GetFundingRatesHistory retrieves funding rates history for a specific market
func (s *Service) GetFundingRatesHistory(
	ctx context.Context,
	marketName string,
	startTime time.Time,
	endTime time.Time,
) ([]models.FundingRateModel, error) {
	// Build URL with path parameters
	path := fmt.Sprintf("/info/%s/funding", marketName)
	
	// Build query parameters
	query := make(url.Values)
	query.Set("startTime", strconv.FormatInt(startTime.UnixMilli(), 10))
	query.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))

	baseUrl := s.Base.EndpointConfig().APIBaseURL + path + "?" + query.Encode()

	var fundingRatesResponse models.FundingRatesResponse
	if err := s.Base.DoRequest(ctx, "GET", baseUrl, nil, &fundingRatesResponse); err != nil {
		return nil, err
	}

	if fundingRatesResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", fundingRatesResponse.Status)
	}

	return fundingRatesResponse.Data, nil
}

// GetOrderbookSnapshot retrieves the orderbook snapshot for a specific market
func (s *Service) GetOrderbookSnapshot(ctx context.Context, marketName string) (*models.OrderbookUpdateModel, error) {
	baseUrl, err := s.Base.GetURL(fmt.Sprintf("/info/markets/%s/orderbook", marketName), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	var orderbookResponse models.OrderbookResponse
	if err := s.Base.DoRequest(ctx, "GET", baseUrl, nil, &orderbookResponse); err != nil {
		return nil, err
	}

	if orderbookResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", orderbookResponse.Status)
	}

	return &orderbookResponse.Data, nil
}
