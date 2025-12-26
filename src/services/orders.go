package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/extended-protocol/extended-sdk-golang/src/client"
	"github.com/extended-protocol/extended-sdk-golang/src/models"
	"github.com/shopspring/decimal"
)

// OrdersService provides order-related API operations.
type OrdersService struct {
	Base *client.BaseClient
}

// PlaceOrder creates an order object and submits it to the exchange.
// Required parameters are passed as function arguments, optional parameters are passed as options.
// It uses the account from the service's BaseClient and always uses account.Sign as the signer.
// If nonce is not provided via WithNonce option, it will be auto-generated.
//
// Example usage:
//
//	order, err := client.Orders.PlaceOrder(ctx,
//		market, amount, price, models.OrderSideBuy, models.OrderTypeLimit,
//		models.TimeInForceGTT, models.SelfTradeProtectionAccount,
//		sdk.WithPostOnly(true),
//		sdk.WithReduceOnly(false),
//		sdk.WithBuilderFee(builderFee),
//		sdk.WithNonce(customNonce), // Optional: omit to auto-generate
//	)
func (s *OrdersService) PlaceOrder(
	ctx context.Context,
	market models.MarketModel,
	syntheticAmount decimal.Decimal,
	price decimal.Decimal,
	side models.OrderSide,
	orderType models.OrderType,
	timeInForce models.TimeInForce,
	selfTradeProtectionLevel models.SelfTradeProtectionLevel,
	opts ...PlaceOrderOption,
) (*models.OrderResponse, error) {
	// Build config from options
	config := buildPlaceOrderConfig(
		market,
		syntheticAmount,
		price,
		side,
		orderType,
		timeInForce,
		selfTradeProtectionLevel,
		opts...,
	)

	// Get account and config from BaseClient
	account, err := s.Base.StarkAccount()
	if err != nil {
		return nil, fmt.Errorf("failed to get stark account: %w", err)
	}

	endpointConfig := s.Base.EndpointConfig()

	// Set default expire time if not provided (1 hour from now)
	expireTime := config.ExpireTime
	if expireTime == nil {
		defaultExpire := time.Now().Add(1 * time.Hour)
		expireTime = &defaultExpire
	}

	// Create order object
	order, err := createOrderObject(createOrderObjectParams{
		Market:                   config.Market,
		Account:                  account,
		SyntheticAmount:          config.SyntheticAmount,
		Price:                    config.Price,
		Side:                     config.Side,
		Type:                     config.Type,
		StarknetDomain:           endpointConfig.StarknetDomain,
		ExpireTime:               *expireTime,
		PostOnly:                 config.PostOnly,
		ReduceOnly:               config.ReduceOnly,
		PreviousOrderExternalID:  config.PreviousOrderExternalID,
		OrderExternalID:          config.OrderExternalID,
		TimeInForce:              config.TimeInForce,
		SelfTradeProtectionLevel: config.SelfTradeProtectionLevel,
		Nonce:                    config.Nonce,
		BuilderFee:               config.BuilderFee,
		BuilderID:                config.BuilderID,
		TpSlType:                 config.TpSlType,
		TakeProfit:               config.TakeProfit,
		StopLoss:                 config.StopLoss,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create order object: %w", err)
	}

	// Submit the order
	return s.submitOrder(ctx, order)
}

// CancelOrder cancels an order by its internal order ID.
// https://api.docs.extended.exchange/#cancel-order
func (s *OrdersService) CancelOrder(ctx context.Context, orderID int) error {
	baseUrl, err := s.Base.GetURL(fmt.Sprintf("/user/order/%d", orderID), nil)
	if err != nil {
		return fmt.Errorf("failed to build URL: %w", err)
	}

	var emptyResponse models.EmptyResponse
	if err := s.Base.DoRequest(ctx, "DELETE", baseUrl, nil, &emptyResponse); err != nil {
		return err
	}

	if emptyResponse.Status != "OK" {
		return fmt.Errorf("API returned error status: %v", emptyResponse.Status)
	}

	return nil
}

// CancelOrderByExternalID cancels an order by its external ID.
// https://api.docs.extended.exchange/#cancel-order
func (s *OrdersService) CancelOrderByExternalID(ctx context.Context, externalID string) error {
	query := map[string]string{
		"externalId": externalID,
	}
	baseUrl, err := s.Base.GetURL("/user/order", query)
	if err != nil {
		return fmt.Errorf("failed to build URL: %w", err)
	}

	var emptyResponse models.EmptyResponse
	if err := s.Base.DoRequest(ctx, "DELETE", baseUrl, nil, &emptyResponse); err != nil {
		return err
	}

	if emptyResponse.Status != "OK" {
		return fmt.Errorf("API returned error status: %v", emptyResponse.Status)
	}

	return nil
}

// MassCancelRequest represents the request parameters for mass cancel operation
type MassCancelRequest struct {
	OrderIDs         []int    `json:"orderIds,omitempty"`
	ExternalOrderIDs []string `json:"externalOrderIds,omitempty"`
	Markets          []string `json:"markets,omitempty"`
	CancelAll        *bool    `json:"cancelAll,omitempty"`
}

// MassCancel cancels multiple orders based on the provided criteria.
// https://api.docs.extended.exchange/#mass-cancel
func (s *OrdersService) MassCancel(
	ctx context.Context,
	orderIDs []int,
	externalOrderIDs []string,
	markets []string,
	cancelAll bool,
) error {
	baseUrl, err := s.Base.GetURL("/user/order/massCancel", nil)
	if err != nil {
		return fmt.Errorf("failed to build URL: %w", err)
	}

	request := MassCancelRequest{
		OrderIDs:         orderIDs,
		ExternalOrderIDs: externalOrderIDs,
		Markets:          markets,
	}
	if cancelAll {
		request.CancelAll = &cancelAll
	}

	// Marshal the request to JSON
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request to JSON: %w", err)
	}

	// Create a buffer with the JSON data
	jsonData := bytes.NewBuffer(requestJSON)

	var emptyResponse models.EmptyResponse
	if err := s.Base.DoRequest(ctx, "POST", baseUrl, jsonData, &emptyResponse); err != nil {
		return err
	}

	if emptyResponse.Status != "OK" {
		return fmt.Errorf("API returned error status: %v", emptyResponse.Status)
	}

	return nil
}

// submitOrder submits a perpetual order to the trading API
func (s *OrdersService) submitOrder(ctx context.Context, order *models.PerpetualOrderModel) (*models.OrderResponse, error) {
	// Validate order object is complete and properly signed
	if order == nil {
		return nil, fmt.Errorf("order is nil")
	}

	baseUrl, err := s.Base.GetURL("/user/order", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	// Marshal the order to JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order to JSON: %w", err)
	}

	// Create a buffer with the JSON data
	jsonData := bytes.NewBuffer(orderJSON)

	// Use the DoRequest method to handle the HTTP request and JSON parsing
	var orderResponse models.OrderResponse
	if err := s.Base.DoRequest(ctx, "POST", baseUrl, jsonData, &orderResponse); err != nil {
		return nil, err
	}

	if orderResponse.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %v", orderResponse.Status)
	}

	if orderResponse.Data.ExternalID != order.ID {
		return nil, fmt.Errorf("mismatched order ID in response: got %s, expected %s", orderResponse.Data.ExternalID, order.ID)
	}

	return &orderResponse, nil
}
