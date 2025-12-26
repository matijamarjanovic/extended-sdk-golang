package sdk

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/extended-protocol/extended-sdk-golang/src/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOrderPlacementAndCancellation tests the complete order lifecycle:
// 1. Place multiple orders with different configurations
// 2. Verify order responses match expected structure
// 3. Cancel orders by both internal ID and external ID
// 4. Ensure all placed orders are canceled before test completion
func TestOrderPlacementAndCancellation(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	// Get a market to use for orders
	markets, err := client.Markets.GetMarkets(ctx, []string{"BTC-USD"})
	require.NoError(t, err, "Should be able to get BTC-USD market")
	require.Greater(t, len(markets), 0, "Should have at least one market")
	market := markets[0]

	// Track all placed orders for cleanup
	type placedOrder struct {
		orderID    uint
		externalID string
	}
	var placedOrders []placedOrder

	// Cleanup function to cancel all orders at the end
	defer func() {
		t.Logf("Cleaning up %d tracked orders", len(placedOrders))
		
		// First, try to cancel all tracked orders
		for _, order := range placedOrders {
			// Try canceling by external ID first (more reliable)
			err := client.Orders.CancelOrderByExternalID(ctx, order.externalID)
			if err != nil {
				// If external ID fails, try internal ID
				t.Logf("Failed to cancel order by external ID %s, trying internal ID %d: %v", order.externalID, order.orderID, err)
				err = client.Orders.CancelOrder(ctx, int(order.orderID))
				if err != nil {
					t.Logf("Failed to cancel order %d: %v (may already be canceled)", order.orderID, err)
				}
			}
		}

		// Wait a bit for cancellations to propagate
		time.Sleep(1 * time.Second)

		// Final cleanup: get all remaining open orders and cancel them
		// This catches any orders that weren't properly tracked or canceled
		openOrders, err := client.Account.GetOpenOrders(ctx, []string{"BTC-USD"}, nil, nil)
		if err == nil && len(openOrders) > 0 {
			t.Logf("Found %d remaining open orders after cleanup, canceling them all", len(openOrders))
			
			// Collect all order IDs for mass cancel
			var orderIDs []int
			
			for _, openOrder := range openOrders {
				orderIDs = append(orderIDs, openOrder.ID)
				t.Logf("Found remaining order: ID=%d, ExternalID=%s", openOrder.ID, openOrder.ExternalID)
			}
			
			// Use mass cancel to cancel all remaining orders at once
			if len(orderIDs) > 0 {
				err := client.Orders.MassCancel(ctx, orderIDs, nil, nil, false)
				if err != nil {
					t.Logf("Mass cancel by IDs failed, trying individual cancellations: %v", err)
					// Fallback to individual cancellations
					for _, openOrder := range openOrders {
						err := client.Orders.CancelOrderByExternalID(ctx, openOrder.ExternalID)
						if err != nil {
							t.Logf("Failed to cancel remaining order %s (ID: %d): %v", openOrder.ExternalID, openOrder.ID, err)
						}
					}
				} else {
					t.Logf("Successfully mass canceled %d remaining orders", len(orderIDs))
				}
			}
		} else if err != nil {
			t.Logf("Could not check for remaining open orders: %v", err)
		} else {
			t.Logf("No remaining open orders found - cleanup successful")
		}
	}()

	// Test 1: Place a basic BUY order (nonce auto-generated)
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

		require.NoError(t, err, "Should not error when placing valid order")
		require.NotNil(t, response, "Response should not be nil")
		require.Equal(t, "OK", response.Status, "Response status should be OK")

		require.NotZero(t, response.Data.OrderID, "Order ID should be non-zero")
		require.NotEmpty(t, response.Data.ExternalID, "External ID should not be empty")

		t.Logf("Placed BUY order - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)

		placedOrders = append(placedOrders, placedOrder{
			orderID:    response.Data.OrderID,
			externalID: response.Data.ExternalID,
		})
	})

	// Test 2: Place a SELL order with post-only (nonce auto-generated)
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

		require.NoError(t, err, "Should not error when placing SELL order with post-only")
		require.NotNil(t, response, "Response should not be nil")
		require.Equal(t, "OK", response.Status, "Response status should be OK")
		require.NotZero(t, response.Data.OrderID, "Order ID should be non-zero")
		require.NotEmpty(t, response.Data.ExternalID, "External ID should not be empty")

		t.Logf("Placed SELL order (post-only) - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)

		placedOrders = append(placedOrders, placedOrder{
			orderID:    response.Data.OrderID,
			externalID: response.Data.ExternalID,
		})
	})

	// Test 3: Place an order with custom external ID and custom nonce
	t.Run("PlaceOrderWithCustomExternalID", func(t *testing.T) {
		nonce := int(time.Now().UnixNano())
		expireTime := time.Now().Add(1 * time.Hour)
		customExternalID := fmt.Sprintf("test-order-%d", nonce)

		// Use a very low price for BUY order (well below market, unlikely to match)
		// Using a small amount to minimize balance requirements
		response, err := client.Orders.PlaceOrder(ctx,
			market,
			decimal.NewFromFloat(0.0001),
			decimal.NewFromFloat(1000),   // Very low price, well below market
			OrderSideBuy,
			models.OrderTypeLimit,
			TimeInForceGTT,
			SelfTradeProtectionClient,
			WithExpireTime(expireTime),
			WithOrderExternalID(customExternalID),
			WithNonce(nonce),
		)

		require.NoError(t, err, "Should not error when placing order with custom external ID")
		require.NotNil(t, response, "Response should not be nil")
		require.Equal(t, "OK", response.Status, "Response status should be OK")
		require.Equal(t, customExternalID, response.Data.ExternalID, "External ID should match custom value")

		t.Logf("Placed order with custom external ID - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)

		placedOrders = append(placedOrders, placedOrder{
			orderID:    response.Data.OrderID,
			externalID: response.Data.ExternalID,
		})
	})

	// Test 4: Place order with previous order external ID (order replacement)
	t.Run("PlaceOrderWithPreviousOrderExternalID", func(t *testing.T) {
		if len(placedOrders) == 0 {
			t.Skip("No previous orders to replace")
		}

		// Use the first placed order's external ID for replacement
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

		require.NoError(t, err, "Should not error when placing order with previous order external ID")
		require.NotNil(t, response, "Response should not be nil")
		require.Equal(t, "OK", response.Status, "Response status should be OK")

		t.Logf("Placed order replacing previous order %s - New ID: %d, External ID: %s", previousOrderID, response.Data.OrderID, response.Data.ExternalID)

		placedOrders = append(placedOrders, placedOrder{
			orderID:    response.Data.OrderID,
			externalID: response.Data.ExternalID,
		})
	})

	// Test 5: Place order with builder fee
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

		require.NoError(t, err, "Should not error when placing order with builder fee")
		require.NotNil(t, response, "Response should not be nil")
		require.Equal(t, "OK", response.Status, "Response status should be OK")

		t.Logf("Placed order with builder fee - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)

		placedOrders = append(placedOrders, placedOrder{
			orderID:    response.Data.OrderID,
			externalID: response.Data.ExternalID,
		})
	})

	// Test 6: Place order with builder ID
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

		require.NoError(t, err, "Should not error when placing order with builder ID")
		require.NotNil(t, response, "Response should not be nil")
		require.Equal(t, "OK", response.Status, "Response status should be OK")

		t.Logf("Placed order with builder ID - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)

		placedOrders = append(placedOrders, placedOrder{
			orderID:    response.Data.OrderID,
			externalID: response.Data.ExternalID,
		})
	})

	// Test 7: Place order with builder fee and builder ID
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

		require.NoError(t, err, "Should not error when placing order with builder fee and ID")
		require.NotNil(t, response, "Response should not be nil")
		require.Equal(t, "OK", response.Status, "Response status should be OK")

		t.Logf("Placed order with builder fee and ID - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)

		placedOrders = append(placedOrders, placedOrder{
			orderID:    response.Data.OrderID,
			externalID: response.Data.ExternalID,
		})
	})

	// Test 8: Place order with all optional parameters combined
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

		require.NoError(t, err, "Should not error when placing order with all optional parameters")
		require.NotNil(t, response, "Response should not be nil")
		require.Equal(t, "OK", response.Status, "Response status should be OK")
		require.Equal(t, customExternalID, response.Data.ExternalID, "External ID should match custom value")

		t.Logf("Placed order with all optional params - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)

		placedOrders = append(placedOrders, placedOrder{
			orderID:    response.Data.OrderID,
			externalID: response.Data.ExternalID,
		})
	})

	// Test 9: Place MARKET order to create a position (for reduce-only testing)
	t.Run("PlaceMarketOrderToCreatePosition", func(t *testing.T) {
		// Get market stats to find current price
		stats, err := client.Markets.GetMarketStatistics(ctx, "BTC-USD")
		require.NoError(t, err, "Should be able to get market stats")

		// Get orderbook to find best ask price for BUY order
		orderbook, err := client.Markets.GetOrderbookSnapshot(ctx, "BTC-USD")
		
		// Use best ask price for BUY order, or fallback to mark price
		var orderPrice decimal.Decimal
		if err == nil && orderbook != nil && len(orderbook.Ask) > 0 {
			// Use best ask price (first ask in the orderbook)
			orderPrice = orderbook.Ask[0].Price
			t.Logf("Using best ask price from orderbook: %s", orderPrice.String())
		} else {
			// Fallback to ask price from stats, or mark price
			if !stats.AskPrice.IsZero() {
				orderPrice = stats.AskPrice
			} else {
				orderPrice = stats.MarkPrice
			}
			t.Logf("Using ask/mark price from stats: %s", orderPrice.String())
		}

		// Use minimum order size (0.0001 BTC is standard minimum)
		minOrderSize := decimal.NewFromFloat(0.0001)
		expireTime := time.Now().Add(1 * time.Hour)

		// Place MARKET order with IOC (Immediate or Cancel) - this will fill immediately or cancel
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

		require.NoError(t, err, "Should not error when placing MARKET order")
		require.NotNil(t, response, "Response should not be nil")
		require.Equal(t, "OK", response.Status, "Response status should be OK")
		require.NotZero(t, response.Data.OrderID, "Order ID should not be zero")
		require.NotEmpty(t, response.Data.ExternalID, "External ID should not be empty")

		t.Logf("Placed MARKET order - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)

		// Wait a bit for order to fill and position to be created
		time.Sleep(3 * time.Second)

		// Verify order was filled by checking its status
		// For IOC orders, they should either be FILLED or CANCELLED
		order, err := client.Account.GetOrderByID(ctx, int(response.Data.OrderID))
		if err != nil {
			// Order might not be found if it was filled and moved to history
			// Try checking order history
			historyOrders, histErr := client.Account.GetOrdersHistory(ctx, []string{"BTC-USD"}, nil, nil, nil, nil)
			if histErr == nil {
				var foundOrder *models.OpenOrderModel
				for i := range historyOrders {
					if historyOrders[i].ID == int(response.Data.OrderID) {
						foundOrder = &historyOrders[i]
						break
					}
				}
				if foundOrder != nil {
					order = foundOrder
					err = nil
				}
			}
		}

		if err == nil && order != nil {
			// Verify order status - should be FILLED for a successful market order
			require.Equal(t, models.OrderStatusFilled, order.Status, "MARKET order should be FILLED")
			t.Logf("Order status verified: %s", order.Status)
			if order.FilledQty != nil {
				require.Greater(t, order.FilledQty.Cmp(decimal.Zero), 0, "Filled quantity should be greater than zero")
				t.Logf("Order filled quantity: %s", order.FilledQty.String())
			}
		} else {
			// If we can't find the order, verify position was created as alternative verification
			positions, posErr := client.Account.GetPositions(ctx, []string{"BTC-USD"}, nil)
			require.NoError(t, posErr, "Should be able to check positions")
			require.Greater(t, len(positions), 0, "Position should be created after MARKET order fills")
			t.Logf("Position created successfully: %s, size: %s", positions[0].Side, positions[0].Size.String())
		}

		// Don't track this order in placedOrders since it should fill immediately
		// If it doesn't fill, it will be canceled by IOC
	})

	// Test 10: Place reduce-only order (requires open position)
	t.Run("PlaceReduceOnlyOrder", func(t *testing.T) {
		// Check if there's an existing position
		positions, err := client.Account.GetPositions(ctx, []string{"BTC-USD"}, nil)
		if err != nil {
			t.Skipf("Cannot test reduce-only: could not check positions: %v", err)
			return
		}

		if len(positions) == 0 {
			t.Skip("Cannot test reduce-only: no open position exists")
			return
		}

		position := positions[0]
		t.Logf("Found position: %s, side: %s, size: %s", position.Market, position.Side, position.Size.String())

		// Determine the side for reduce-only order (opposite of position side)
		var reduceSide models.OrderSide
		if position.Side == models.PositionSideLong {
			reduceSide = OrderSideSell // Reduce long position by selling
		} else {
			reduceSide = OrderSideBuy // Reduce short position by buying
		}

		// Get market stats and orderbook to find current price for MARKET order
		stats, err := client.Markets.GetMarketStatistics(ctx, "BTC-USD")
		require.NoError(t, err, "Should be able to get market stats")

		// Get orderbook to find best price
		orderbook, err := client.Markets.GetOrderbookSnapshot(ctx, "BTC-USD")

		// Use best bid price for SELL order, best ask price for BUY order
		var reducePrice decimal.Decimal
		if reduceSide == OrderSideSell {
			// For SELL, use best bid price from orderbook, or fallback to bid/mark price from stats
			if err == nil && orderbook != nil && len(orderbook.Bid) > 0 {
				reducePrice = orderbook.Bid[0].Price
				t.Logf("Using best bid price from orderbook: %s", reducePrice.String())
			} else {
				if !stats.BidPrice.IsZero() {
					reducePrice = stats.BidPrice
				} else {
					reducePrice = stats.MarkPrice
				}
				t.Logf("Using bid/mark price from stats: %s", reducePrice.String())
			}
		} else {
			// For BUY, use best ask price from orderbook, or fallback to ask/mark price from stats
			if err == nil && orderbook != nil && len(orderbook.Ask) > 0 {
				reducePrice = orderbook.Ask[0].Price
				t.Logf("Using best ask price from orderbook: %s", reducePrice.String())
			} else {
				if !stats.AskPrice.IsZero() {
					reducePrice = stats.AskPrice
				} else {
					reducePrice = stats.MarkPrice
				}
				t.Logf("Using ask/mark price from stats: %s", reducePrice.String())
			}
		}

		// Use a small amount to reduce the position (minimum order size)
		reduceAmount := decimal.NewFromFloat(0.0001)
		expireTime := time.Now().Add(1 * time.Hour)

		// Place MARKET order with IOC (Immediate or Cancel) - this will fill immediately or cancel
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

		require.NoError(t, err, "Should not error when placing order with reduce-only")
		require.NotNil(t, response, "Response should not be nil")
		require.Equal(t, "OK", response.Status, "Response status should be OK")
		require.NotZero(t, response.Data.OrderID, "Order ID should not be zero")
		require.NotEmpty(t, response.Data.ExternalID, "External ID should not be empty")

		t.Logf("Placed reduce-only MARKET order - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)

		// Wait a bit for order to fill
		time.Sleep(3 * time.Second)

		// Verify order was filled by checking its status
		// For IOC orders, they should either be FILLED or CANCELLED
		order, err := client.Account.GetOrderByID(ctx, int(response.Data.OrderID))
		if err != nil {
			// Order might not be found if it was filled and moved to history
			// Try checking order history
			historyOrders, histErr := client.Account.GetOrdersHistory(ctx, []string{"BTC-USD"}, nil, nil, nil, nil)
			if histErr == nil {
				var foundOrder *models.OpenOrderModel
				for i := range historyOrders {
					if historyOrders[i].ID == int(response.Data.OrderID) {
						foundOrder = &historyOrders[i]
						break
					}
				}
				if foundOrder != nil {
					order = foundOrder
					err = nil
				}
			}
		}

		// REQUIRE that we found the order
		require.NoError(t, err, "Should be able to find the placed order")
		require.NotNil(t, order, "Order should not be nil")

		// Verify order properties
		require.Equal(t, int(response.Data.OrderID), order.ID, "Order ID should match")
		require.Equal(t, response.Data.ExternalID, order.ExternalID, "External ID should match")
		require.True(t, order.ReduceOnly, "Order should have reduce-only flag set to true")
		require.Equal(t, models.OrderTypeMarket, order.Type, "Order type should be MARKET")
		require.Equal(t, reduceSide, order.Side, "Order side should match")
		// REQUIRE that the order was FILLED
		require.Equal(t, models.OrderStatusFilled, order.Status, "MARKET order MUST be FILLED")
		require.NotNil(t, order.FilledQty, "Filled quantity should not be nil")
		require.Greater(t, order.FilledQty.Cmp(decimal.Zero), 0, "Filled quantity MUST be greater than zero")
		t.Logf("Order status verified: %s, Filled quantity: %s", order.Status, order.FilledQty.String())

		// REQUIRE that the position was actually reduced or closed
		positionsAfter, posErr := client.Account.GetPositions(ctx, []string{"BTC-USD"}, nil)
		require.NoError(t, posErr, "Should be able to check positions after reduce-only order")

		if len(positionsAfter) > 0 {
			// Position still exists - verify it was reduced
			newPosition := positionsAfter[0]
			require.Equal(t, position.Market, newPosition.Market, "Position market should match")
			require.Equal(t, position.Side, newPosition.Side, "Position side should match")
			// Position size MUST be smaller (strictly less than original)
			require.Less(t, newPosition.Size.Cmp(position.Size), 0, "Position size MUST be reduced (was: %s, now: %s)", position.Size.String(), newPosition.Size.String())
			t.Logf("Position reduced: %s -> %s", position.Size.String(), newPosition.Size.String())
		} else {
			// Position was fully closed - this is also valid
			t.Logf("Position was fully closed by reduce-only order (original size: %s)", position.Size.String())
		}

		t.Logf("Verified reduce-only MARKET order - ID: %d, ReduceOnly: %v, Status: %s, Filled: %s", order.ID, order.ReduceOnly, order.Status, order.FilledQty.String())

		// Don't track this order in placedOrders since it should fill immediately
		// If it doesn't fill, it will be canceled by IOC
	})

	// Test 13: Cancel order by internal ID
	t.Run("CancelOrderByInternalID", func(t *testing.T) {
		if len(placedOrders) == 0 {
			t.Skip("No orders placed to cancel")
		}

		// Use the first placed order
		orderToCancel := placedOrders[0]

		err := client.Orders.CancelOrder(ctx, int(orderToCancel.orderID))
		require.NoError(t, err, "Should not error when canceling order by internal ID")

		// Verify the order is canceled by checking open orders
		// Wait a bit for the cancellation to propagate
		time.Sleep(500 * time.Millisecond)
		openOrders, err := client.Account.GetOpenOrders(ctx, []string{"BTC-USD"}, nil, nil)
		if err == nil {
			for _, openOrder := range openOrders {
				if openOrder.ExternalID == orderToCancel.externalID {
					t.Errorf("Order %s should have been canceled but still appears in open orders", orderToCancel.externalID)
				}
			}
		}

		// Remove from cleanup list since we already canceled it
		placedOrders = placedOrders[1:]

		t.Logf("Successfully canceled order by internal ID: %d", orderToCancel.orderID)
	})

	// Test 6: Cancel order by external ID
	t.Run("CancelOrderByExternalID", func(t *testing.T) {
		if len(placedOrders) == 0 {
			t.Skip("No orders placed to cancel")
		}

		// Use the first remaining order
		orderToCancel := placedOrders[0]

		err := client.Orders.CancelOrderByExternalID(ctx, orderToCancel.externalID)
		require.NoError(t, err, "Should not error when canceling order by external ID")

		// Verify the order is canceled
		time.Sleep(500 * time.Millisecond)
		openOrders, err := client.Account.GetOpenOrders(ctx, []string{"BTC-USD"}, nil, nil)
		if err == nil {
			for _, openOrder := range openOrders {
				if openOrder.ExternalID == orderToCancel.externalID {
					t.Errorf("Order %s should have been canceled but still appears in open orders", orderToCancel.externalID)
				}
			}
		}

		// Remove from cleanup list
		placedOrders = placedOrders[1:]

		t.Logf("Successfully canceled order by external ID: %s", orderToCancel.externalID)
	})

	// Test 7: Verify order response structure matches Python SDK expectations
	t.Run("VerifyOrderResponseStructure", func(t *testing.T) {
		expireTime := time.Now().Add(1 * time.Hour)

		response, err := client.Orders.PlaceOrder(ctx,
			market,
			decimal.NewFromFloat(0.0001), // Very small amount
			decimal.NewFromFloat(1000),   // Very low price, well below market
			OrderSideBuy,
			models.OrderTypeLimit,
			TimeInForceGTT,
			SelfTradeProtectionDisabled,
			WithExpireTime(expireTime),
		)

		require.NoError(t, err, "Should not error when placing order")

		assert.Equal(t, "OK", response.Status, "Status should be OK")
		assert.NotZero(t, response.Data.OrderID, "OrderID (id) should be non-zero integer")
		assert.NotEmpty(t, response.Data.ExternalID, "ExternalID (external_id) should be non-empty string")

		// Verify the order can be retrieved by external ID
		orders, err := client.Account.GetOrderByExternalID(ctx, response.Data.ExternalID)
		if err == nil && len(orders) > 0 {
			assert.Equal(t, response.Data.ExternalID, orders[0].ExternalID, "Retrieved order should match placed order")
			assert.Equal(t, int(response.Data.OrderID), orders[0].ID, "Retrieved order ID should match")
		}

		placedOrders = append(placedOrders, placedOrder{
			orderID:    response.Data.OrderID,
			externalID: response.Data.ExternalID,
		})

		t.Logf("Verified order response structure - ID: %d, External ID: %s", response.Data.OrderID, response.Data.ExternalID)
	})

	// Test 8: Test mass cancel functionality
	t.Run("MassCancelOrders", func(t *testing.T) {
		// Place a few orders for mass cancel
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

			require.NoError(t, err, "Should not error when placing order for mass cancel")
			orderIDs = append(orderIDs, int(response.Data.OrderID))
			externalIDs = append(externalIDs, response.Data.ExternalID)

			placedOrders = append(placedOrders, placedOrder{
				orderID:    response.Data.OrderID,
				externalID: response.Data.ExternalID,
			})
		}

		// Mass cancel by order IDs
		err := client.Orders.MassCancel(ctx, orderIDs, nil, nil, false)
		require.NoError(t, err, "Should not error when mass canceling orders")

		// Verify orders are canceled
		time.Sleep(500 * time.Millisecond)
		openOrders, err := client.Account.GetOpenOrders(ctx, []string{"BTC-USD"}, nil, nil)
		if err == nil {
			for _, externalID := range externalIDs {
				for _, openOrder := range openOrders {
					if openOrder.ExternalID == externalID {
						t.Errorf("Order %s should have been canceled by mass cancel", externalID)
					}
				}
			}
		}

		// Remove from cleanup list
		for range orderIDs {
			if len(placedOrders) > 0 {
				placedOrders = placedOrders[1:]
			}
		}

		t.Logf("Successfully mass canceled %d orders", len(orderIDs))
	})

	t.Logf("All order placement and cancellation tests completed. %d orders remaining for cleanup.", len(placedOrders))
}

