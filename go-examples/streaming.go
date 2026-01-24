package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/extended-protocol/extended-sdk-golang/x10"
	"github.com/extended-protocol/extended-sdk-golang/x10/models"
)

//TODO: MAKE ALIASES FOR MODELS

func streamingExample(client *x10.Client) {
	ctx := context.Background()

	// subscribe to various websocket streams
	orderbookStream, _ := client.Streaming.SubscribeToOrderbooks(ctx, "BTC-USD", 10)
	defer orderbookStream.Close()

	tradesStream, _ := client.Streaming.SubscribeToPublicTrades(ctx, "BTC-USD")
	defer tradesStream.Close()

	fundingStream, _ := client.Streaming.SubscribeToFundingRates(ctx, "BTC-USD")
	defer fundingStream.Close()

	candlesStream, _ := client.Streaming.SubscribeToCandles(
		ctx,
		"BTC-USD",
		x10.CandleTypeTrades,
		x10.CandleIntervalPT1M,
	)
	defer candlesStream.Close()

	accountStream, _ := client.Streaming.SubscribeToAccountUpdates(ctx)
	defer accountStream.Close()

	markPricesStream, _ := client.Streaming.SubscribeToMarkPrices(ctx, "BTC-USD")
	defer markPricesStream.Close()

	indexPricesStream, _ := client.Streaming.SubscribeToIndexPrices(ctx, "BTC-USD")
	defer indexPricesStream.Close()

	// receive messages from streams
	var orderbookUpdate models.OrderbookUpdateModel
	orderbookStream.Recv(ctx, &orderbookUpdate)

	var trade models.StreamPublicTradeModel
	tradesStream.Recv(ctx, &trade)

	var fundingRate models.FundingRateModel
	fundingStream.Recv(ctx, &fundingRate)

	var candle models.CandleModel
	candlesStream.Recv(ctx, &candle)

	var accountUpdate models.AccountStreamDataModel
	accountStream.Recv(ctx, &accountUpdate)

	var markPrice models.MarkPriceModel
	markPricesStream.Recv(ctx, &markPrice)

	var indexPrice models.IndexPriceModel
	indexPricesStream.Recv(ctx, &indexPrice)

	time.Sleep(100 * time.Millisecond)

	// print all stream data as json
	results := []interface{}{
		orderbookUpdate, trade, fundingRate, candle, accountUpdate, markPrice, indexPrice,
	}
	jsonData, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(jsonData))
}
