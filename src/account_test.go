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

func TestAccountService_GetAccount(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	account, err := client.Account.GetAccount(ctx)

	require.NoError(t, err, "should not error when getting account")
	require.NotNil(t, account, "account should not be nil")
	assert.Greater(t, account.ID, -1, "account ID should be greater than -1")
	t.Logf("account ID: %d, description: %s, status: %s", account.ID, account.Description, account.Status)
}

func TestAccountService_GetBalance(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	balance, err := client.Account.GetBalance(ctx)

	require.NoError(t, err, "should not error when getting balance")
	require.NotNil(t, balance, "balance should not be nil")
	assert.NotEmpty(t, balance.CollateralName, "collateral name should not be empty")
	t.Logf("balance: %s, equity: %s, availableForTrade: %s", balance.Balance.String(), balance.Equity.String(), balance.AvailableForTrade.String())
}

func TestAccountService_GetMarketFee(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	fees, err := client.Account.GetMarketFee(ctx, "BTC-USD")

	require.NoError(t, err, "should not error when getting market fee")
	require.Greater(t, len(fees), 0, "should return at least one fee")
	t.Logf("got %d fees for BTC-USD", len(fees))
}

func TestAccountService_GetFees(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	fees, err := client.Account.GetFees(ctx, []string{"BTC-USD", "ETH-USD"}, nil)

	require.NoError(t, err, "should not error when getting fees")
	require.Greater(t, len(fees), 0, "should return at least one fee")
	t.Logf("got %d fees", len(fees))
}

func TestAccountService_GetFees_WithBuilderID(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	builderID := 1
	fees, err := client.Account.GetFees(ctx, []string{"BTC-USD"}, &builderID)

	require.NoError(t, err, "should not error when getting fees with builder ID")
	require.Greater(t, len(fees), 0, "should return at least one fee")
	t.Logf("got %d fees with builder ID", len(fees))
}

func TestAccountService_GetPositions(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	positions, err := client.Account.GetPositions(ctx, []string{}, nil)

	require.NoError(t, err, "should not error when getting positions")
	require.NotNil(t, positions, "positions should not be nil")
	t.Logf("got %d positions", len(positions))
}

func TestAccountService_GetPositions_WithMarketFilter(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	positions, err := client.Account.GetPositions(ctx, []string{"BTC-USD"}, nil)

	require.NoError(t, err, "should not error when getting positions with market filter")
	require.NotNil(t, positions, "positions should not be nil")
	t.Logf("got %d positions for BTC-USD", len(positions))
}

func TestAccountService_GetPositions_WithSideFilter(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	side := models.PositionSideLong
	positions, err := client.Account.GetPositions(ctx, []string{}, &side)

	require.NoError(t, err, "should not error when getting positions with side filter")
	require.NotNil(t, positions, "positions should not be nil")
	t.Logf("got %d LONG positions", len(positions))
}

func TestAccountService_GetPositionsHistory(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	history, err := client.Account.GetPositionsHistory(ctx, []string{}, nil, nil, nil)

	require.NoError(t, err, "should not error when getting positions history")
	require.NotNil(t, history, "history should not be nil")
	t.Logf("got %d position history entries", len(history))
}

func TestAccountService_GetPositionsHistory_WithFilters(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	side := models.PositionSideLong
	limit := 10
	history, err := client.Account.GetPositionsHistory(ctx, []string{"BTC-USD"}, &side, nil, &limit)

	require.NoError(t, err, "should not error when getting positions history with filters")
	require.NotNil(t, history, "history should not be nil")
	assert.LessOrEqual(t, len(history), limit, "should respect limit")
	t.Logf("got %d position history entries", len(history))
}

func TestAccountService_GetOpenOrders(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	orders, err := client.Account.GetOpenOrders(ctx, []string{}, nil, nil)

	require.NoError(t, err, "should not error when getting open orders")
	require.NotNil(t, orders, "orders should not be nil")
	t.Logf("got %d open orders", len(orders))
}

func TestAccountService_GetOpenOrders_WithFilters(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	orderType := models.OrderTypeLimit
	orderSide := models.OrderSideBuy
	orders, err := client.Account.GetOpenOrders(ctx, []string{"BTC-USD"}, &orderType, &orderSide)

	require.NoError(t, err, "should not error when getting open orders with filters")
	require.NotNil(t, orders, "orders should not be nil")
	t.Logf("got %d open orders with filters", len(orders))
}

