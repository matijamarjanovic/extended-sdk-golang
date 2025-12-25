package sdk

import (
	"context"
	"strconv"
	"testing"

	"github.com/extended-protocol/extended-sdk-golang/src/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Account Service Tests

func TestAccountService_GetAccount(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	account, err := client.Account.GetAccount(ctx)

	require.NoError(t, err, "Should not error when getting account")
	require.NotNil(t, account, "Account should not be nil")
	assert.Greater(t, account.ID, -1, "Account ID should be greater than -1")
	t.Logf("Account ID: %d, Description: %s, Status: %s", account.ID, account.Description, account.Status)
}

func TestAccountService_GetBalance(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	balance, err := client.Account.GetBalance(ctx)

	require.NoError(t, err, "Should not error when getting balance")
	require.NotNil(t, balance, "Balance should not be nil")
	assert.NotEmpty(t, balance.CollateralName, "Collateral name should not be empty")
	t.Logf("Balance: %s, Equity: %s, AvailableForTrade: %s", balance.Balance.String(), balance.Equity.String(), balance.AvailableForTrade.String())
}

func TestAccountService_GetMarketFee(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	fees, err := client.Account.GetMarketFee(ctx, "BTC-USD")

	require.NoError(t, err, "Should not error when getting market fee")
	require.Greater(t, len(fees), 0, "Should return at least one fee")
	t.Logf("Got %d fees for BTC-USD", len(fees))
}

func TestAccountService_GetFees(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	fees, err := client.Account.GetFees(ctx, []string{"BTC-USD", "ETH-USD"}, nil)

	require.NoError(t, err, "Should not error when getting fees")
	require.Greater(t, len(fees), 0, "Should return at least one fee")
	t.Logf("Got %d fees", len(fees))
}

func TestAccountService_GetFees_WithBuilderID(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	builderID := 1
	fees, err := client.Account.GetFees(ctx, []string{"BTC-USD"}, &builderID)

	require.NoError(t, err, "Should not error when getting fees with builder ID")
	require.Greater(t, len(fees), 0, "Should return at least one fee")
	t.Logf("Got %d fees with builder ID", len(fees))
}

func TestAccountService_GetPositions(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	positions, err := client.Account.GetPositions(ctx, []string{}, nil)

	require.NoError(t, err, "Should not error when getting positions")
	require.NotNil(t, positions, "Positions should not be nil")
	t.Logf("Got %d positions", len(positions))
}

func TestAccountService_GetPositions_WithMarketFilter(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	positions, err := client.Account.GetPositions(ctx, []string{"BTC-USD"}, nil)

	require.NoError(t, err, "Should not error when getting positions with market filter")
	require.NotNil(t, positions, "Positions should not be nil")
	t.Logf("Got %d positions for BTC-USD", len(positions))
}

func TestAccountService_GetPositions_WithSideFilter(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	side := models.PositionSideLong
	positions, err := client.Account.GetPositions(ctx, []string{}, &side)

	require.NoError(t, err, "Should not error when getting positions with side filter")
	require.NotNil(t, positions, "Positions should not be nil")
	t.Logf("Got %d LONG positions", len(positions))
}

func TestAccountService_GetPositionsHistory(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	history, err := client.Account.GetPositionsHistory(ctx, []string{}, nil, nil, nil)

	require.NoError(t, err, "Should not error when getting positions history")
	require.NotNil(t, history, "History should not be nil")
	t.Logf("Got %d position history entries", len(history))
}

func TestAccountService_GetPositionsHistory_WithFilters(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	side := models.PositionSideLong
	limit := 10
	history, err := client.Account.GetPositionsHistory(ctx, []string{"BTC-USD"}, &side, nil, &limit)

	require.NoError(t, err, "Should not error when getting positions history with filters")
	require.NotNil(t, history, "History should not be nil")
	assert.LessOrEqual(t, len(history), limit, "Should respect limit")
	t.Logf("Got %d position history entries", len(history))
}

func TestAccountService_GetOpenOrders(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	orders, err := client.Account.GetOpenOrders(ctx, []string{}, nil, nil)

	require.NoError(t, err, "Should not error when getting open orders")
	require.NotNil(t, orders, "Orders should not be nil")
	t.Logf("Got %d open orders", len(orders))
}

func TestAccountService_GetOpenOrders_WithFilters(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	orderType := models.OrderTypeLimit
	orderSide := models.OrderSideBuy
	orders, err := client.Account.GetOpenOrders(ctx, []string{"BTC-USD"}, &orderType, &orderSide)

	require.NoError(t, err, "Should not error when getting open orders with filters")
	require.NotNil(t, orders, "Orders should not be nil")
	t.Logf("Got %d open orders with filters", len(orders))
}

func TestAccountService_GetOrdersHistory(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	history, err := client.Account.GetOrdersHistory(ctx, []string{}, nil, nil, nil, nil)

	require.NoError(t, err, "Should not error when getting orders history")
	require.NotNil(t, history, "History should not be nil")
	t.Logf("Got %d order history entries", len(history))
}

func TestAccountService_GetOrdersHistory_WithFilters(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	orderType := models.OrderTypeLimit
	limit := 10
	history, err := client.Account.GetOrdersHistory(ctx, []string{"BTC-USD"}, &orderType, nil, nil, &limit)

	require.NoError(t, err, "Should not error when getting orders history with filters")
	require.NotNil(t, history, "History should not be nil")
	assert.LessOrEqual(t, len(history), limit, "Should respect limit")
	t.Logf("Got %d order history entries", len(history))
}

func TestAccountService_GetOrderByID(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	// First get open orders to find a valid order ID
	orders, err := client.Account.GetOpenOrders(ctx, []string{}, nil, nil)
	if err != nil || len(orders) == 0 {
		t.Skip("No open orders available for testing GetOrderByID")
		return
	}

	orderID := orders[0].ID
	order, err := client.Account.GetOrderByID(ctx, orderID)

	require.NoError(t, err, "Should not error when getting order by ID")
	require.NotNil(t, order, "Order should not be nil")
	assert.Equal(t, orderID, order.ID, "Order ID should match")
	t.Logf("Retrieved order ID: %d", order.ID)
}

func TestAccountService_GetOrderByExternalID(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	// First get open orders to find a valid external ID
	orders, err := client.Account.GetOpenOrders(ctx, []string{}, nil, nil)
	if err != nil || len(orders) == 0 {
		t.Skip("No open orders available for testing GetOrderByExternalID")
		return
	}

	externalID := orders[0].ID
	ordersByExtID, err := client.Account.GetOrderByExternalID(ctx, strconv.Itoa(externalID))

	require.NoError(t, err, "Should not error when getting order by external ID")
	require.NotNil(t, ordersByExtID, "Orders should not be nil")
	t.Logf("Retrieved %d orders by external ID", len(ordersByExtID))
}

func TestAccountService_GetTrades(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	trades, err := client.Account.GetTrades(ctx, []string{"BTC-USD"}, nil, nil, nil, nil)

	require.NoError(t, err, "Should not error when getting trades")
	require.NotNil(t, trades, "Trades should not be nil")
	t.Logf("Got %d trades", len(trades))
}

func TestAccountService_GetTrades_WithFilters(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	tradeSide := models.OrderSideBuy
	tradeType := models.TradeTypeTrade
	limit := 10
	trades, err := client.Account.GetTrades(ctx, []string{"BTC-USD", "ETH-USD"}, &tradeSide, &tradeType, nil, &limit)

	require.NoError(t, err, "Should not error when getting trades with filters")
	require.NotNil(t, trades, "Trades should not be nil")
	assert.LessOrEqual(t, len(trades), limit, "Should respect limit")
	t.Logf("Got %d trades with filters", len(trades))
}

func TestAccountService_GetLeverage(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	leverage, err := client.Account.GetLeverage(ctx, []string{"BTC-USD"})

	require.NoError(t, err, "Should not error when getting leverage")
	require.NotNil(t, leverage, "Leverage should not be nil")
	require.Greater(t, len(leverage), 0, "Should return at least one leverage entry")
	t.Logf("Got %d leverage entries", len(leverage))
	for _, l := range leverage {
		t.Logf("Market: %s, Leverage: %s", l.Market, l.Leverage.String())
	}
}

func TestAccountService_GetLeverage_MultipleMarkets(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	leverage, err := client.Account.GetLeverage(ctx, []string{"BTC-USD", "ETH-USD"})

	require.NoError(t, err, "Should not error when getting leverage for multiple markets")
	require.NotNil(t, leverage, "Leverage should not be nil")
	t.Logf("Got %d leverage entries", len(leverage))
}

func TestAccountService_GetClient(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	clientInfo, err := client.Account.GetClient(ctx)

	require.NoError(t, err, "Should not error when getting client")
	require.NotNil(t, clientInfo, "Client should not be nil")
	assert.Greater(t, clientInfo.ID, 0, "Client ID should be greater than 0")
	t.Logf("Client ID: %d", clientInfo.ID)
	if clientInfo.EvmWalletAddress != nil {
		t.Logf("EVM Wallet Address: %s", *clientInfo.EvmWalletAddress)
	}
	if clientInfo.StarknetWalletAddress != nil {
		t.Logf("Starknet Wallet Address: %s", *clientInfo.StarknetWalletAddress)
	}
}

func TestAccountService_GetBridgeConfig(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	config, err := client.Account.GetBridgeConfig(ctx)

	require.NoError(t, err, "Should not error when getting bridge config")
	require.NotNil(t, config, "Bridge config should not be nil")
	require.NotNil(t, config.Chains, "Chains should not be nil")
	t.Logf("Got bridge config with %d chains", len(config.Chains))
	for _, chain := range config.Chains {
		t.Logf("Chain: %s, Contract: %s", chain.Chain, chain.ContractAddress)
	}
}

func TestAccountService_GetBridgeQuote(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	// Get bridge config first to see available chains
	config, err := client.Account.GetBridgeConfig(ctx)
	if err != nil || len(config.Chains) == 0 {
		t.Skip("No bridge config available for testing GetBridgeQuote")
		return
	}

	// Use STRK as chainIn and the first available chain as chainOut
	chainIn := "STRK"
	chainOut := config.Chains[0].Chain
	amount := decimal.NewFromFloat(1.0) // 1.0 amount for testing

	quote, err := client.Account.GetBridgeQuote(ctx, chainIn, chainOut, amount)

	require.NoError(t, err, "Should not error when getting bridge quote")
	require.NotNil(t, quote, "Quote should not be nil")
	assert.NotEmpty(t, quote.ID, "Quote ID should not be empty")
	t.Logf("Quote ID: %s, Fee: %s", quote.ID, quote.Fee.String())
}

func TestAccountService_UpdateLeverage(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	// First get current leverage
	leverage, err := client.Account.GetLeverage(ctx, []string{"BTC-USD"})
	require.NoError(t, err, "Should not error when getting leverage")
	require.Greater(t, len(leverage), 0, "Should have at least one leverage entry")

	currentLeverage := leverage[0].Leverage
	newLeverage := decimal.NewFromFloat(2.0) // Set to 2x leverage

	// Update leverage
	err = client.Account.UpdateLeverage(ctx, "BTC-USD", newLeverage)
	require.NoError(t, err, "Should not error when updating leverage")

	// Verify the update
	updatedLeverage, err := client.Account.GetLeverage(ctx, []string{"BTC-USD"})
	require.NoError(t, err, "Should not error when getting updated leverage")
	require.Greater(t, len(updatedLeverage), 0, "Should have at least one leverage entry")
	
	// Restore original leverage
	err = client.Account.UpdateLeverage(ctx, "BTC-USD", currentLeverage)
	require.NoError(t, err, "Should not error when restoring leverage")
	
	t.Logf("Updated leverage from %s to %s and restored", currentLeverage.String(), newLeverage.String())
}

func TestAccountService_CommitBridgeQuote(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	// First get bridge config and a quote
	config, err := client.Account.GetBridgeConfig(ctx)
	if err != nil || len(config.Chains) == 0 {
		t.Skip("No bridge config available for testing CommitBridgeQuote")
		return
	}

	chainIn := "STRK"
	chainOut := config.Chains[0].Chain
	amount := decimal.NewFromFloat(1.0)

	quote, err := client.Account.GetBridgeQuote(ctx, chainIn, chainOut, amount)
	if err != nil {
		t.Skip("Failed to get bridge quote for testing CommitBridgeQuote")
		return
	}

	// Commit the quote
	err = client.Account.CommitBridgeQuote(ctx, quote.ID)
	require.NoError(t, err, "Should not error when committing bridge quote")
	t.Logf("Successfully committed bridge quote with ID: %s", quote.ID)
}

func TestAccountService_Withdraw(t *testing.T) {
	t.Skip("Withdraw is not yet implemented")
}

func TestAccountService_Transfer(t *testing.T) {
	t.Skip("Transfer is not yet implemented")
}

func TestAccountService_AssetOperations(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	// Test getting all asset operations
	operations, err := client.Account.AssetOperations(ctx, nil, nil, nil, nil, nil, nil, nil)

	require.NoError(t, err, "Should not error when getting asset operations")
	require.NotNil(t, operations, "Operations should not be nil")
	t.Logf("Got %d asset operations", len(operations))
}

func TestAccountService_AssetOperations_WithFilters(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	// Test with filters
	operationTypes := []models.AssetOperationType{models.AssetOperationTypeDeposit, models.AssetOperationTypeWithdrawal}
	operationStatuses := []models.AssetOperationStatus{models.AssetOperationStatusCompleted}
	limit := 10

	operations, err := client.Account.AssetOperations(ctx, nil, operationTypes, operationStatuses, nil, nil, nil, &limit)

	require.NoError(t, err, "Should not error when getting asset operations with filters")
	require.NotNil(t, operations, "Operations should not be nil")
	assert.LessOrEqual(t, len(operations), limit, "Should respect limit")
	t.Logf("Got %d asset operations with filters", len(operations))
	
	for _, op := range operations {
		t.Logf("Operation ID: %s, Type: %s, Status: %s, Amount: %s", op.ID, op.Type, op.Status, op.Amount.String())
	}
}

