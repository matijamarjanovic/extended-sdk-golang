package account

import (
	"context"
	"fmt"

	"github.com/extended-protocol/extended-sdk-golang/src/client"
	"github.com/extended-protocol/extended-sdk-golang/src/models"
	"github.com/shopspring/decimal"
)

// Service provides account-related API operations.
type Service struct {
	Base *client.BaseClient
}

// GetAccount retrieves account information
func (s *Service) GetAccount(ctx context.Context) (*models.AccountModel, error) {
	baseUrl, err := s.Base.GetURL("/user/account/info", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	var accountResponse models.AccountResponse
	if err := s.Base.DoRequest(ctx, "GET", baseUrl, nil, &accountResponse); err != nil {
		return nil, err
	}

	if accountResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", accountResponse.Status)
	}

	return &accountResponse.Data, nil
}

// GetClient retrieves client information
func (s *Service) GetClient(ctx context.Context) (*models.ClientModel, error) {
	baseUrl, err := s.Base.GetURL("/user/client/info", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	var clientResponse models.ClientResponse
	if err := s.Base.DoRequest(ctx, "GET", baseUrl, nil, &clientResponse); err != nil {
		return nil, err
	}

	if clientResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", clientResponse.Status)
	}

	return &clientResponse.Data, nil
}

// GetBalance retrieves account balance information
func (s *Service) GetBalance(ctx context.Context) (*models.BalanceModel, error) {
	baseUrl, err := s.Base.GetURL("/user/balance", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	var balanceResponse models.BalanceResponse
	if err := s.Base.DoRequest(ctx, "GET", baseUrl, nil, &balanceResponse); err != nil {
		return nil, err
	}

	if balanceResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", balanceResponse.Status)
	}

	return &balanceResponse.Data, nil
}

// GetMarketFee retrieves current trading fees for a specific market
func (s *Service) GetMarketFee(ctx context.Context, market string) ([]models.TradingFeeModel, error) {
	baseUrl, err := s.Base.GetURL("/user/fees", map[string]string{"market": market})
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	var feeResponse models.FeeResponse
	if err := s.Base.DoRequest(ctx, "GET", baseUrl, nil, &feeResponse); err != nil {
		return nil, err
	}

	if feeResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", feeResponse.Status)
	}

	return feeResponse.Data, nil
}

// GetFees retrieves trading fees for specified markets (matches Python SDK signature)
func (s *Service) GetFees(ctx context.Context, marketNames []string, builderID *int) ([]models.TradingFeeModel, error) {
	// Build URL with query parameters
	baseUrl := s.Base.EndpointConfig().APIBaseURL + "/user/fees"
	
	// Build query string manually to handle multiple market parameters
	queryParts := []string{}
	for _, market := range marketNames {
		queryParts = append(queryParts, "market="+market)
	}
	if builderID != nil {
		queryParts = append(queryParts, fmt.Sprintf("builderId=%d", *builderID))
	}
	
	url := baseUrl
	if len(queryParts) > 0 {
		url += "?" + queryParts[0]
		for i := 1; i < len(queryParts); i++ {
			url += "&" + queryParts[i]
		}
	}

	var feeResponse models.FeeResponse
	if err := s.Base.DoRequest(ctx, "GET", url, nil, &feeResponse); err != nil {
		return nil, err
	}

	if feeResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", feeResponse.Status)
	}

	return feeResponse.Data, nil
}

// GetPositions retrieves current positions, optionally filtered by market names and position side
func (s *Service) GetPositions(ctx context.Context, marketNames []string, positionSide *models.PositionSide) ([]models.PositionModel, error) {
	// Build URL manually to handle multiple market parameters
	baseUrl := s.Base.EndpointConfig().APIBaseURL + "/user/positions"
	queryParts := []string{}
	for _, market := range marketNames {
		queryParts = append(queryParts, "market="+market)
	}
	if positionSide != nil {
		queryParts = append(queryParts, "side="+string(*positionSide))
	}

	url := baseUrl
	if len(queryParts) > 0 {
		url += "?" + queryParts[0]
		for i := 1; i < len(queryParts); i++ {
			url += "&" + queryParts[i]
		}
	}

	var positionsResponse models.PositionsResponse
	if err := s.Base.DoRequest(ctx, "GET", url, nil, &positionsResponse); err != nil {
		return nil, err
	}

	if positionsResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", positionsResponse.Status)
	}

	return positionsResponse.Data, nil
}

