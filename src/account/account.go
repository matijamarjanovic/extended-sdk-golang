package account

import (
	"github.com/extended-protocol/extended-sdk-golang/src"
	"github.com/extended-protocol/extended-sdk-golang/src/models"
)

// Service provides account-related API operations.
// It holds a reference to the main client to access shared infrastructure.
type Service struct {
	Client *sdk.Client // Reference to main client
}

// Methods to be implemented:
// - GetMarketFee (move from api_client.go)
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


