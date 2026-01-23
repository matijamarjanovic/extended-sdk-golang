package x10

import (
	"context"
	"testing"
	"time"

	"github.com/extended-protocol/extended-sdk-golang/x10/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarketsService_GetMarketsDict(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	marketsDict, err := client.Markets.GetMarketsDict(ctx)

	require.NoError(t, err, "should not error when getting markets dict")
	require.NotNil(t, marketsDict, "markets dict should not be nil")
	assert.Greater(t, len(marketsDict), 0, "should return at least one market")

	if btcMarket, exists := marketsDict["BTC-USD"]; exists {
		assert.Equal(t, "BTC-USD", btcMarket.Name, "market name should match key")
		t.Logf("found BTC-USD market: %+v", btcMarket)
	}

	t.Logf("got %d markets in dict", len(marketsDict))
}

func TestMarketsService_GetMarketStatistics(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	stats, err := client.Markets.GetMarketStatistics(ctx, "BTC-USD")

	require.NoError(t, err, "should not error when getting market statistics")
	require.NotNil(t, stats, "market statistics should not be nil")
	assert.Greater(t, stats.LastPrice.Cmp(decimal.Zero), 0, "last price should be greater than zero")
	t.Logf("BTC-USD stats - last price: %s, mark price: %s, funding rate: %s",
		stats.LastPrice.String(), stats.MarkPrice.String(), stats.FundingRate.String())
}

func TestMarketsService_GetMarketStatistics_InvalidMarket(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	stats, err := client.Markets.GetMarketStatistics(ctx, "INVALID-MARKET")

	require.NoError(t, err, "should not error for invalid market")
	require.NotNil(t, stats, "stats should not be nil (API returns stats object)")
	t.Logf("got stats for invalid market (API behavior)")
}

func TestMarketsService_GetCandlesHistory(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	limit := 50
	candles, err := client.Markets.GetCandlesHistory(
		ctx,
		"BTC-USD",
		models.CandleTypeTrades,
		models.CandleIntervalPT1H,
		&limit,
		nil,
	)

	require.NoError(t, err, "should not error when getting candles history")
	require.NotNil(t, candles, "candles should not be nil")
	assert.Greater(t, len(candles), 0, "should return at least one candle")

	if len(candles) > 0 {
		candle := candles[0]
		assert.Greater(t, candle.High.Cmp(candle.Low), 0, "high should be >= low")
		assert.Greater(t, candle.Timestamp, int64(0), "Timestamp should be positive")
		t.Logf("got %d candles, first candle: Open=%s, High=%s, Low=%s, Close=%s, Timestamp=%d",
			len(candles), candle.Open.String(), candle.High.String(), candle.Low.String(), candle.Close.String(), candle.Timestamp)
	}
}

func TestMarketsService_GetCandlesHistory_WithLimit(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	limit := 10
	candles, err := client.Markets.GetCandlesHistory(
		ctx,
		"BTC-USD",
		models.CandleTypeTrades,
		models.CandleIntervalPT5M,
		&limit,
		nil,
	)

	require.NoError(t, err, "should not error when getting candles history with limit")
	require.NotNil(t, candles, "candles should not be nil")
	assert.LessOrEqual(t, len(candles), limit, "should respect limit")
	t.Logf("got %d candles with limit of %d", len(candles), limit)
}

func TestMarketsService_GetCandlesHistory_WithEndTime(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	endTime := time.Now()
	limit := 5
	candles, err := client.Markets.GetCandlesHistory(
		ctx,
		"BTC-USD",
		models.CandleTypeMarkPrices,
		models.CandleIntervalPT15M,
		&limit,
		&endTime,
	)

	require.NoError(t, err, "should not error when getting candles history with end time")
	require.NotNil(t, candles, "candles should not be nil")

	if len(candles) > 0 {
		for _, candle := range candles {
			candleTime := time.UnixMilli(candle.Timestamp)
			assert.True(t, candleTime.Before(endTime) || candleTime.Equal(endTime),
				"candle timestamp should be before or equal to endTime")
		}
		t.Logf("got %d candles before endTime %v", len(candles), endTime)
	}
}

func TestMarketsService_GetCandlesHistory_DifferentTypes(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	candleTypes := []models.CandleType{
		models.CandleTypeTrades,
		models.CandleTypeMarkPrices,
		models.CandleTypeIndexPrices,
	}

	for _, candleType := range candleTypes {
		limit := 50
		candles, err := client.Markets.GetCandlesHistory(
			ctx,
			"BTC-USD",
			candleType,
			models.CandleIntervalPT1H,
			&limit,
			nil,
		)

		require.NoError(t, err, "should not error when getting candles for type %s", candleType)
		require.NotNil(t, candles, "candles should not be nil for type %s", candleType)
		t.Logf("got %d candles for type %s", len(candles), candleType)
	}
}

