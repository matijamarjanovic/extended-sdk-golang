# Account Service

The Account Service provides access to account information, positions, orders, trades, and asset operations.

## Key Differences from Python SDK

In Python, many filter parameters are optional (can be `None`). In Go, we use a value-based approach with sentinel values:

- **Enum filters**: Use `x10.PositionSideAll`, `x10.OrderTypeAll`, `x10.OrderSideAll`, `x10.TradeTypeAll` to omit filters
- **Integer filters**: Use `0` to omit `cursor`, `limit`, `builderID` parameters
- **String filters**: Use `""` to omit `id` parameter
- **Time filters**: Use `time.Time{}` (zero value) to omit time-based filters

This approach avoids `nil` pointers and makes the API more explicit and type-safe.

## Methods

### GetAccount

Retrieves account information.

```go
account, err := client.Account.GetAccount(ctx)
```

### GetBalance

Retrieves account balance.

```go
balance, err := client.Account.GetBalance(ctx)
```

### GetPositions

Retrieves current positions. Filter by markets and position side.

```go
// All positions for specified markets
positions, err := client.Account.GetPositions(ctx, []string{"BTC-USD", "ETH-USD"}, x10.PositionSideAll)

// Only long positions
positions, err := client.Account.GetPositions(ctx, []string{"BTC-USD"}, x10.PositionSideLong)
```

**Parameters:**
- `marketNames`: List of market names to filter by
- `positionSide`: `x10.PositionSideAll` (no filter), `x10.PositionSideLong`, or `x10.PositionSideShort`

### GetPositionsHistory

Retrieves position history with optional filters.

```go
// All history, no pagination
history, err := client.Account.GetPositionsHistory(ctx, []string{"BTC-USD"}, x10.PositionSideAll, 0, 0)

// Long positions only, with pagination
history, err := client.Account.GetPositionsHistory(ctx, []string{"BTC-USD"}, x10.PositionSideLong, 0, 10)
```

**Parameters:**
- `marketNames`: List of market names
- `positionSide`: Filter by side (use `x10.PositionSideAll` to omit)
- `cursor`: Pagination cursor (use `0` to omit)
- `limit`: Maximum results (use `0` to omit)

### GetOpenOrders

Retrieves open orders with optional filters.

```go
// All open orders
orders, err := client.Account.GetOpenOrders(ctx, []string{"BTC-USD"}, x10.OrderTypeAll, x10.OrderSideAll)

// Only limit buy orders
orders, err := client.Account.GetOpenOrders(ctx, []string{"BTC-USD"}, x10.OrderTypeLimit, x10.OrderSideBuy)
```

**Parameters:**
- `marketNames`: List of market names
- `orderType`: Filter by type (use `x10.OrderTypeAll` to omit)
- `orderSide`: Filter by side (use `x10.OrderSideAll` to omit)

### GetOrdersHistory

Retrieves order history with optional filters and pagination.

```go
// All history, no pagination
history, err := client.Account.GetOrdersHistory(ctx, []string{"BTC-USD"}, x10.OrderTypeAll, x10.OrderSideAll, 0, 0)

// Filtered with pagination
history, err := client.Account.GetOrdersHistory(ctx, []string{"BTC-USD"}, x10.OrderTypeLimit, x10.OrderSideBuy, 0, 10)
```

**Parameters:**
- `marketNames`: List of market names
- `orderType`: Filter by type (use `x10.OrderTypeAll` to omit)
- `orderSide`: Filter by side (use `x10.OrderSideAll` to omit)
- `cursor`: Pagination cursor (use `0` to omit)
- `limit`: Maximum results (use `0` to omit)

### GetTrades

Retrieves trade history with optional filters.

```go
// All trades, no pagination
trades, err := client.Account.GetTrades(ctx, []string{"BTC-USD"}, x10.OrderSideAll, x10.TradeTypeAll, 0, 0)

// Filtered trades with pagination
trades, err := client.Account.GetTrades(ctx, []string{"BTC-USD"}, x10.OrderSideBuy, x10.TradeTypeTrade, 0, 10)
```

**Parameters:**
- `marketNames`: List of market names
- `tradeSide`: Filter by side (use `x10.OrderSideAll` to omit)
- `tradeType`: Filter by type (use `x10.TradeTypeAll` to omit)
- `cursor`: Pagination cursor (use `0` to omit)
- `limit`: Maximum results (use `0` to omit)

### GetFees

Retrieves trading fees for specified markets.

```go
// All fees for markets
fees, err := client.Account.GetFees(ctx, []string{"BTC-USD", "ETH-USD"}, 0)

// Fees with builder ID
fees, err := client.Account.GetFees(ctx, []string{"BTC-USD"}, 2017)
```

**Parameters:**
- `marketNames`: List of market names
- `builderID`: Builder ID filter (use `0` to omit)

### AssetOperations

Retrieves asset operations history (deposits, withdrawals, transfers, claims).

```go
// All operations
operations, err := client.Account.AssetOperations(ctx, "", []x10.AssetOperationType{}, []x10.AssetOperationStatus{}, 0, 0, 0, 0)

// Filtered operations
operations, err := client.Account.AssetOperations(
    ctx,
    "",  // id filter (empty = omit)
    []x10.AssetOperationType{x10.AssetOperationTypeDeposit},
    []x10.AssetOperationStatus{x10.AssetOperationStatusCompleted},
    0,   // startTime (0 = omit)
    0,   // endTime (0 = omit)
    0,   // cursor (0 = omit)
    10,  // limit
)
```

**Parameters:**
- `id`: Operation ID filter (use `""` to omit)
- `operationTypes`: Filter by types (use `[]x10.AssetOperationType{}` to omit)
- `operationStatuses`: Filter by statuses (use `[]x10.AssetOperationStatus{}` to omit)
- `startTime`: Start timestamp in milliseconds (use `0` to omit)
- `endTime`: End timestamp in milliseconds (use `0` to omit)
- `cursor`: Pagination cursor (use `0` to omit)
- `limit`: Maximum results (use `0` to omit)
