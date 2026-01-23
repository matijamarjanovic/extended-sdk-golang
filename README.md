# Golang Extended SDK

## Introduction

This is a Go SDK for the Extended Exchange perpetual trading platform. It provides a clean interface for account management, order placement, market data retrieval, and streaming functionality.

## Prerequisites

- Go 1.19 or later
- Rust toolchain (rustc, cargo) (development on the rust lib)
- GCC or compatible C compiler
- Git

Note: Currently, the SDK is only compatible for use on linux-based x86_64 machines (including WSL)

## Project Structure

```
extended-sdk-golang/
├── README.md           # This file
├── src/
│   ├── client.go       # Main SDK client with service access
│   ├── config.go       # Pre-configured endpoint configs
│   ├── client/         # Core client functionality
│   │   ├── base.go     # Base HTTP client
│   │   ├── sign.go     # Cryptographic signing
│   │   ├── starknet_account.go  # Account management
│   │   └── utils.go    # Utility functions
│   ├── models/         # Data models
│   │   ├── account.go
│   │   ├── common.go
│   │   ├── markets.go
│   │   ├── orders.go
│   │   ├── responses.go
│   │   └── streaming.go
│   └── services/       # API service implementations
│       ├── account.go   # Account operations
│       ├── markets.go  # Market data
│       ├── orders.go   # Order management
│       └── streaming.go # WebSocket streaming
└── rust-lib/          # Rust library source code
    └── target/
        └── release/   # Built Rust library (.so file)
```

## Building the Rust Library

1. Navigate to the rust-lib directory:
```bash
cd rust-lib
```

2. Build the release version of the Rust library:
```bash
cargo build --release
```

This will generate the shared library at `rust-lib/target/release/liborderffi.so` (Linux) or equivalent for your platform. You must then copy the library to the source directory.

Alternatively, run `build-lib.sh` in the root directory.

## Running Tests

After building the Rust library, set the library path and test credentials:

```bash
export LD_LIBRARY_PATH="$(pwd):${LD_LIBRARY_PATH:-}"
export TEST_API_KEY=<your_api_key>
export TEST_PRIVATE_KEY=<your_private_key_hex>
export TEST_PUBLIC_KEY=<your_public_key_hex>
export TEST_VAULT=<your_vault_id>
```

Then run tests from the `src/` directory:

```bash
cd src
go test -v
```

## Key Features

- **Account Management**: Get account info, balances, positions, orders, and trades
- **Order Placement**: Place, cancel, and manage orders with flexible options
- **Market Data**: Retrieve markets, orderbooks, candles, funding rates, and statistics
- **Streaming**: WebSocket support for real-time market data
- **Type Safety**: Strongly typed models and services

## API Overview

The SDK is organized into services accessible through the main `Client`:

- `client.Account` - Account operations (info, balance, positions, orders, trades)
- `client.Orders` - Order management (place, cancel, mass cancel)
- `client.Markets` - Market data (markets, orderbook, candles, funding rates)
- `client.Streaming` - WebSocket streaming functionality

## Troubleshooting

### Build Issues
- **CGO errors**: Ensure you have a C compiler installed and the Rust library is built
- **Library not found**: Check that `liborderffi.so` exists in the root directory and `LD_LIBRARY_PATH` is set correctly

### API Issues
- **Authentication errors**: Verify your API key is valid and properly set
- **Order validation errors**: Ensure order parameters are within acceptable ranges

### Testing Issues
- **Missing environment variables**: Tests require API keys and account credentials
- **Library loading errors**: Ensure the Rust library is built and `LD_LIBRARY_PATH` is set

