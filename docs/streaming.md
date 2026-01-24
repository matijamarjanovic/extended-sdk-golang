# Streaming Service

The Streaming Service provides real-time WebSocket connections for market data and account updates.

## Key Differences from Python SDK

The Go streaming API follows a similar pattern to Python, but uses Go's type system for message handling. The main difference is in how you receive and unmarshal messages.

**Python approach:**
```python
stream = client.streaming.subscribe_to_orderbooks("BTC-USD", depth=10)
message = stream.recv()  # Returns dict
```

**Go approach:**
```go
stream, _ := client.Streaming.SubscribeToOrderbooks(ctx, "BTC-USD", 10)
var orderbook x10.OrderbookUpdateModel
stream.Recv(ctx, &orderbook)  // Unmarshals into struct
```

Go requires you to declare the expected message type and pass a pointer for unmarshaling. This provides type safety and better IDE support.

## Connection Management

All stream connections should be closed with `defer` to prevent resource leaks:

```go
stream, err := client.Streaming.SubscribeToOrderbooks(ctx, "BTC-USD", 10)
if err != nil {
    return err
}
defer stream.Close()  // Always close connections
```

## Methods

### SubscribeToOrderbooks

Subscribes to orderbook updates for a market.

```go
// Specific market with depth
stream, err := client.Streaming.SubscribeToOrderbooks(ctx, "BTC-USD", 10)
defer stream.Close()

var orderbook x10.OrderbookUpdateModel
stream.Recv(ctx, &orderbook)
```

**Parameters:**
- `marketName`: Market name (use `""` to subscribe to all markets)
- `depth`: Orderbook depth (use `0` to use API default)

**Message Type:** `x10.OrderbookUpdateModel`

### SubscribeToPublicTrades

Subscribes to public trade updates.

```go
stream, err := client.Streaming.SubscribeToPublicTrades(ctx, "BTC-USD")
defer stream.Close()

var trade x10.StreamPublicTradeModel
stream.Recv(ctx, &trade)
```

**Parameters:**
- `marketName`: Market name (use `""` to subscribe to all markets)

**Message Type:** `x10.StreamPublicTradeModel`

### SubscribeToFundingRates

Subscribes to funding rate updates.

```go
stream, err := client.Streaming.SubscribeToFundingRates(ctx, "BTC-USD")
defer stream.Close()

var rate x10.FundingRateModel
stream.Recv(ctx, &rate)
```

**Parameters:**
- `marketName`: Market name (use `""` to subscribe to all markets)

**Message Type:** `x10.FundingRateModel`

### SubscribeToCandles

Subscribes to real-time candle updates.

```go
stream, err := client.Streaming.SubscribeToCandles(
    ctx,
    "BTC-USD",
    x10.CandleTypeTrades,
    x10.CandleIntervalPT1M,
)
defer stream.Close()

var candle x10.CandleModel
stream.Recv(ctx, &candle)
```

**Parameters:**
- `marketName`: Market name
- `candleType`: `x10.CandleTypeTrades`, `x10.CandleTypeMarkPrices`, or `x10.CandleTypeIndexPrices`
- `interval`: Candle interval (e.g., `x10.CandleIntervalPT1M`)

**Message Type:** `x10.CandleModel`

### SubscribeToAccountUpdates

Subscribes to account updates (orders, positions, trades, balance).

```go
stream, err := client.Streaming.SubscribeToAccountUpdates(ctx)
defer stream.Close()

var update x10.AccountStreamDataModel
stream.Recv(ctx, &update)
```

**Message Type:** `x10.AccountStreamDataModel`

### SubscribeToMarkPrices

Subscribes to mark price updates.

```go
stream, err := client.Streaming.SubscribeToMarkPrices(ctx, "BTC-USD")
defer stream.Close()

var price x10.MarkPriceModel
stream.Recv(ctx, &price)
```

**Parameters:**
- `marketName`: Market name (use `""` to subscribe to all markets)

**Message Type:** `x10.MarkPriceModel`

### SubscribeToIndexPrices

Subscribes to index price updates.

```go
stream, err := client.Streaming.SubscribeToIndexPrices(ctx, "BTC-USD")
defer stream.Close()

var price x10.IndexPriceModel
stream.Recv(ctx, &price)
```

**Parameters:**
- `marketName`: Market name (use `""` to subscribe to all markets)

**Message Type:** `x10.IndexPriceModel`

## Receiving Messages

The `Recv` method blocks until a message is received or the context is cancelled. Always pass a pointer to the expected message type:

```go
var orderbook x10.OrderbookUpdateModel
err := stream.Recv(ctx, &orderbook)
if err != nil {
    // Handle error (connection closed, context cancelled, etc.)
    return err
}
// Use orderbook data
fmt.Printf("Market: %s\n", orderbook.Market)
```

## Connection Status

Check connection status and message count:

```go
if stream.IsClosed() {
    // Connection is closed
}

count := stream.MessagesCount()  // Number of messages received
```
