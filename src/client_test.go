package sdk

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/extended-protocol/extended-sdk-golang/src/orders"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() { load() }

func load() {
	wd, _ := os.Getwd()
	for {
		p := filepath.Join(wd, ".env")
		if _, err := os.Stat(p); err == nil {
			_ = godotenv.Load(p)
			return
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			return
		}
		wd = parent
	}
}
func createTestClient() *Client {
	apiKey := os.Getenv("TEST_API_KEY")
	vaultStr := os.Getenv("TEST_VAULT")
	vault, _ := strconv.ParseUint(vaultStr, 10, 64)
	publicKey := os.Getenv("TEST_PUBLIC_KEY")
	privateKey := os.Getenv("TEST_PRIVATE_KEY")

	account, err := NewStarkPerpetualAccount(vault, privateKey, publicKey, apiKey)

	if err != nil {
		panic("Failed to create StarkPerpetualAccount: " + err.Error())
	}

	return NewClient(STARKNET_MAINNET_CONFIG, account, 30*time.Second)
}

func TestClient_GetMarkets_SingleValidMarket(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	markets, err := client.Markets.GetMarkets(ctx, []string{"BTC-USD"})

	require.NoError(t, err, "Should not error when requesting BTC-USD market")
	require.Equal(t, len(markets), 1, "Should return one market for valid request")
}

func TestClient_GetMarkets_MultipleValidMarkets(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()
	requestedMarkets := []string{"BTC-USD", "ETH-USD"}

	markets, err := client.Markets.GetMarkets(ctx, requestedMarkets)

	require.NoError(t, err, "Should not error when requesting multiple valid markets")
	t.Logf("Requested %v, got %d markets", requestedMarkets, len(markets))

	require.Equal(t, len(markets), len(requestedMarkets), "Should return correct number of markets")
}

func TestClient_GetMarkets_InvalidMarket(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	markets, err := client.Markets.GetMarkets(ctx, []string{"INVALID-MARKET-NAME"})

	require.Error(t, err, "Should error when requesting invalid market")
	assert.Equal(t, len(markets), 0, "Should return zero markets for invalid request")
}

func TestClient_GetMarkets_ContextTimeout(t *testing.T) {
	client := createTestClient()

	// Create context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	_, err := client.Markets.GetMarkets(ctx, []string{"BTC-USD"})

	require.Error(t, err, "Should error when context times out")
	t.Logf("Got expected timeout error: %v", err)
}

func TestClient_GetMarkets_NetworkError(t *testing.T) {
	// Create client with invalid URL
	cfg := STARKNET_MAINNET_CONFIG
	cfg.APIBaseURL = "http://invalid-url-that-does-not-exist.com"
	account, _ := NewStarkPerpetualAccount(0, "0x0", "0x0", "")
	client := NewClient(cfg, account, 5*time.Second)
	ctx := context.Background()

	_, err := client.Markets.GetMarkets(ctx, []string{"BTC-USD"})

	require.Error(t, err, "Should error when network request fails")
	t.Logf("Got expected network error: %v", err)
}

func TestClient_GetMarketFee_ValidMarket(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	fees, err := client.Account.GetMarketFee(ctx, "BTC-USD")

	require.NoError(t, err, "Should not error when requesting fees for BTC-USD market")
	require.Equal(t, len(fees), 1, "Should return one fee entry for valid market")
	t.Logf("Got %d fees for BTC-USD", len(fees))

	for _, fee := range fees {
		t.Logf("Fee: %+v", fee)
	}
}

func TestClient_GetMarketFee_InvalidMarket(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	fees, err := client.Account.GetMarketFee(ctx, "INVALID-MARKET-NAME")

	// If no error, should return empty list or no matching fees
	assert.Error(t, err, "Should error when requesting fees for invalid market")
	assert.Equal(t, len(fees), 0, "Should return zero fees for invalid market")
}

func TestClient_GetMarketFee_ContextTimeout(t *testing.T) {
	client := createTestClient()

	// Create context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	_, err := client.Account.GetMarketFee(ctx, "BTC-USD")

	require.Error(t, err, "Should error when context times out")
	t.Logf("Got expected timeout error: %v", err)
}

func TestClient_GetMarketFee_NetworkError(t *testing.T) {
	// Create client with invalid URL
	cfg := STARKNET_MAINNET_CONFIG
	cfg.APIBaseURL = "http://invalid-url-that-does-not-exist.com"
	account, _ := NewStarkPerpetualAccount(0, "0x0", "0x0", "")
	client := NewClient(cfg, account, 5*time.Second)
	ctx := context.Background()

	_, err := client.Account.GetMarketFee(ctx, "BTC-USD")

	require.Error(t, err, "Should error when network request fails")
	t.Logf("Got expected network error: %v", err)
}

func TestClient_SubmitOrder_ValidOrder(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	// First get a market to use for the order
	markets, err := client.Markets.GetMarkets(ctx, []string{"BTC-USD"})
	require.NoError(t, err, "Should be able to get BTC-USD market")
	require.Greater(t, len(markets), 0, "Should have at least one market")

	market := markets[0]

	// Get account from client
	account, err := client.StarkAccount()
	require.NoError(t, err, "Should be able to get account")

	// Create order parameters
	nonce := int(time.Now().Unix()) // Use timestamp as nonce for uniqueness
	expireTime := time.Now().Add(1 * time.Hour)

	// Get config from client to use its StarknetDomain
	cfg := client.EndpointConfig()

	params := orders.CreateOrderObjectParams{
		Market:          market,
		Account:         account,
		SyntheticAmount: decimal.NewFromFloat(0.001), // Small BTC amount
		Price:           decimal.NewFromFloat(1),     // Place a low price so that it doesn't match
		Side:            OrderSideBuy,
		Signer:          account.Sign,
		StarknetDomain:  cfg.StarknetDomain,
		ExpireTime:               &expireTime,
		PostOnly:                 false,
		TimeInForce:              TimeInForceGTT,
		SelfTradeProtectionLevel: SelfTradeProtectionDisabled,
		Nonce:                    &nonce,
	}

	// Create the order object
	order, err := orders.CreateOrderObject(params)
	require.NoError(t, err, "Should be able to create valid order")
	require.NotNil(t, order, "Order should not be nil")

	// Submit the order
	response, err := client.Orders.SubmitOrder(ctx, order)

	require.NoError(t, err, "Should not error when submitting valid order")

	require.NotNil(t, response, "Response should not be nil")
	require.Equal(t, "OK", response.Status, "Response status should be OK")

	t.Logf("Successfully submitted order with ID: %s", response.Data.ExternalID)
}
