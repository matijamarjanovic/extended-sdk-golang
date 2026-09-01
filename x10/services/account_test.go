package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/matijamarjanovic/extended-sdk-golang/x10/client"
	"github.com/matijamarjanovic/extended-sdk-golang/x10/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// GET /user/orders/{id} returns a single order object in `data`, not a list
// (reference: Python SDK account_module.get_order_by_id -> WrappedApiResponseModel[OpenOrderModel]).
func TestAccountService_GetOrderByID_SingleObjectData(t *testing.T) {
	const orderID = 275854637944209408
	const responseJSON = `{
		"status": "OK",
		"data": {
			"id": 275854637944209408,
			"accountId": 10002,
			"externalId": "0xabcdef",
			"market": "BTC-USD",
			"type": "LIMIT",
			"side": "BUY",
			"status": "NEW",
			"price": "43445.1168",
			"qty": "0.001",
			"reduceOnly": false,
			"postOnly": false,
			"createdTime": 1704420537000,
			"updatedTime": 1704420537000,
			"timeInForce": "GTT"
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/user/orders/275854637944209408", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(responseJSON))
	}))
	defer server.Close()

	cfg := models.EndpointConfig{APIBaseURL: server.URL}
	baseClient := client.NewBaseClient(cfg, TestAPIKey, nil, nil, 10*time.Second)
	service := &AccountService{Base: baseClient}

	order, err := service.GetOrderByID(context.Background(), orderID)
	require.NoError(t, err)
	require.NotNil(t, order)
	assert.Equal(t, orderID, order.ID)
	assert.Equal(t, "BTC-USD", order.Market)
	assert.Equal(t, models.OrderSideBuy, order.Side)
	assert.True(t, order.Price.Equal(decimal.RequireFromString("43445.1168")), "price mismatch: %s", order.Price)
	assert.True(t, order.Qty.Equal(decimal.RequireFromString("0.001")), "qty mismatch: %s", order.Qty)
}
