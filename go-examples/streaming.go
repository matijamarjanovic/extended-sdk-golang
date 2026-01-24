package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/extended-protocol/extended-sdk-golang/x10"
)

func streamingExample(client *x10.Client) {
	ctx := context.Background()

	// subscribe to orderbook updates
	orderbookStream, _ := client.Streaming.SubscribeToOrderbooks(ctx, "BTC-USD", 10)
	defer orderbookStream.Close()

	// subscribe to public trades
	tradesStream, _ := client.Streaming.SubscribeToPublicTrades(ctx, "BTC-USD")
	defer tradesStream.Close()

	// receive messages from streams
	var orderbookUpdate x10.OrderbookUpdateModel
	orderbookStream.Recv(ctx, &orderbookUpdate)

	var trade x10.StreamPublicTradeModel
	tradesStream.Recv(ctx, &trade)

	time.Sleep(100 * time.Millisecond)

	// print stream data as json
	results := []interface{}{orderbookUpdate, trade}
	jsonData, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(jsonData))
}
