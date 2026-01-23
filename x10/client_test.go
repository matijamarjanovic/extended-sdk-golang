package sdk

import (
	"context"
	"testing"
	"time"

	"github.com/extended-protocol/extended-sdk-golang/x10/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetMarkets_SingleValidMarket(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	markets, err := client.Markets.GetMarkets(ctx, []string{"BTC-USD"})

	require.NoError(t, err, "should not error when requesting BTC-USD market")
	require.Equal(t, len(markets), 1, "should return one market for valid request")
}

func TestClient_GetMarkets_MultipleValidMarkets(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()
	requestedMarkets := []string{"BTC-USD", "ETH-USD"}

	markets, err := client.Markets.GetMarkets(ctx, requestedMarkets)

	require.NoError(t, err, "should not error when requesting multiple valid markets")
	t.Logf("requested %v, got %d markets", requestedMarkets, len(markets))

	require.Equal(t, len(markets), len(requestedMarkets), "should return correct number of markets")
}

func TestClient_GetMarkets_InvalidMarket(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	markets, err := client.Markets.GetMarkets(ctx, []string{"INVALID-MARKET-NAME"})

	require.Error(t, err, "should error when requesting invalid market")
	assert.Equal(t, len(markets), 0, "should return zero markets for invalid request")
}

func TestClient_GetMarkets_ContextTimeout(t *testing.T) {
	client := createTestClient()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	_, err := client.Markets.GetMarkets(ctx, []string{"BTC-USD"})

	require.Error(t, err, "should error when context times out")
	t.Logf("got expected timeout error: %v", err)
}

func TestClient_GetMarkets_NetworkError(t *testing.T) {
	cfg := STARKNET_MAINNET_CONFIG
	cfg.APIBaseURL = "http://invalid-url-that-does-not-exist.com"
	account, _ := NewStarkPerpetualAccount(0, "0x0", "0x0", "")
	client := NewClient(cfg, account, 5*time.Second)
	ctx := context.Background()

	_, err := client.Markets.GetMarkets(ctx, []string{"BTC-USD"})

	require.Error(t, err, "should error when network request fails")
	t.Logf("got expected network error: %v", err)
}

func TestClient_GetMarketFee_ValidMarket(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	fees, err := client.Account.GetMarketFee(ctx, "BTC-USD")

	require.NoError(t, err, "should not error when requesting fees for BTC-USD market")
	require.Equal(t, len(fees), 1, "should return one fee entry for valid market")
	t.Logf("got %d fees for BTC-USD", len(fees))

	for _, fee := range fees {
		t.Logf("fee: %+v", fee)
	}
}

func TestClient_GetMarketFee_InvalidMarket(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	fees, err := client.Account.GetMarketFee(ctx, "INVALID-MARKET-NAME")

	assert.Error(t, err, "should error when requesting fees for invalid market")
	assert.Equal(t, len(fees), 0, "should return zero fees for invalid market")
}

func TestClient_GetMarketFee_ContextTimeout(t *testing.T) {
	client := createTestClient()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	_, err := client.Account.GetMarketFee(ctx, "BTC-USD")

	require.Error(t, err, "should error when context times out")
	t.Logf("got expected timeout error: %v", err)
}

func TestClient_GetMarketFee_NetworkError(t *testing.T) {
	cfg := STARKNET_MAINNET_CONFIG
	cfg.APIBaseURL = "http://invalid-url-that-does-not-exist.com"
	account, _ := NewStarkPerpetualAccount(0, "0x0", "0x0", "")
	client := NewClient(cfg, account, 5*time.Second)
	ctx := context.Background()

	_, err := client.Account.GetMarketFee(ctx, "BTC-USD")

	require.Error(t, err, "should error when network request fails")
	t.Logf("got expected network error: %v", err)
}

func TestClient_PlaceOrder_ValidOrder(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	markets, err := client.Markets.GetMarkets(ctx, []string{"BTC-USD"})
	require.NoError(t, err, "should be able to get BTC-USD market")
	require.Greater(t, len(markets), 0, "should have at least one market")

	market := markets[0]
	expireTime := time.Now().Add(1 * time.Hour)

	response, err := client.Orders.PlaceOrder(ctx,
		market,
		decimal.NewFromFloat(0.001),
		decimal.NewFromFloat(1),
		OrderSideBuy,
		models.OrderTypeLimit,
		TimeInForceGTT,
		SelfTradeProtectionDisabled,
		WithExpireTime(expireTime),
	)

	require.NoError(t, err, "should not error when placing valid order")
	require.NotNil(t, response, "response should not be nil")
	require.Equal(t, "OK", response.Status, "response status should be OK")

	t.Logf("successfully placed order with ID: %s", response.Data.ExternalID)
}
