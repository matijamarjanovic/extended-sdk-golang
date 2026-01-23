package examples

import (
	"context"
	"fmt"

	sdk "github.com/matijamarjanovic/extended-sdk-golang/x10"
	"github.com/shopspring/decimal"
)

func accountExample(client *sdk.Client) {
	ctx := context.Background()

	// retrieve account information
	account, _ := client.Account.GetAccount(ctx)
	balance, _ := client.Account.GetBalance(ctx)
	positions, _ := client.Account.GetPositions(ctx, []string{}, nil)
	openOrders, _ := client.Account.GetOpenOrders(ctx, []string{}, nil, nil)
	ordersHistory, _ := client.Account.GetOrdersHistory(ctx, []string{}, nil, nil, nil, nil)
	trades, _ := client.Account.GetTrades(ctx, []string{}, nil, nil, nil, nil)
	fees, _ := client.Account.GetFees(ctx, []string{"BTC-USD"}, nil)
	leverage, _ := client.Account.GetLeverage(ctx, []string{"BTC-USD"})
	bridgeConfig, _ := client.Account.GetBridgeConfig(ctx)

	fmt.Printf("%+v\n%+v\n%+v\n%+v\n%+v\n%+v\n%+v\n%+v\n%+v\n",
		account, balance, positions, openOrders, ordersHistory, trades, fees, leverage, bridgeConfig)

	// update leverage for a market
	newLeverage := decimal.NewFromFloat(10.0)
	client.Account.UpdateLeverage(ctx, "BTC-USD", newLeverage)
}