func TestAccountService_GetOrdersHistory(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	history, err := client.Account.GetOrdersHistory(ctx, []string{}, nil, nil, nil, nil)

	require.NoError(t, err, "should not error when getting orders history")
	require.NotNil(t, history, "history should not be nil")
	t.Logf("got %d order history entries", len(history))
}

func TestAccountService_GetOrdersHistory_WithFilters(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	orderType := models.OrderTypeLimit
	limit := 10
	history, err := client.Account.GetOrdersHistory(ctx, []string{"BTC-USD"}, &orderType, nil, nil, &limit)

	require.NoError(t, err, "should not error when getting orders history with filters")
	require.NotNil(t, history, "history should not be nil")
	assert.LessOrEqual(t, len(history), limit, "should respect limit")
	t.Logf("got %d order history entries", len(history))
}

func TestAccountService_GetOrderByID(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	orders, err := client.Account.GetOpenOrders(ctx, []string{}, nil, nil)
	if err != nil || len(orders) == 0 {
		t.Skip("No open orders available for testing GetOrderByID")
		return
	}

	orderID := orders[0].ID
	order, err := client.Account.GetOrderByID(ctx, orderID)

	require.NoError(t, err, "should not error when getting order by ID")
	require.NotNil(t, order, "order should not be nil")
	assert.Equal(t, orderID, order.ID, "order ID should match")
	t.Logf("retrieved order ID: %d", order.ID)
}

func TestAccountService_GetOrderByExternalID(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	orders, err := client.Account.GetOpenOrders(ctx, []string{}, nil, nil)
	if err != nil || len(orders) == 0 {
		t.Skip("No open orders available for testing GetOrderByExternalID")
		return
	}

	externalID := orders[0].ID
	ordersByExtID, err := client.Account.GetOrderByExternalID(ctx, strconv.Itoa(externalID))

	require.NoError(t, err, "should not error when getting order by external ID")
	require.NotNil(t, ordersByExtID, "orders should not be nil")
	t.Logf("retrieved %d orders by external ID", len(ordersByExtID))
}

func TestAccountService_GetTrades(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	trades, err := client.Account.GetTrades(ctx, []string{"BTC-USD"}, nil, nil, nil, nil)

	require.NoError(t, err, "should not error when getting trades")
	require.NotNil(t, trades, "trades should not be nil")
	t.Logf("got %d trades", len(trades))
}

func TestAccountService_GetTrades_WithFilters(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	tradeSide := models.OrderSideBuy
	tradeType := models.TradeTypeTrade
	limit := 10
	trades, err := client.Account.GetTrades(ctx, []string{"BTC-USD", "ETH-USD"}, &tradeSide, &tradeType, nil, &limit)

	require.NoError(t, err, "should not error when getting trades with filters")
	require.NotNil(t, trades, "trades should not be nil")
	assert.LessOrEqual(t, len(trades), limit, "should respect limit")
	t.Logf("got %d trades with filters", len(trades))
}

func TestAccountService_GetLeverage(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	leverage, err := client.Account.GetLeverage(ctx, []string{"BTC-USD"})

	require.NoError(t, err, "should not error when getting leverage")
	require.NotNil(t, leverage, "leverage should not be nil")
	require.Greater(t, len(leverage), 0, "should return at least one leverage entry")
	t.Logf("got %d leverage entries", len(leverage))
	for _, l := range leverage {
		t.Logf("market: %s, leverage: %s", l.Market, l.Leverage.String())
	}
}

func TestAccountService_GetLeverage_MultipleMarkets(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	leverage, err := client.Account.GetLeverage(ctx, []string{"BTC-USD", "ETH-USD"})

	require.NoError(t, err, "should not error when getting leverage for multiple markets")
	require.NotNil(t, leverage, "leverage should not be nil")
	t.Logf("got %d leverage entries", len(leverage))
}

func TestAccountService_GetClient(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	clientInfo, err := client.Account.GetClient(ctx)

	require.NoError(t, err, "should not error when getting client")
	require.NotNil(t, clientInfo, "client should not be nil")
	assert.Greater(t, clientInfo.ID, 0, "client ID should be greater than 0")
	t.Logf("client ID: %d", clientInfo.ID)
	if clientInfo.EvmWalletAddress != nil {
		t.Logf("EVM wallet address: %s", *clientInfo.EvmWalletAddress)
	}
	if clientInfo.StarknetWalletAddress != nil {
		t.Logf("Starknet wallet address: %s", *clientInfo.StarknetWalletAddress)
	}
}

