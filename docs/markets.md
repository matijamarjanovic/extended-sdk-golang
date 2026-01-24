# Markets Service

The Markets Service provides access to market data, statistics, orderbooks, candles, and funding rates.

## Key Differences from Python SDK

Similar to the Account Service, optional parameters use value-based sentinels:

- **Integer filters**: Use `0` to omit `limit` parameter
- **Time filters**: Use `time.Time{}` (zero value) to omit `endTime` parameter

## Methods

### GetMarkets

Retrieves market information for specified markets.

```go
// All markets
markets, err := client.Markets.GetMarkets(ctx, []string{})

// Specific markets
markets, err := client.Markets.GetMarkets(ctx, []string{"BTC-USD", "ETH-USD"})
```

### GetMarketsDict

Retrieves all markets as a map keyed by market name.

```go
marketsDict, err := client.Markets.GetMarketsDict(ctx)
// Access: btcMarket := marketsDict["BTC-USD"]
```

### GetMarketStatistics

Retrieves market statistics (volume, prices, funding rate, etc.).

```go
stats, err := client.Markets.GetMarketStatistics(ctx, "BTC-USD")
```

### GetOrderbookSnapshot

Retrieves current orderbook snapshot.

```go
orderbook, err := client.Markets.GetOrderbookSnapshot(ctx, "BTC-USD")
```

### GetCandlesHistory

Retrieves historical candle data.

```go
// All available candles (no limit, no end time)
candles, err := client.Markets.GetCandlesHistory(
    ctx,
    "BTC-USD",
    x10.CandleTypeTrades,
    x10.CandleIntervalPT1M,
    0,          // limit (0 = omit, use API default)
    time.Time{}, // endTime (zero value = omit)
)

// Limited candles with end time
endTime := time.Now()
candles, err := client.Markets.GetCandlesHistory(
    ctx,
    "BTC-USD",
    x10.CandleTypeTrades,
    x10.CandleIntervalPT1H,
    50,    // limit
    endTime,
)
```

**Parameters:**
- `marketName`: Market name (e.g., "BTC-USD")
- `candleType`: `x10.CandleTypeTrades`, `x10.CandleTypeMarkPrices`, or `x10.CandleTypeIndexPrices`
- `interval`: `x10.CandleIntervalPT1M`, `x10.CandleIntervalPT5M`, `x10.CandleIntervalPT1H`, etc.
- `limit`: Maximum number of candles (use `0` to omit, API will use default)
- `endTime`: End time for the range (use `time.Time{}` to omit)

**Candle Intervals:**
- `x10.CandleIntervalPT1M` - 1 minute
- `x10.CandleIntervalPT5M` - 5 minutes
- `x10.CandleIntervalPT15M` - 15 minutes
- `x10.CandleIntervalPT30M` - 30 minutes
- `x10.CandleIntervalPT1H` - 1 hour
- `x10.CandleIntervalPT2H` - 2 hours
- `x10.CandleIntervalPT4H` - 4 hours
- `x10.CandleIntervalP1D` - 1 day

### GetFundingRatesHistory

Retrieves historical funding rates.

```go
endTime := time.Now()
startTime := endTime.Add(-24 * time.Hour)
rates, err := client.Markets.GetFundingRatesHistory(ctx, "BTC-USD", startTime, endTime)
```

**Parameters:**
- `marketName`: Market name
- `startTime`: Start time (required)
- `endTime`: End time (required)
