# Orders Service

The Orders Service handles order placement, cancellation, and management.

## Key Differences from Python SDK

Order placement in Go uses a **functional options pattern** instead of optional keyword arguments. This is more complex than the filter approach but necessary because orders have many optional parameters that interact with each other.

**Python approach:**
```python
order = client.orders.place_order(
    market=market,
    amount=amount,
    price=price,
    side=OrderSide.BUY,
    order_type=OrderType.LIMIT,
    time_in_force=TimeInForce.GTT,
    self_trade_protection=SelfTradeProtection.DISABLED,
    post_only=True,  # optional
    reduce_only=False,  # optional
    builder_fee=fee,  # optional
)
```

**Go approach:**
```go
order, err := client.Orders.PlaceOrder(
    ctx,
    market, amount, price,
    x10.OrderSideBuy,
    x10.OrderTypeLimit,
    x10.TimeInForceGTT,
    x10.SelfTradeProtectionDisabled,
    x10.WithPostOnly(true),      // optional
    x10.WithReduceOnly(false),   // optional
    x10.WithBuilderFee(fee),     // optional
)
```

Required parameters are function arguments. Optional parameters use `With*` option functions.

## Methods

### PlaceOrder

Places an order on the exchange. Required parameters are passed as function arguments, optional parameters as option functions.

```go
order, err := client.Orders.PlaceOrder(
    ctx,
    market,                    // MarketModel
    decimal.NewFromFloat(0.1), // amount
    decimal.NewFromFloat(50000), // price
    x10.OrderSideBuy,          // side
    x10.OrderTypeLimit,        // order type
    x10.TimeInForceGTT,        // time in force
    x10.SelfTradeProtectionDisabled,
    // Optional parameters:
    x10.WithPostOnly(true),
    x10.WithReduceOnly(false),
    x10.WithExpireTime(time.Now().Add(1 * time.Hour)),
    x10.WithBuilderFee(builderFee),
    x10.WithBuilderID(2017),
    x10.WithNonce(12345),      // omit to auto-generate
    x10.WithOrderExternalID("custom-id"),
)
```

**Required Parameters:**
- `market`: Market model from `GetMarkets()`
- `syntheticAmount`: Order size
- `price`: Order price
- `side`: `x10.OrderSideBuy` or `x10.OrderSideSell`
- `orderType`: `x10.OrderTypeLimit`, `x10.OrderTypeMarket`, etc.
- `timeInForce`: `x10.TimeInForceGTT`, `x10.TimeInForceIOC`, `x10.TimeInForceFOK`
- `selfTradeProtectionLevel`: `x10.SelfTradeProtectionDisabled`, `x10.SelfTradeProtectionAccount`, `x10.SelfTradeProtectionClient`

**Optional Parameters (Option Functions):**
- `WithPostOnly(bool)`: Post-only order
- `WithReduceOnly(bool)`: Reduce-only order
- `WithExpireTime(time.Time)`: Order expiration time
- `WithBuilderFee(decimal.Decimal)`: Builder fee
- `WithBuilderID(int)`: Builder ID
- `WithNonce(int)`: Custom nonce (omit to auto-generate)
- `WithOrderExternalID(string)`: Custom external ID
- `WithPreviousOrderExternalID(string)`: For order replacement
- `WithTpSlType(x10.TpSlType)`: Take profit / stop loss type
- `WithTakeProfit(x10.TpSlTriggerParam)`: Take profit parameters
- `WithStopLoss(x10.TpSlTriggerParam)`: Stop loss parameters

**Response:**
Returns `*models.OrderResponse` with `OrderID` and `ExternalID`.

### CancelOrder

Cancels an order by internal order ID.

```go
err := client.Orders.CancelOrder(ctx, orderID)
```

### CancelOrderByExternalID

Cancels an order by external ID.

```go
err := client.Orders.CancelOrderByExternalID(ctx, externalID)
```

### MassCancel

Cancels multiple orders by IDs, external IDs, or markets.

```go
// Cancel by order IDs
err := client.Orders.MassCancel(ctx, []int{1, 2, 3}, []string{}, []string{}, false)

// Cancel by external IDs
err := client.Orders.MassCancel(ctx, []int{}, []string{"ext1", "ext2"}, []string{}, false)

// Cancel by markets
err := client.Orders.MassCancel(ctx, []int{}, []string{}, []string{"BTC-USD"}, false)

// Cancel all orders
err := client.Orders.MassCancel(ctx, []int{}, []string{}, []string{}, true)
```

**Parameters:**
- `orderIDs`: List of internal order IDs (use `[]int{}` to omit)
- `externalOrderIDs`: List of external IDs (use `[]string{}` to omit)
- `markets`: List of market names (use `[]string{}` to omit)
- `cancelAll`: If `true`, cancels all orders (ignores other parameters)