// GetPositionsHistory retrieves position history with optional filters
func (s *Service) GetPositionsHistory(ctx context.Context, marketNames []string, positionSide *models.PositionSide, cursor *int, limit *int) ([]models.PositionHistoryModel, error) {
	// Build URL manually to handle multiple market parameters
	baseUrl := s.Base.EndpointConfig().APIBaseURL + "/user/positions/history"
	queryParts := []string{}
	for _, market := range marketNames {
		queryParts = append(queryParts, "market="+market)
	}
	if positionSide != nil {
		queryParts = append(queryParts, "side="+string(*positionSide))
	}
	if cursor != nil {
		queryParts = append(queryParts, fmt.Sprintf("cursor=%d", *cursor))
	}
	if limit != nil {
		queryParts = append(queryParts, fmt.Sprintf("limit=%d", *limit))
	}

	url := baseUrl
	if len(queryParts) > 0 {
		url += "?" + queryParts[0]
		for i := 1; i < len(queryParts); i++ {
			url += "&" + queryParts[i]
		}
	}

	var positionsHistoryResponse models.PositionsHistoryResponse
	if err := s.Base.DoRequest(ctx, "GET", url, nil, &positionsHistoryResponse); err != nil {
		return nil, err
	}

	if positionsHistoryResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", positionsHistoryResponse.Status)
	}

	return positionsHistoryResponse.Data, nil
}

// GetOpenOrders retrieves open orders with optional filters
func (s *Service) GetOpenOrders(ctx context.Context, marketNames []string, orderType *models.OrderType, orderSide *models.OrderSide) ([]models.OpenOrderModel, error) {
	// Build URL manually to handle multiple market parameters
	baseUrl := s.Base.EndpointConfig().APIBaseURL + "/user/orders"
	queryParts := []string{}
	for _, market := range marketNames {
		queryParts = append(queryParts, "market="+market)
	}
	if orderType != nil {
		queryParts = append(queryParts, "type="+string(*orderType))
	}
	if orderSide != nil {
		queryParts = append(queryParts, "side="+string(*orderSide))
	}

	url := baseUrl
	if len(queryParts) > 0 {
		url += "?" + queryParts[0]
		for i := 1; i < len(queryParts); i++ {
			url += "&" + queryParts[i]
		}
	}

	var openOrdersResponse models.OpenOrdersResponse
	if err := s.Base.DoRequest(ctx, "GET", url, nil, &openOrdersResponse); err != nil {
		return nil, err
	}

	if openOrdersResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", openOrdersResponse.Status)
	}

	return openOrdersResponse.Data, nil
}

// GetOrdersHistory retrieves order history with optional filters
func (s *Service) GetOrdersHistory(ctx context.Context, marketNames []string, orderType *models.OrderType, orderSide *models.OrderSide, cursor *int, limit *int) ([]models.OpenOrderModel, error) {
	// Build URL manually to handle multiple market parameters
	baseUrl := s.Base.EndpointConfig().APIBaseURL + "/user/orders/history"
	queryParts := []string{}
	for _, market := range marketNames {
		queryParts = append(queryParts, "market="+market)
	}
	if orderType != nil {
		queryParts = append(queryParts, "type="+string(*orderType))
	}
	if orderSide != nil {
		queryParts = append(queryParts, "side="+string(*orderSide))
	}
	if cursor != nil {
		queryParts = append(queryParts, fmt.Sprintf("cursor=%d", *cursor))
	}
	if limit != nil {
		queryParts = append(queryParts, fmt.Sprintf("limit=%d", *limit))
	}

	url := baseUrl
	if len(queryParts) > 0 {
		url += "?" + queryParts[0]
		for i := 1; i < len(queryParts); i++ {
			url += "&" + queryParts[i]
		}
	}

	var ordersHistoryResponse models.OrdersHistoryResponse
	if err := s.Base.DoRequest(ctx, "GET", url, nil, &ordersHistoryResponse); err != nil {
		return nil, err
	}

	if ordersHistoryResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", ordersHistoryResponse.Status)
	}

	return ordersHistoryResponse.Data, nil
}

