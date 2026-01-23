package main

import (
	"context"
	"fmt"
	"time"

	"github.com/extended-protocol/extended-sdk-golang/x10"
	"github.com/extended-protocol/extended-sdk-golang/x10/models"
)

func marketsExample(client *x10.Client) {
	ctx := context.Background()

	// get market data
	markets, _ := client.Markets.GetMarkets(ctx, []string{"BTC-USD", "ETH-USD"})
	marketsDict, _ := client.Markets.GetMarketsDict(ctx)
	stats, _ := client.Markets.GetMarketStatistics(ctx, "BTC-USD")
	orderbook, _ := client.Markets.GetOrderbookSnapshot(ctx, "BTC-USD")

	// get candle history
	limit := 100
	candles, _ := client.Markets.GetCandlesHistory(
		ctx,
		"BTC-USD",
		models.CandleTypeTrades,
		models.CandleIntervalPT1M,
		&limit,
		nil,
	)

	// get funding rates history
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)
	fundingRates, _ := client.Markets.GetFundingRatesHistory(ctx, "BTC-USD", startTime, endTime)

	fmt.Printf("%+v\n%+v\n%+v\n%+v\n%+v\n%+v\n",
		markets, marketsDict, stats, orderbook, candles, fundingRates)
}
