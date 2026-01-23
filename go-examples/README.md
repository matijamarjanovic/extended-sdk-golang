# Go SDK Examples

Simple examples demonstrating SDK functionality.

## Examples

- **account** - Account operations (balance, positions, orders, trades, fees, leverage)
- **markets** - Market data (markets, statistics, orderbook, candles, funding rates)
- **orders** - Order management (place, cancel, mass cancel)
- **streaming** - WebSocket streaming (orderbooks, trades, funding, candles, account updates, prices)

## Running Examples

Copy `.env.example` to `.env` and fill in your credentials:

```bash
cp .env.example .env
```

Then run:

```bash
go run . account
go run . markets
go run . orders
go run . streaming
```