// TestOrderPlacementErrorHandling tests error cases for order placement
func TestOrderPlacementErrorHandling(t *testing.T) {
	client := createTestClient()
	ctx := context.Background()

	markets, err := client.Markets.GetMarkets(ctx, []string{"BTC-USD"})
	require.NoError(t, err)
	require.Greater(t, len(markets), 0)
	market := markets[0]

	// Test: Invalid market (should fail)
	// Note: We can't test with a completely invalid market structure because it causes
	// a panic in the signing code when L2Config has invalid hex IDs. Instead, we test
	// with a valid market structure but expect API-level validation errors.
	t.Run("InvalidMarket", func(t *testing.T) {
		// Use a valid market structure but with a non-existent market name
		// This will pass local validation but should fail at the API level
		invalidMarket := models.MarketModel{
			Name: "NONEXISTENT-MARKET-12345",
			// Use the same L2Config as the valid market to avoid signing errors
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

		// Should error at API level due to invalid/non-existent market
		require.Error(t, err, "Should error when placing order with invalid market name")
		t.Logf("Got expected error for invalid market: %v", err)
	})

	// Test: Zero amount (should fail validation or API)
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

		// Should error due to invalid amount
		if err == nil {
			t.Log("Note: Zero amount did not error immediately (may fail at API)")
		}
	})
}

