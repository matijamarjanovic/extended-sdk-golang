package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/extended-protocol/extended-sdk-golang/src/client"
	"github.com/extended-protocol/extended-sdk-golang/src/models"
)

// OrdersService provides order-related API operations.
type OrdersService struct {
	Base *client.BaseClient
}

// SubmitOrder submits a perpetual order to the trading API
func (s *OrdersService) SubmitOrder(ctx context.Context, order *models.PerpetualOrderModel) (*models.OrderResponse, error) {
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

// Methods to be implemented:
// - PlaceOrder (alias/wrapper for SubmitOrder)
// - CancelOrder (new)
// - CancelOrderByExternalID (new)
// - MassCancel (new)
//
// Split into multiple files (e.g., orders_cancel.go) as code grows