// GetOrderByID retrieves an order by its ID
func (s *Service) GetOrderByID(ctx context.Context, orderID int) (*models.OpenOrderModel, error) {
	baseUrl, err := s.Base.GetURL(fmt.Sprintf("/user/orders/%d", orderID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	var orderResponse models.OpenOrdersResponse
	if err := s.Base.DoRequest(ctx, "GET", baseUrl, nil, &orderResponse); err != nil {
		return nil, err
	}

	if orderResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", orderResponse.Status)
	}

	if len(orderResponse.Data) == 0 {
		return nil, fmt.Errorf("order not found")
	}

	return &orderResponse.Data[0], nil
}

// GetOrderByExternalID retrieves orders by external ID
func (s *Service) GetOrderByExternalID(ctx context.Context, externalID string) ([]models.OpenOrderModel, error) {
	baseUrl, err := s.Base.GetURL(fmt.Sprintf("/user/orders/external/%s", externalID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	var orderResponse models.OpenOrdersResponse
	if err := s.Base.DoRequest(ctx, "GET", baseUrl, nil, &orderResponse); err != nil {
		return nil, err
	}

	if orderResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", orderResponse.Status)
	}

	return orderResponse.Data, nil
}

// GetTrades retrieves trades with optional filters
func (s *Service) GetTrades(ctx context.Context, marketNames []string, tradeSide *models.OrderSide, tradeType *models.TradeType, cursor *int, limit *int) ([]models.AccountTradeModel, error) {
	// Build URL manually to handle multiple market parameters
	baseUrl := s.Base.EndpointConfig().APIBaseURL + "/user/trades"
	queryParts := []string{}
	for _, market := range marketNames {
		queryParts = append(queryParts, "market="+market)
	}
	if tradeSide != nil {
		queryParts = append(queryParts, "side="+string(*tradeSide))
	}
	if tradeType != nil {
		queryParts = append(queryParts, "type="+string(*tradeType))
	}
	if cursor != nil {
		queryParts = append(queryParts, fmt.Sprintf("cursor=%d", *cursor))
	}
	if limit != nil {
		queryParts = append(queryParts, fmt.Sprintf("limit=%d", *limit))
	}

	url := baseUrl
	if len(queryParts) > 0 {
		url += "?" + queryParts[0]
		for i := 1; i < len(queryParts); i++ {
			url += "&" + queryParts[i]
		}
	}

	var tradesResponse models.TradesResponse
	if err := s.Base.DoRequest(ctx, "GET", url, nil, &tradesResponse); err != nil {
		return nil, err
	}

	if tradesResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", tradesResponse.Status)
	}

	return tradesResponse.Data, nil
}

// GetLeverage retrieves leverage for specified markets
func (s *Service) GetLeverage(ctx context.Context, marketNames []string) ([]models.AccountLeverage, error) {
	// Build URL manually to handle multiple market parameters
	baseUrl := s.Base.EndpointConfig().APIBaseURL + "/user/leverage"
	queryParts := []string{}
	for _, market := range marketNames {
		queryParts = append(queryParts, "market="+market)
	}

	url := baseUrl
	if len(queryParts) > 0 {
		url += "?" + queryParts[0]
		for i := 1; i < len(queryParts); i++ {
			url += "&" + queryParts[i]
		}
	}

	var leverageResponse models.LeverageResponse
	if err := s.Base.DoRequest(ctx, "GET", url, nil, &leverageResponse); err != nil {
		return nil, err
	}

	if leverageResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", leverageResponse.Status)
	}

	return leverageResponse.Data, nil
}

// GetBridgeConfig retrieves the bridge configuration
func (s *Service) GetBridgeConfig(ctx context.Context) (*models.BridgesConfig, error) {
	baseUrl, err := s.Base.GetURL("/user/bridge/config", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	var bridgesConfigResponse models.BridgesConfigResponse
	if err := s.Base.DoRequest(ctx, "GET", baseUrl, nil, &bridgesConfigResponse); err != nil {
		return nil, err
	}

	if bridgesConfigResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", bridgesConfigResponse.Status)
	}

	return &bridgesConfigResponse.Data, nil
}

// GetBridgeQuote retrieves a bridge quote
func (s *Service) GetBridgeQuote(ctx context.Context, chainIn, chainOut string, amount decimal.Decimal) (*models.Quote, error) {
	query := map[string]string{
		"chainIn":  chainIn,
		"chainOut": chainOut,
		"amount":   amount.String(),
	}

	baseUrl, err := s.Base.GetURL("/user/bridge/quote", query)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	var quoteResponse models.QuoteResponse
	if err := s.Base.DoRequest(ctx, "GET", baseUrl, nil, &quoteResponse); err != nil {
		return nil, err
	}

	if quoteResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", quoteResponse.Status)
	}

	return &quoteResponse.Data, nil
}

// Methods to be implemented:
// - UpdateLeverage (new)
// - CommitBridgeQuote (new)
// - Transfer (new)
// - Withdraw (new)
// - AssetOperations (new)
