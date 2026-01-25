package main

import (
	"context"
	"fmt"

	"github.com/matijamarjanovic/extended-sdk-golang/x10"
	"github.com/shopspring/decimal"
)

func accountExample(client *x10.Client) {
	ctx := context.Background()
	markets := []string{"BTC-USD", "ETH-USD"}

	account, err := client.Account.GetAccount(ctx)
	balance, _ := client.Account.GetBalance(ctx)
	positions, _ := client.Account.GetPositions(ctx, markets, x10.PositionSideLong)
	openOrders, _ := client.Account.GetOpenOrders(ctx, markets, x10.OrderTypeLimit, x10.OrderSideBuy)
	ordersHistory, _ := client.Account.GetOrdersHistory(ctx, markets, x10.OrderTypeLimit, x10.OrderSideBuy, 0, 5)
	trades, _ := client.Account.GetTrades(ctx, markets, x10.OrderSideBuy, x10.TradeTypeTrade, 0, 5)
	fees, _ := client.Account.GetFees(ctx, markets, 0)
	leverage, _ := client.Account.GetLeverage(ctx, markets)
	bridgeConfig, _ := client.Account.GetBridgeConfig(ctx)

	if err != nil {
		fmt.Printf("account error: %v\n", err)
		return
	}

	fmt.Printf("%+v\n%+v\n%+v\n%+v\n%+v\n%+v\n%+v\n%+v\n%+v\n",
		account, balance, positions, openOrders, ordersHistory, trades, fees, leverage, bridgeConfig)

	newLeverage := decimal.NewFromFloat(10.0)
	if err := client.Account.UpdateLeverage(ctx, "BTC-USD", newLeverage); err != nil { // set leverage for market
		fmt.Printf("update leverage error: %v\n", err)
		return
	}
}