func TestAccountService_GetBridgeConfig(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	config, err := client.Account.GetBridgeConfig(ctx)

	require.NoError(t, err, "should not error when getting bridge config")
	require.NotNil(t, config, "bridge config should not be nil")
	require.NotNil(t, config.Chains, "chains should not be nil")
	t.Logf("got bridge config with %d chains", len(config.Chains))
	for _, chain := range config.Chains {
		t.Logf("chain: %s, contract: %s", chain.Chain, chain.ContractAddress)
	}
}

func TestAccountService_GetBridgeQuote(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	config, err := client.Account.GetBridgeConfig(ctx)
	if err != nil || len(config.Chains) == 0 {
		t.Skip("no bridge config available for testing GetBridgeQuote")
		return
	}

	chainIn := config.Chains[0].Chain
	chainOut := "STRK"
	amount := decimal.NewFromFloat(1.0)

	quote, err := client.Account.GetBridgeQuote(ctx, chainIn, chainOut, amount)

	require.NoError(t, err, "should not error when getting bridge quote")
	require.NotNil(t, quote, "quote should not be nil")
	assert.NotEmpty(t, quote.ID, "quote ID should not be empty")
	t.Logf("quote ID: %s, fee: %s", quote.ID, quote.Fee.String())
}

func TestAccountService_UpdateLeverage(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	leverage, err := client.Account.GetLeverage(ctx, []string{"BTC-USD"})
	require.NoError(t, err, "should not error when getting leverage")
	require.Greater(t, len(leverage), 0, "should have at least one leverage entry")

	currentLeverage := leverage[0].Leverage
	newLeverage := decimal.NewFromFloat(2.0)

	err = client.Account.UpdateLeverage(ctx, "BTC-USD", newLeverage)
	require.NoError(t, err, "should not error when updating leverage")

	updatedLeverage, err := client.Account.GetLeverage(ctx, []string{"BTC-USD"})
	require.NoError(t, err, "should not error when getting updated leverage")
	require.Greater(t, len(updatedLeverage), 0, "should have at least one leverage entry")

	err = client.Account.UpdateLeverage(ctx, "BTC-USD", currentLeverage)
	require.NoError(t, err, "should not error when restoring leverage")

	t.Logf("updated leverage from %s to %s and restored", currentLeverage.String(), newLeverage.String())
}

func TestAccountService_CommitBridgeQuote(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	config, err := client.Account.GetBridgeConfig(ctx)
	if err != nil || len(config.Chains) == 0 {
		t.Skip("no bridge config available for testing CommitBridgeQuote")
		return
	}

	chainIn := "STRK"
	chainOut := config.Chains[0].Chain
	amount := decimal.NewFromFloat(1.0)

	quote, err := client.Account.GetBridgeQuote(ctx, chainIn, chainOut, amount)
	if err != nil {
		t.Skip("failed to get bridge quote for testing CommitBridgeQuote")
		return
	}

	err = client.Account.CommitBridgeQuote(ctx, quote.ID)
	require.NoError(t, err, "should not error when committing bridge quote")
	t.Logf("successfully committed bridge quote with ID: %s", quote.ID)
}

func TestAccountService_Withdraw(t *testing.T) {
	t.Skip("withdraw is not yet implemented")
}

func TestAccountService_Transfer(t *testing.T) {
	t.Skip("transfer is not yet implemented")
}

func TestAccountService_AssetOperations(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	operations, err := client.Account.AssetOperations(ctx, nil, nil, nil, nil, nil, nil, nil)

	require.NoError(t, err, "should not error when getting asset operations")
	require.NotNil(t, operations, "operations should not be nil")
	t.Logf("got %d asset operations", len(operations))
}

func TestAccountService_AssetOperations_WithFilters(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	operationTypes := []models.AssetOperationType{models.AssetOperationTypeDeposit, models.AssetOperationTypeWithdrawal}
	operationStatuses := []models.AssetOperationStatus{models.AssetOperationStatusCompleted}
	limit := 10

	operations, err := client.Account.AssetOperations(ctx, nil, operationTypes, operationStatuses, nil, nil, nil, &limit)

	require.NoError(t, err, "should not error when getting asset operations with filters")
	require.NotNil(t, operations, "operations should not be nil")
	assert.LessOrEqual(t, len(operations), limit, "should respect limit")
	t.Logf("got %d asset operations with filters", len(operations))

	for _, op := range operations {
		t.Logf("operation ID: %s, type: %s, status: %s, amount: %s", op.ID, op.Type, op.Status, op.Amount.String())
	}
}
