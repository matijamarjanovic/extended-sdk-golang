package account

import (
	"context"
	"fmt"

	"github.com/extended-protocol/extended-sdk-golang/src/client"
	"github.com/extended-protocol/extended-sdk-golang/src/models"
)

// Service provides account-related API operations.
type Service struct {
	Base *client.BaseModule
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

// Methods to be implemented:
// - GetBalance (new)
// - GetAccount (new)
// - GetClient (new)
// - GetPositions (new)
// - GetPositionsHistory (new)
// - GetOpenOrders (new)
// - GetOrdersHistory (new)
// - GetOrderByID (new)
// - GetOrderByExternalID (new)
// - GetTrades (new)
// - GetFees (new)
// - GetLeverage (new)
// - UpdateLeverage (new)
// - GetBridgeConfig (new)
// - GetBridgeQuote (new)
// - CommitBridgeQuote (new)
// - Transfer (new)
// - Withdraw (new)
// - AssetOperations (new)


