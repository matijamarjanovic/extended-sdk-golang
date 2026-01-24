package main

import (
	"context"
	"fmt"

	"github.com/extended-protocol/extended-sdk-golang/x10"
	"github.com/shopspring/decimal"
)

func accountExample(client *x10.Client) {
	ctx := context.Background()
	markets := []string{"BTC-USD", "ETH-USD"}

	account, err1 := client.Account.GetAccount(ctx)
	balance, err2 := client.Account.GetBalance(ctx)
	positions, err3 := client.Account.GetPositions(ctx, markets, x10.PositionSideLong)
	openOrders, err4 := client.Account.GetOpenOrders(ctx, markets, x10.OrderTypeLimit, x10.OrderSideBuy)
	ordersHistory, err5 := client.Account.GetOrdersHistory(ctx, markets, x10.OrderTypeLimit, x10.OrderSideBuy, 0, 5)
	trades, err6 := client.Account.GetTrades(ctx, markets, x10.OrderSideBuy, x10.TradeTypeTrade, 0, 5)
	fees, err7 := client.Account.GetFees(ctx, markets, 0)
	leverage, err8 := client.Account.GetLeverage(ctx, markets)
	bridgeConfig, err9 := client.Account.GetBridgeConfig(ctx)

	for _, e := range []error{err1, err2, err3, err4, err5, err6, err7, err8, err9} {
		if e != nil {
			fmt.Printf("account error: %v\n", e)
			return
		}
	}

	fmt.Printf("%+v\n%+v\n%+v\n%+v\n%+v\n%+v\n%+v\n%+v\n%+v\n",
		account, balance, positions, openOrders, ordersHistory, trades, fees, leverage, bridgeConfig)

	newLeverage := decimal.NewFromFloat(10.0)
	if err := client.Account.UpdateLeverage(ctx, "BTC-USD", newLeverage); err != nil { // set leverage for market
		fmt.Printf("update leverage error: %v\n", err)
		return
	}
}
