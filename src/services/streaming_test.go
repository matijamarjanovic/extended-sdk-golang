package services

import (
	"context"
	"testing"
	"time"

	"github.com/extended-protocol/extended-sdk-golang/src/client"
	"github.com/extended-protocol/extended-sdk-golang/src/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestStreamingService() *StreamingService {
	cfg := models.EndpointConfig{
		StreamURL: "wss://api.starknet.sepolia.extended.exchange/stream.extended.exchange/v1",
	}
	baseClient := client.NewBaseClient(cfg, "test-api-key", nil, nil, 0)
	return &StreamingService{Base: baseClient}
}

func TestStreamingService_buildStreamURL(t *testing.T) {
	service := createTestStreamingService()

	tests := []struct {
		name     string
		path     string
		query    map[string]string
		expected string
	}{
		{
			name:     "simple path",
			path:     "/orderbooks",
			query:    nil,
			expected: "wss://api.starknet.sepolia.extended.exchange/stream.extended.exchange/v1/orderbooks",
		},
		{
			name:     "path with market",
			path:     "/orderbooks/BTC-USD",
			query:    nil,
			expected: "wss://api.starknet.sepolia.extended.exchange/stream.extended.exchange/v1/orderbooks/BTC-USD",
		},
		{
			name:  "path with query params",
			path:  "/orderbooks",
			query: map[string]string{"depth": "10"},
			expected: "wss://api.starknet.sepolia.extended.exchange/stream.extended.exchange/v1/orderbooks?depth=10",
		},
		{
			name:  "candles with interval",
			path:  "/candles/BTC-USD/trades",
			query: map[string]string{"interval": "PT1M"},
			expected: "wss://api.starknet.sepolia.extended.exchange/stream.extended.exchange/v1/candles/BTC-USD/trades?interval=PT1M",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.buildStreamURL(tt.path, tt.query)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Integration tests - these connect to real WebSocket endpoints
func createRealStreamingService() *StreamingService {
	cfg := models.EndpointConfig{
		StreamURL: "wss://api.starknet.sepolia.extended.exchange/stream.extended.exchange/v1",
	}
	baseClient := client.NewBaseClient(cfg, "", nil, nil, 10*time.Second)
	return &StreamingService{Base: baseClient}
}

func TestStreamingService_SubscribeToOrderbooks_RealConnection(t *testing.T) {

	service := createRealStreamingService()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Subscribe to orderbook for BTC-USD
	conn, err := service.SubscribeToOrderbooks(ctx, "BTC-USD", nil)
	require.NoError(t, err, "Should successfully connect to orderbook stream")
	defer conn.Close()

	// Try to receive at least one message
	var orderbook models.OrderbookUpdateModel
	err = conn.Recv(ctx, &orderbook)
	require.NoError(t, err, "Should receive at least one orderbook message")

	// Verify we got valid data
	assert.NotEmpty(t, orderbook.Market, "Market should be set")
	assert.Greater(t, conn.MessagesCount(), int64(0), "Should have received at least one message")
	t.Logf("Received orderbook update for market: %s", orderbook.Market)
}

func TestStreamingService_SubscribeToMarkPrices_RealConnection(t *testing.T) {

	service := createRealStreamingService()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Subscribe to mark prices for BTC-USD
	conn, err := service.SubscribeToMarkPrices(ctx, "BTC-USD")
	require.NoError(t, err, "Should successfully connect to mark price stream")
	defer conn.Close()

	// Try to receive at least one message
	var markPrice models.MarkPriceModel
	err = conn.Recv(ctx, &markPrice)
	require.NoError(t, err, "Should receive at least one mark price message")

	// Verify we got valid data
	assert.NotEmpty(t, markPrice.Market, "Market should be set")
	assert.True(t, markPrice.Price.GreaterThan(decimal.Zero), "Price should be positive")
	assert.Greater(t, conn.MessagesCount(), int64(0), "Should have received at least one message")
	t.Logf("Received mark price for %s: %s", markPrice.Market, markPrice.Price.String())
}

func TestStreamingService_SubscribeToIndexPrices_RealConnection(t *testing.T) {

	service := createRealStreamingService()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Subscribe to index prices for BTC-USD
	conn, err := service.SubscribeToIndexPrices(ctx, "BTC-USD")
	require.NoError(t, err, "Should successfully connect to index price stream")
	defer conn.Close()

	// Try to receive at least one message
	var indexPrice models.IndexPriceModel
	err = conn.Recv(ctx, &indexPrice)
	require.NoError(t, err, "Should receive at least one index price message")

	// Verify we got valid data
	assert.NotEmpty(t, indexPrice.Market, "Market should be set")
	assert.True(t, indexPrice.Price.GreaterThan(decimal.Zero), "Price should be positive")
	assert.Greater(t, conn.MessagesCount(), int64(0), "Should have received at least one message")
	t.Logf("Received index price for %s: %s", indexPrice.Market, indexPrice.Price.String())
}

func TestStreamingService_SubscribeToPublicTrades_RealConnection(t *testing.T) {

	service := createRealStreamingService()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Subscribe to public trades for BTC-USD
	conn, err := service.SubscribeToPublicTrades(ctx, "BTC-USD")
	require.NoError(t, err, "Should successfully connect to public trades stream")
	defer conn.Close()

	// Try to receive at least one message (may timeout if no trades occur)
	// Public trades stream returns an array of trades
	var trades []models.StreamPublicTradeModel
	err = conn.Recv(ctx, &trades)
	if err != nil {
		// If we timeout, that's okay - just verify connection was established
		if ctx.Err() == context.DeadlineExceeded {
			t.Log("No trades received within timeout, but connection was successful")
			return
		}
		require.NoError(t, err, "Should receive trades or timeout gracefully")
	}

	// If we got trades, verify them
	if err == nil {
		require.Greater(t, len(trades), 0, "Should receive at least one trade")
		trade := trades[0]
		assert.NotEmpty(t, trade.Market, "Market should be set")
		assert.True(t, trade.Price.GreaterThan(decimal.Zero), "Price should be positive")
		t.Logf("Received %d trades, first trade for %s: %s @ %s", len(trades), trade.Market, trade.Qty.String(), trade.Price.String())
	}
}

func TestStreamingService_SubscribeToFundingRates_RealConnection(t *testing.T) {

	service := createRealStreamingService()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Subscribe to funding rates for BTC-USD
	conn, err := service.SubscribeToFundingRates(ctx, "BTC-USD")
	require.NoError(t, err, "Should successfully connect to funding rates stream")
	defer conn.Close()

	// Try to receive at least one message (may timeout if no updates occur)
	var fundingRate models.FundingRateModel
	err = conn.Recv(ctx, &fundingRate)
	if err != nil {
		// If we timeout, that's okay - just verify connection was established
		if ctx.Err() == context.DeadlineExceeded {
			t.Log("No funding rate updates received within timeout, but connection was successful")
			return
		}
		require.NoError(t, err, "Should receive funding rate or timeout gracefully")
	}

	// If we got a funding rate, verify it
	if err == nil {
		assert.NotEmpty(t, fundingRate.Market, "Market should be set")
		t.Logf("Received funding rate for %s: %s", fundingRate.Market, fundingRate.FundingRate.String())
	}
}

func TestStreamingService_ConnectionClose(t *testing.T) {

	service := createRealStreamingService()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := service.SubscribeToMarkPrices(ctx, "BTC-USD")
	require.NoError(t, err)

	// Verify connection is open
	assert.False(t, conn.IsClosed(), "Connection should be open")

	// Close the connection
	err = conn.Close()
	require.NoError(t, err, "Should close connection successfully")

	// Verify connection is closed
	assert.True(t, conn.IsClosed(), "Connection should be closed")

	// Try to receive after close should fail
	var markPrice models.MarkPriceModel
	err = conn.Recv(ctx, &markPrice)
	assert.Error(t, err, "Should error when trying to receive on closed connection")
	assert.Contains(t, err.Error(), "closed", "Error should mention connection is closed")
}