func TestMarketsService_GetFundingRatesHistory(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	endTime := time.Now()
	startTime := endTime.Add(-30 * 24 * time.Hour)

	fundingRates, err := client.Markets.GetFundingRatesHistory(ctx, "BTC-USD", startTime, endTime)

	require.NoError(t, err, "should not error when getting funding rates history")
	require.NotNil(t, fundingRates, "funding rates should not be nil")

	if len(fundingRates) > 0 {
		rate := fundingRates[0]
		assert.Equal(t, "BTC-USD", rate.Market, "market should match")
		assert.Greater(t, rate.Timestamp, int64(0), "timestamp should be positive")
		t.Logf("got %d funding rates, first rate: Market=%s, Rate=%s, Timestamp=%d",
			len(fundingRates), rate.Market, rate.FundingRate.String(), rate.Timestamp)
	} else {
		t.Logf("got 0 funding rates (this may be expected if no rates in the time range)")
	}
}

func TestMarketsService_GetFundingRatesHistory_LongRange(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	endTime := time.Now()
	startTime := endTime.Add(-30 * 24 * time.Hour)

	fundingRates, err := client.Markets.GetFundingRatesHistory(ctx, "BTC-USD", startTime, endTime)

	require.NoError(t, err, "should not error when getting funding rates history for long range")
	require.NotNil(t, fundingRates, "funding rates should not be nil")

	if len(fundingRates) > 0 {
		startMillis := startTime.UnixMilli()
		endMillis := endTime.UnixMilli()
		for _, rate := range fundingRates {
			assert.GreaterOrEqual(t, rate.Timestamp, startMillis, "rate timestamp should be >= startTime")
			assert.LessOrEqual(t, rate.Timestamp, endMillis, "rate timestamp should be <= endTime")
		}
		t.Logf("got %d funding rates in 7-day range", len(fundingRates))
	}
}

func TestMarketsService_GetOrderbookSnapshot(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	orderbook, err := client.Markets.GetOrderbookSnapshot(ctx, "BTC-USD")

	require.NoError(t, err, "should not error when getting orderbook snapshot")
	require.NotNil(t, orderbook, "orderbook should not be nil")

	if orderbook.Market != "" {
		assert.Equal(t, "BTC-USD", orderbook.Market, "market should match")
	}

	bidLen := 0
	askLen := 0
	if orderbook.Bid != nil {
		bidLen = len(orderbook.Bid)
	}
	if orderbook.Ask != nil {
		askLen = len(orderbook.Ask)
	}

	if bidLen > 0 {
		bid := orderbook.Bid[0]
		assert.Greater(t, bid.Price.Cmp(decimal.Zero), 0, "Bid price should be positive")
		assert.Greater(t, bid.Qty.Cmp(decimal.Zero), 0, "Bid quantity should be positive")
		t.Logf("top bid: price=%s, qty=%s", bid.Price.String(), bid.Qty.String())
	}

	if askLen > 0 {
		ask := orderbook.Ask[0]
		assert.Greater(t, ask.Price.Cmp(decimal.Zero), 0, "Ask price should be positive")
		assert.Greater(t, ask.Qty.Cmp(decimal.Zero), 0, "Ask quantity should be positive")
		t.Logf("top ask: price=%s, qty=%s", ask.Price.String(), ask.Qty.String())
	}

	if bidLen > 0 && askLen > 0 {
		topBid := orderbook.Bid[0].Price
		topAsk := orderbook.Ask[0].Price
		assert.LessOrEqual(t, topBid.Cmp(topAsk), 0, "top bid should be <= top ask")
		t.Logf("orderbook spread: %s (ask - bid)", topAsk.Sub(topBid).String())
	}

	t.Logf("orderbook has %d bids and %d asks", bidLen, askLen)
}

func TestMarketsService_GetOrderbookSnapshot_InvalidMarket(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	orderbook, err := client.Markets.GetOrderbookSnapshot(ctx, "INVALID-MARKET")

	require.NoError(t, err, "should not error for invalid market")
	require.NotNil(t, orderbook, "orderbook should not be nil")
	assert.Empty(t, orderbook.Market, "market should be empty for invalid market")
	t.Logf("got empty orderbook for invalid market (expected behavior)")
}
