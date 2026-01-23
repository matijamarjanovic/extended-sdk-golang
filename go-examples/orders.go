package main

import (
	"context"
	"fmt"
	"time"

	"github.com/extended-protocol/extended-sdk-golang/x10"
	"github.com/extended-protocol/extended-sdk-golang/x10/models"
	"github.com/shopspring/decimal"
)

func ordersExample(client *x10.Client) {
	ctx := context.Background()

	// get market info
	markets, _ := client.Markets.GetMarkets(ctx, []string{"BTC-USD"})
	market := markets[0]

	// place a basic limit order
	expireTime := time.Now().Add(1 * time.Hour)
	order, _ := client.Orders.PlaceOrder(
		ctx,
		market,
		decimal.NewFromFloat(0.001),
		decimal.NewFromFloat(50000),
		models.OrderSideBuy,
		models.OrderTypeLimit,
		models.TimeInForceGTT,
		models.SelfTradeProtectionDisabled,
		x10.WithExpireTime(expireTime),
		x10.WithPostOnly(false),
	)

	// place order with additional options
	nonce := 12345
	builderFee := decimal.NewFromFloat(0.0001)
	orderWithOptions, _ := client.Orders.PlaceOrder(
		ctx,
		market,
		decimal.NewFromFloat(0.001),
		decimal.NewFromFloat(50001),
		models.OrderSideBuy,
		models.OrderTypeLimit,
		models.TimeInForceGTT,
		models.SelfTradeProtectionDisabled,
		x10.WithPostOnly(true),
		x10.WithReduceOnly(false),
		x10.WithNonce(nonce),
		x10.WithOrderExternalID("custom-id"),
		x10.WithBuilderFee(builderFee),
	)

	fmt.Printf("%+v\n%+v\n", order, orderWithOptions)

	// cancel orders
	client.Orders.CancelOrder(ctx, int(order.Data.OrderID))
	client.Orders.CancelOrderByExternalID(ctx, orderWithOptions.Data.ExternalID)

	// mass cancel orders
	orderIDs := []int{1, 2, 3}
	externalIDs := []string{"ext1", "ext2"}
	marketsToCancel := []string{"BTC-USD"}
	client.Orders.MassCancel(ctx, orderIDs, externalIDs, marketsToCancel, false)
}
