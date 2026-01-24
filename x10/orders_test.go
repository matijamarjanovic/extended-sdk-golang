package x10

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/extended-protocol/extended-sdk-golang/x10/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testOrderPlacementAndCancellation tests the complete order lifecycle
func TestOrderPlacementAndCancellation(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	markets, err := client.Markets.GetMarkets(ctx, []string{"BTC-USD"})
	require.NoError(t, err, "should be able to get BTC-USD market")
	require.Greater(t, len(markets), 0, "should have at least one market")
	market := markets[0]

	// track all placed orders for cleanup
	type placedOrder struct {
		orderID    uint
		externalID string
	}
	var placedOrders []placedOrder

	// cleanup function to cancel all orders at the end
	defer func() {
		t.Logf("cleaning up %d tracked orders", len(placedOrders))

		for _, order := range placedOrders {
			err := client.Orders.CancelOrderByExternalID(ctx, order.externalID)
			if err != nil {
				t.Logf("failed to cancel order by external ID %s, trying internal ID %d: %v", order.externalID, order.orderID, err)
				err = client.Orders.CancelOrder(ctx, int(order.orderID))
				if err != nil {
					t.Logf("failed to cancel order %d: %v (may already be canceled)", order.orderID, err)
				}
			}
		}

		time.Sleep(1 * time.Second)

		openOrders, err := client.Account.GetOpenOrders(ctx, []string{"BTC-USD"}, models.OrderTypeAll, models.OrderSideAll)
		if err == nil && len(openOrders) > 0 {
			t.Logf("found %d remaining open orders after cleanup, canceling them all", len(openOrders))

			var orderIDs []int

			for _, openOrder := range openOrders {
				orderIDs = append(orderIDs, openOrder.ID)
				t.Logf("found remaining order: ID=%d, ExternalID=%s", openOrder.ID, openOrder.ExternalID)
			}

			if len(orderIDs) > 0 {
				err := client.Orders.MassCancel(ctx, orderIDs, nil, nil, false)
				if err != nil {
					t.Logf("mass cancel by IDs failed, trying individual cancellations: %v", err)
					for _, openOrder := range openOrders {
						err := client.Orders.CancelOrderByExternalID(ctx, openOrder.ExternalID)
						if err != nil {
							t.Logf("failed to cancel remaining order %s (ID: %d): %v", openOrder.ExternalID, openOrder.ID, err)
						}
					}
				} else {
					t.Logf("successfully mass canceled %d remaining orders", len(orderIDs))
				}
			}
		} else if err != nil {
			t.Logf("could not check for remaining open orders: %v", err)
		} else {
			t.Logf("no remaining open orders found - cleanup successful")
		}
	}()

	// test 1: place a basic BUY order
	t.Run("PlaceBasicBuyOrder", func(t *testing.T) {
		expireTime := time.Now().Add(1 * time.Hour)

		response, err := client.Orders.PlaceOrder(ctx,
			market,
			decimal.NewFromFloat(0.0001),
			decimal.NewFromFloat(1000),
			OrderSideBuy,
			models.OrderTypeLimit,
			TimeInForceGTT,
			SelfTradeProtectionDisabled,
			WithExpireTime(expireTime),
		)

		require.NoError(t, err, "should not error when placing valid order")
		require.NotNil(t, response, "response should not be nil")
		require.Equal(t, "OK", response.Status, "response status should be OK")

		require.NotZero(t, response.Data.OrderID, "order ID should be non-zero")
		require.NotEmpty(t, response.Data.ExternalID, "external ID should not be empty")

		t.Logf("placed BUY order - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)

		placedOrders = append(placedOrders, placedOrder{
			orderID:    response.Data.OrderID,
			externalID: response.Data.ExternalID,
		})
	})

	// test 2: place a SELL order with post-only
	t.Run("PlaceSellOrderWithPostOnly", func(t *testing.T) {
		expireTime := time.Now().Add(1 * time.Hour)

		response, err := client.Orders.PlaceOrder(ctx,
			market,
			decimal.NewFromFloat(0.001),
			decimal.NewFromFloat(100000),
			OrderSideSell,
			models.OrderTypeLimit,
			TimeInForceGTT,
			SelfTradeProtectionAccount,
			WithExpireTime(expireTime),
			WithPostOnly(true),
		)

		require.NoError(t, err, "should not error when placing SELL order with post-only")
		require.NotNil(t, response, "response should not be nil")
		require.Equal(t, "OK", response.Status, "response status should be OK")
		require.NotZero(t, response.Data.OrderID, "order ID should be non-zero")
		require.NotEmpty(t, response.Data.ExternalID, "external ID should not be empty")

		t.Logf("placed SELL order (post-only) - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)

		placedOrders = append(placedOrders, placedOrder{
			orderID:    response.Data.OrderID,
			externalID: response.Data.ExternalID,
		})
	})

	// test 3: place an order with custom external ID and custom nonce
	t.Run("PlaceOrderWithCustomExternalID", func(t *testing.T) {
		nonce := int(time.Now().UnixNano())
		expireTime := time.Now().Add(1 * time.Hour)
		customExternalID := fmt.Sprintf("test-order-%d", nonce)

		response, err := client.Orders.PlaceOrder(ctx,
			market,
			decimal.NewFromFloat(0.0001),
			decimal.NewFromFloat(1000),
			OrderSideBuy,
			models.OrderTypeLimit,
			TimeInForceGTT,
			SelfTradeProtectionClient,
			WithExpireTime(expireTime),
			WithOrderExternalID(customExternalID),
			WithNonce(nonce),
		)

		require.NoError(t, err, "should not error when placing order with custom external ID")
		require.NotNil(t, response, "response should not be nil")
		require.Equal(t, "OK", response.Status, "response status should be OK")
		require.Equal(t, customExternalID, response.Data.ExternalID, "external ID should match custom value")

		t.Logf("placed order with custom external ID - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)

		placedOrders = append(placedOrders, placedOrder{
			orderID:    response.Data.OrderID,
			externalID: response.Data.ExternalID,
		})
	})

	// test 4: place order with previous order external ID (order replacement)
	t.Run("PlaceOrderWithPreviousOrderExternalID", func(t *testing.T) {
		if len(placedOrders) == 0 {
			t.Skip("no previous orders to replace")
		}

		previousOrderID := placedOrders[0].externalID
		expireTime := time.Now().Add(1 * time.Hour)

		response, err := client.Orders.PlaceOrder(ctx,
			market,
			decimal.NewFromFloat(0.0001),
			decimal.NewFromFloat(1000),
			OrderSideBuy,
			models.OrderTypeLimit,
			TimeInForceGTT,
			SelfTradeProtectionDisabled,
			WithExpireTime(expireTime),
			WithPreviousOrderExternalID(previousOrderID),
		)

		require.NoError(t, err, "should not error when placing order with previous order external ID")
		require.NotNil(t, response, "response should not be nil")
		require.Equal(t, "OK", response.Status, "response status should be OK")

		t.Logf("placed order replacing previous order %s - New ID: %d, External ID: %s", previousOrderID, response.Data.OrderID, response.Data.ExternalID)

		placedOrders = append(placedOrders, placedOrder{
			orderID:    response.Data.OrderID,
			externalID: response.Data.ExternalID,
		})
	})

	// test 5: place order with builder fee
	t.Run("PlaceOrderWithBuilderFee", func(t *testing.T) {
		expireTime := time.Now().Add(1 * time.Hour)
		builderFee := decimal.NewFromFloat(0.0001)
		builderID := 2017

		response, err := client.Orders.PlaceOrder(ctx,
			market,
			decimal.NewFromFloat(0.0001),
			decimal.NewFromFloat(1000),
			OrderSideBuy,
			models.OrderTypeLimit,
			TimeInForceGTT,
			SelfTradeProtectionDisabled,
			WithExpireTime(expireTime),
			WithBuilderFee(builderFee),
			WithBuilderID(builderID),
		)

		require.NoError(t, err, "should not error when placing order with builder fee")
		require.NotNil(t, response, "response should not be nil")
		require.Equal(t, "OK", response.Status, "response status should be OK")

		t.Logf("placed order with builder fee - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)

		placedOrders = append(placedOrders, placedOrder{
			orderID:    response.Data.OrderID,
			externalID: response.Data.ExternalID,
		})
	})

	// test 6: place order with builder ID
	t.Run("PlaceOrderWithBuilderID", func(t *testing.T) {
		expireTime := time.Now().Add(1 * time.Hour)
		builderID := 2017

		response, err := client.Orders.PlaceOrder(ctx,
			market,
			decimal.NewFromFloat(0.0001),
			decimal.NewFromFloat(1000),
			OrderSideBuy,
			models.OrderTypeLimit,
			TimeInForceGTT,
			SelfTradeProtectionDisabled,
			WithExpireTime(expireTime),
			WithBuilderID(builderID),
		)

		require.NoError(t, err, "should not error when placing order with builder ID")
		require.NotNil(t, response, "response should not be nil")
		require.Equal(t, "OK", response.Status, "response status should be OK")

		t.Logf("placed order with builder ID - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)

		placedOrders = append(placedOrders, placedOrder{
			orderID:    response.Data.OrderID,
			externalID: response.Data.ExternalID,
		})
	})

	// test 7: place order with builder fee and builder ID
	t.Run("PlaceOrderWithBuilderFeeAndID", func(t *testing.T) {
		expireTime := time.Now().Add(1 * time.Hour)
		builderFee := decimal.NewFromFloat(0.0001)
		builderID := 2017

		response, err := client.Orders.PlaceOrder(ctx,
			market,
			decimal.NewFromFloat(0.0001),
			decimal.NewFromFloat(1000),
			OrderSideBuy,
			models.OrderTypeLimit,
			TimeInForceGTT,
			SelfTradeProtectionDisabled,
			WithExpireTime(expireTime),
			WithBuilderFee(builderFee),
			WithBuilderID(builderID),
		)

		require.NoError(t, err, "should not error when placing order with builder fee and ID")
		require.NotNil(t, response, "response should not be nil")
		require.Equal(t, "OK", response.Status, "response status should be OK")

		t.Logf("placed order with builder fee and ID - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)

		placedOrders = append(placedOrders, placedOrder{
			orderID:    response.Data.OrderID,
			externalID: response.Data.ExternalID,
		})
	})

	// test 8: place order with all optional parameters combined
	t.Run("PlaceOrderWithAllOptionalParams", func(t *testing.T) {
		expireTime := time.Now().Add(1 * time.Hour)

		customExternalID := fmt.Sprintf("test-all-params-%d", time.Now().UnixNano())
		builderFee := decimal.NewFromFloat(0.0001)
		builderID := 2017

		response, err := client.Orders.PlaceOrder(ctx,
			market,
			decimal.NewFromFloat(0.0001),
			decimal.NewFromFloat(1000),
			OrderSideBuy,
			models.OrderTypeLimit,
			TimeInForceGTT,
			SelfTradeProtectionDisabled,
			WithExpireTime(expireTime),
			WithPostOnly(true),
			WithOrderExternalID(customExternalID),
			WithBuilderFee(builderFee),
			WithBuilderID(builderID),
		)

		require.NoError(t, err, "should not error when placing order with all optional parameters")
		require.NotNil(t, response, "response should not be nil")
		require.Equal(t, "OK", response.Status, "response status should be OK")
		require.Equal(t, customExternalID, response.Data.ExternalID, "External ID should match custom value")

		t.Logf("placed order with all optional params - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)

		placedOrders = append(placedOrders, placedOrder{
			orderID:    response.Data.OrderID,
			externalID: response.Data.ExternalID,
		})
	})

	// test 9: place MARKET order to create a position (for reduce-only testing)
	t.Run("PlaceMarketOrderToCreatePosition", func(t *testing.T) {
		stats, err := client.Markets.GetMarketStatistics(ctx, "BTC-USD")
		require.NoError(t, err, "should be able to get market stats")

		orderbook, err := client.Markets.GetOrderbookSnapshot(ctx, "BTC-USD")

		var orderPrice decimal.Decimal
		if err == nil && orderbook != nil && len(orderbook.Ask) > 0 {
			orderPrice = orderbook.Ask[1].Price
			t.Logf("Using best ask price from orderbook: %s", orderPrice.String())
		} else {
			if !stats.AskPrice.IsZero() {
				orderPrice = stats.AskPrice
			} else {
				orderPrice = stats.MarkPrice
			}
			t.Logf("Using ask/mark price from stats: %s", orderPrice.String())
		}

		minOrderSize := decimal.NewFromFloat(0.0001)
		expireTime := time.Now().Add(1 * time.Hour)
		response, err := client.Orders.PlaceOrder(ctx,
			market,
			minOrderSize,
			orderPrice,
			OrderSideBuy,
			models.OrderTypeMarket,
			models.TimeInForceIOC,
			SelfTradeProtectionDisabled,
			WithExpireTime(expireTime),
		)

		require.NoError(t, err, "should not error when placing MARKET order")
		require.NotNil(t, response, "response should not be nil")
		require.Equal(t, "OK", response.Status, "response status should be OK")
		require.NotZero(t, response.Data.OrderID, "order ID should not be zero")
		require.NotEmpty(t, response.Data.ExternalID, "external ID should not be empty")

		t.Logf("placed MARKET order - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)

		time.Sleep(5 * time.Second)

		// verify order was filled by checking trades
		trades, err := client.Account.GetTrades(ctx, []string{"BTC-USD"}, models.OrderSideAll, models.TradeTypeAll, 0, 0)
		require.NoError(t, err, "should be able to get trades")

		var foundTrade *models.AccountTradeModel
		for i := range trades {
			if trades[i].OrderID == int(response.Data.OrderID) {
				foundTrade = &trades[i]
				break
			}
		}

		require.NotNil(t, foundTrade, "should find trade for placed order")
		require.Equal(t, int(response.Data.OrderID), foundTrade.OrderID, "trade order ID should match")
	})

	// test 10: place reduce-only order (requires open position)
	t.Run("PlaceReduceOnlyOrder", func(t *testing.T) {
		positions, err := client.Account.GetPositions(ctx, []string{"BTC-USD"}, models.PositionSideAll)
		if err != nil {
			t.Skipf("cannot test reduce-only: could not check positions: %v", err)
			return
		}

		if len(positions) == 0 {
			t.Skip("cannot test reduce-only: no open position exists")
			return
		}

		position := positions[0]
		t.Logf("found position: %s, side: %s, size: %s", position.Market, position.Side, position.Size.String())

		var reduceSide models.OrderSide
		if position.Side == models.PositionSideLong {
			reduceSide = OrderSideSell
		} else {
			reduceSide = OrderSideBuy
		}

		stats, err := client.Markets.GetMarketStatistics(ctx, "BTC-USD")
		require.NoError(t, err, "should be able to get market stats")

		orderbook, err := client.Markets.GetOrderbookSnapshot(ctx, "BTC-USD")

		var reducePrice decimal.Decimal
		if reduceSide == OrderSideSell {
			if err == nil && orderbook != nil && len(orderbook.Bid) > 0 {
				reducePrice = orderbook.Bid[1].Price
				t.Logf("using best bid price from orderbook: %s", reducePrice.String())
			} else {
				if !stats.BidPrice.IsZero() {
					reducePrice = stats.BidPrice
				} else {
					reducePrice = stats.MarkPrice
				}
				t.Logf("using bid/mark price from stats: %s", reducePrice.String())
			}
		} else {
			if err == nil && orderbook != nil && len(orderbook.Ask) > 0 {
				reducePrice = orderbook.Ask[1].Price
				t.Logf("using best ask price from orderbook: %s", reducePrice.String())
			} else {
				if !stats.AskPrice.IsZero() {
					reducePrice = stats.AskPrice
				} else {
					reducePrice = stats.MarkPrice
				}
				t.Logf("using ask/mark price from stats: %s", reducePrice.String())
			}
		}

		reduceAmount := decimal.NewFromFloat(0.0001)
		expireTime := time.Now().Add(1 * time.Hour)
		response, err := client.Orders.PlaceOrder(ctx,
			market,
			reduceAmount,
			reducePrice,
			reduceSide,
			models.OrderTypeMarket,
			models.TimeInForceIOC,
			SelfTradeProtectionDisabled,
			WithExpireTime(expireTime),
			WithReduceOnly(true),
		)

		require.NoError(t, err, "should not error when placing order with reduce-only")
		require.NotNil(t, response, "response should not be nil")
		require.Equal(t, "OK", response.Status, "response status should be OK")
		require.NotZero(t, response.Data.OrderID, "order ID should not be zero")
		require.NotEmpty(t, response.Data.ExternalID, "external ID should not be empty")

		t.Logf("placed reduce-only MARKET order - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)

		time.Sleep(5 * time.Second)

		trades, err := client.Account.GetTrades(ctx, []string{"BTC-USD"}, models.OrderSideAll, models.TradeTypeAll, 0, 0)
		require.NoError(t, err, "should be able to get trades")

		var foundTrade *models.AccountTradeModel
		for i := range trades {
			if trades[i].OrderID == int(response.Data.OrderID) {
				foundTrade = &trades[i]
				break
			}
		}

		require.NotNil(t, foundTrade, "should find trade for placed order")
		require.Equal(t, int(response.Data.OrderID), foundTrade.OrderID, "trade order ID should match")

		positionsAfter, posErr := client.Account.GetPositions(ctx, []string{"BTC-USD"}, models.PositionSideAll)
		require.NoError(t, posErr, "should be able to check positions after reduce-only order")

		if len(positionsAfter) > 0 {
			newPosition := positionsAfter[0]
			require.Equal(t, position.Market, newPosition.Market, "position market should match")
			require.Equal(t, position.Side, newPosition.Side, "position side should match")
			require.Less(t, newPosition.Size.Cmp(position.Size), 0, "position size MUST be reduced (was: %s, now: %s)", position.Size.String(), newPosition.Size.String())
			t.Logf("position reduced: %s -> %s", position.Size.String(), newPosition.Size.String())
		} else {
			t.Logf("position was fully closed by reduce-only order (original size: %s)", position.Size.String())
		}
	})

	// test 13: cancel order by internal ID
	t.Run("CancelOrderByInternalID", func(t *testing.T) {
		if len(placedOrders) == 0 {
			t.Skip("no orders placed to cancel")
		}

		orderToCancel := placedOrders[0]

		err := client.Orders.CancelOrder(ctx, int(orderToCancel.orderID))
		require.NoError(t, err, "should not error when canceling order by internal ID")

		time.Sleep(500 * time.Millisecond)
		openOrders, err := client.Account.GetOpenOrders(ctx, []string{"BTC-USD"}, models.OrderTypeAll, models.OrderSideAll)
		if err == nil {
			for _, openOrder := range openOrders {
				if openOrder.ExternalID == orderToCancel.externalID {
					t.Errorf("order %s should have been canceled but still appears in open orders", orderToCancel.externalID)
				}
			}
		}

		placedOrders = placedOrders[1:]

		t.Logf("successfully canceled order by internal ID: %d", orderToCancel.orderID)
	})

	// test 6: cancel order by external ID
	t.Run("CancelOrderByExternalID", func(t *testing.T) {
		if len(placedOrders) == 0 {
			t.Skip("no orders placed to cancel")
		}

		orderToCancel := placedOrders[0]

		err := client.Orders.CancelOrderByExternalID(ctx, orderToCancel.externalID)
		require.NoError(t, err, "should not error when canceling order by external ID")

		time.Sleep(500 * time.Millisecond)
		openOrders, err := client.Account.GetOpenOrders(ctx, []string{"BTC-USD"}, models.OrderTypeAll, models.OrderSideAll)
		if err == nil {
			for _, openOrder := range openOrders {
				if openOrder.ExternalID == orderToCancel.externalID {
					t.Errorf("order %s should have been canceled but still appears in open orders", orderToCancel.externalID)
				}
			}
		}

		placedOrders = placedOrders[1:]

		t.Logf("successfully canceled order by external ID: %s", orderToCancel.externalID)
	})

	// test 7: verify order response structure
	t.Run("VerifyOrderResponseStructure", func(t *testing.T) {
		expireTime := time.Now().Add(1 * time.Hour)

		response, err := client.Orders.PlaceOrder(ctx,
			market,
			decimal.NewFromFloat(0.0001),
			decimal.NewFromFloat(1000),
			OrderSideBuy,
			models.OrderTypeLimit,
			TimeInForceGTT,
			SelfTradeProtectionDisabled,
			WithExpireTime(expireTime),
		)

		require.NoError(t, err, "should not error when placing order")

		assert.Equal(t, "OK", response.Status, "status should be OK")
		assert.NotZero(t, response.Data.OrderID, "order ID (id) should be non-zero integer")
		assert.NotEmpty(t, response.Data.ExternalID, "external ID (external_id) should be non-empty string")
		orders, err := client.Account.GetOrderByExternalID(ctx, response.Data.ExternalID)
		if err == nil && len(orders) > 0 {
			assert.Equal(t, response.Data.ExternalID, orders[0].ExternalID, "retrieved order should match placed order")
			assert.Equal(t, int(response.Data.OrderID), orders[0].ID, "retrieved order ID should match")
		}

		placedOrders = append(placedOrders, placedOrder{
			orderID:    response.Data.OrderID,
			externalID: response.Data.ExternalID,
		})

		t.Logf("verified order response structure - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)
	})

	// test 8: test mass cancel
	t.Run("MassCancelOrders", func(t *testing.T) {
		var orderIDs []int
		var externalIDs []string

		for i := 0; i < 2; i++ {
			expireTime := time.Now().Add(1 * time.Hour)

			response, err := client.Orders.PlaceOrder(ctx,
				market,
				decimal.NewFromFloat(0.0001),
				decimal.NewFromFloat(1000),
				OrderSideBuy,
				models.OrderTypeLimit,
				TimeInForceGTT,
				SelfTradeProtectionDisabled,
				WithExpireTime(expireTime),
			)

			require.NoError(t, err, "should not error when placing order for mass cancel")
			orderIDs = append(orderIDs, int(response.Data.OrderID))
			externalIDs = append(externalIDs, response.Data.ExternalID)

			placedOrders = append(placedOrders, placedOrder{
				orderID:    response.Data.OrderID,
				externalID: response.Data.ExternalID,
			})
		}

		err := client.Orders.MassCancel(ctx, orderIDs, nil, nil, false)
		require.NoError(t, err, "should not error when mass canceling orders")

		time.Sleep(500 * time.Millisecond)
		openOrders, err := client.Account.GetOpenOrders(ctx, []string{"BTC-USD"}, models.OrderTypeAll, models.OrderSideAll)
		if err == nil {
			for _, externalID := range externalIDs {
				for _, openOrder := range openOrders {
					if openOrder.ExternalID == externalID {
						t.Errorf("order %s should have been canceled by mass cancel", externalID)
					}
				}
			}
		}

		for range orderIDs {
			if len(placedOrders) > 0 {
				placedOrders = placedOrders[1:]
			}
		}

		t.Logf("successfully mass canceled %d orders", len(orderIDs))
	})

	t.Logf("all order placement and cancellation tests completed. %d orders remaining for cleanup.", len(placedOrders))
}

// testOrderPlacementWithTPSL tests placing orders with take profit and stop loss
func TestOrderPlacementWithTPSL(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	markets, err := client.Markets.GetMarkets(ctx, []string{"BTC-USD"})
	require.NoError(t, err, "should be able to get BTC-USD market")
	require.Greater(t, len(markets), 0, "should have at least one market")
	market := markets[0]

	type placedOrder struct {
		orderID    uint
		externalID string
	}
	var placedOrders []placedOrder

	defer func() {
		t.Logf("cleaning up %d TPSL orders", len(placedOrders))
		for _, order := range placedOrders {
			err := client.Orders.CancelOrderByExternalID(ctx, order.externalID)
			if err != nil {
				t.Logf("failed to cancel TPSL order %s: %v", order.externalID, err)
			}
		}
		time.Sleep(1 * time.Second)
	}()

	roundPrice := func(price decimal.Decimal) decimal.Decimal {
		minPriceChange := decimal.NewFromFloat(0.1)
		if market.TradingConfig != nil && !market.TradingConfig.MinPriceChange.IsZero() {
			minPriceChange = market.TradingConfig.MinPriceChange
		}
		return price.Div(minPriceChange).Round(0).Mul(minPriceChange)
	}

	// test 1: place order with take profit only
	t.Run("PlaceOrderWithTakeProfit", func(t *testing.T) {
		expireTime := time.Now().Add(1 * time.Hour)

		stats, err := client.Markets.GetMarketStatistics(ctx, "BTC-USD")
		require.NoError(t, err, "should be able to get market stats")

		basePrice := stats.MarkPrice
		if basePrice.IsZero() {
			basePrice = decimal.NewFromFloat(50000)
		}
		basePrice = roundPrice(basePrice)
		tpPrice := roundPrice(basePrice.Mul(decimal.NewFromFloat(1.05)))
		tpTriggerPrice := roundPrice(basePrice.Mul(decimal.NewFromFloat(1.03)))
		orderPrice := roundPrice(basePrice.Mul(decimal.NewFromFloat(0.95)))

		response, err := client.Orders.PlaceOrder(ctx,
			market,
			decimal.NewFromFloat(0.0001),
			orderPrice,
			OrderSideBuy,
			models.OrderTypeLimit,
			TimeInForceGTT,
			SelfTradeProtectionDisabled,
			WithExpireTime(expireTime),
			WithTpSlType(models.TpSlTypeOrder),
			WithTakeProfit(models.TpSlTriggerParam{
				TriggerPrice:     tpTriggerPrice,
				TriggerPriceType: models.TriggerPriceTypeMark,
				Price:            tpPrice,
				PriceType:        models.ExecutionPriceTypeLimit,
			}),
		)

		require.NoError(t, err, "should not error when placing order with take profit")
		require.NotNil(t, response, "response should not be nil")
		require.Equal(t, "OK", response.Status, "response status should be OK")
		require.NotZero(t, response.Data.OrderID, "order ID should be non-zero")
		require.NotEmpty(t, response.Data.ExternalID, "external ID should not be empty")

		t.Logf("placed order with take profit - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)

		placedOrders = append(placedOrders, placedOrder{
			orderID:    response.Data.OrderID,
			externalID: response.Data.ExternalID,
		})

		time.Sleep(500 * time.Millisecond)
		orders, err := client.Account.GetOrderByExternalID(ctx, response.Data.ExternalID)
		if err == nil && len(orders) > 0 {
			order := orders[0]
			require.NotNil(t, order.TakeProfit, "order should have take profit set")
			require.NotNil(t, order.TpSlType, "order should have TPSL type set")
			require.Equal(t, models.TpSlTypeOrder, *order.TpSlType, "TPSL type should be ORDER")
			t.Logf("verified order has take profit: trigger=%s, price=%s",
				order.TakeProfit.TriggerPrice.String(), order.TakeProfit.Price.String())
		}
	})

	// test 2: place order with stop loss only
	t.Run("PlaceOrderWithStopLoss", func(t *testing.T) {
		expireTime := time.Now().Add(1 * time.Hour)

		stats, err := client.Markets.GetMarketStatistics(ctx, "BTC-USD")
		require.NoError(t, err, "should be able to get market stats")

		basePrice := stats.MarkPrice
		if basePrice.IsZero() {
			basePrice = decimal.NewFromFloat(50000)
		}
		basePrice = roundPrice(basePrice)
		orderPrice := roundPrice(basePrice.Mul(decimal.NewFromFloat(0.95)))
		slPrice := roundPrice(orderPrice.Mul(decimal.NewFromFloat(0.98)))
		slTriggerPrice := roundPrice(orderPrice.Mul(decimal.NewFromFloat(0.99)))

		response, err := client.Orders.PlaceOrder(ctx,
			market,
			decimal.NewFromFloat(0.0001),
			orderPrice,
			OrderSideBuy,
			models.OrderTypeLimit,
			TimeInForceGTT,
			SelfTradeProtectionDisabled,
			WithExpireTime(expireTime),
			WithTpSlType(models.TpSlTypeOrder),
			WithStopLoss(models.TpSlTriggerParam{
				TriggerPrice:     slTriggerPrice,
				TriggerPriceType: models.TriggerPriceTypeMark,
				Price:            slPrice,
				PriceType:        models.ExecutionPriceTypeLimit,
			}),
		)

		require.NoError(t, err, "should not error when placing order with stop loss")
		require.NotNil(t, response, "response should not be nil")
		require.Equal(t, "OK", response.Status, "response status should be OK")
		require.NotZero(t, response.Data.OrderID, "order ID should be non-zero")

		t.Logf("placed order with stop loss - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)

		placedOrders = append(placedOrders, placedOrder{
			orderID:    response.Data.OrderID,
			externalID: response.Data.ExternalID,
		})

		time.Sleep(500 * time.Millisecond)
		orders, err := client.Account.GetOrderByExternalID(ctx, response.Data.ExternalID)
		if err == nil && len(orders) > 0 {
			order := orders[0]
			require.NotNil(t, order.StopLoss, "order should have stop loss set")
			t.Logf("verified order has stop loss: trigger=%s, price=%s",
				order.StopLoss.TriggerPrice.String(), order.StopLoss.Price.String())
		}
	})

	// test 3: place order with both take profit and stop loss
	t.Run("PlaceOrderWithBothTPSL", func(t *testing.T) {
		expireTime := time.Now().Add(1 * time.Hour)

		stats, err := client.Markets.GetMarketStatistics(ctx, "BTC-USD")
		require.NoError(t, err, "should be able to get market stats")

		basePrice := stats.MarkPrice
		if basePrice.IsZero() {
			basePrice = decimal.NewFromFloat(50000)
		}
		basePrice = roundPrice(basePrice)
		orderPrice := roundPrice(basePrice.Mul(decimal.NewFromFloat(0.95)))

		tpPrice := roundPrice(orderPrice.Mul(decimal.NewFromFloat(1.05)))
		tpTriggerPrice := roundPrice(orderPrice.Mul(decimal.NewFromFloat(1.03)))
		slPrice := roundPrice(orderPrice.Mul(decimal.NewFromFloat(0.98)))
		slTriggerPrice := roundPrice(orderPrice.Mul(decimal.NewFromFloat(0.99)))

		response, err := client.Orders.PlaceOrder(ctx,
			market,
			decimal.NewFromFloat(0.0001),
			orderPrice,
			OrderSideBuy,
			models.OrderTypeLimit,
			TimeInForceGTT,
			SelfTradeProtectionDisabled,
			WithExpireTime(expireTime),
			WithTpSlType(models.TpSlTypeOrder),
			WithTakeProfit(models.TpSlTriggerParam{
				TriggerPrice:     tpTriggerPrice,
				TriggerPriceType: models.TriggerPriceTypeMark,
				Price:            tpPrice,
				PriceType:        models.ExecutionPriceTypeLimit,
			}),
			WithStopLoss(models.TpSlTriggerParam{
				TriggerPrice:     slTriggerPrice,
				TriggerPriceType: models.TriggerPriceTypeMark,
				Price:            slPrice,
				PriceType:        models.ExecutionPriceTypeLimit,
			}),
		)

		require.NoError(t, err, "should not error when placing order with both TPSL")
		require.NotNil(t, response, "response should not be nil")
		require.Equal(t, "OK", response.Status, "response status should be OK")
		require.NotZero(t, response.Data.OrderID, "order ID should be non-zero")

		t.Logf("placed order with both TPSL - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)

		placedOrders = append(placedOrders, placedOrder{
			orderID:    response.Data.OrderID,
			externalID: response.Data.ExternalID,
		})

		time.Sleep(500 * time.Millisecond)
		orders, err := client.Account.GetOrderByExternalID(ctx, response.Data.ExternalID)
		if err == nil && len(orders) > 0 {
			order := orders[0]
			require.NotNil(t, order.TakeProfit, "order should have take profit set")
			require.NotNil(t, order.StopLoss, "order should have stop loss set")
			t.Logf("verified order has both TPSL: TP trigger=%s, SL trigger=%s",
				order.TakeProfit.TriggerPrice.String(), order.StopLoss.TriggerPrice.String())
		}
	})
}

// testOrderPlacementErrorHandling tests error cases for order placement
func TestOrderPlacementErrorHandling(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	markets, err := client.Markets.GetMarkets(ctx, []string{"BTC-USD"})
	require.NoError(t, err, "should be able to get markets")
	require.Greater(t, len(markets), 0, "should have at least one market")
	market := markets[0]

	t.Run("InvalidMarket", func(t *testing.T) {
		invalidMarket := models.MarketModel{
			Name:     "NONEXISTENT-MARKET-12345",
			L2Config: market.L2Config,
		}

		_, err := client.Orders.PlaceOrder(ctx,
			invalidMarket,
			decimal.NewFromFloat(0.0001),
			decimal.NewFromFloat(1000),
			OrderSideBuy,
			models.OrderTypeLimit,
			TimeInForceGTT,
			SelfTradeProtectionDisabled,
		)

		require.Error(t, err, "should error when placing order with invalid market name")
		t.Logf("got expected error for invalid market: %v", err)
	})

	t.Run("ZeroAmount", func(t *testing.T) {
		_, err := client.Orders.PlaceOrder(ctx,
			market,
			decimal.Zero,
			decimal.NewFromFloat(1),
			OrderSideBuy,
			models.OrderTypeLimit,
			TimeInForceGTT,
			SelfTradeProtectionDisabled,
		)

		if err == nil {
			t.Log("note: zero amount did not error immediately (may fail at API)")
		}
	})
}
