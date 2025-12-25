package orders

import (
	"github.com/extended-protocol/extended-sdk-golang/src"
	"github.com/extended-protocol/extended-sdk-golang/src/models"
)

// Service provides order-related API operations.
// It holds a reference to the main client to access shared infrastructure.
type Service struct {
	client *sdk.Client // Reference to main client
}

// Methods to be implemented:
// - SubmitOrder (move from api_client.go)
// - PlaceOrder (alias/wrapper for SubmitOrder)
// - CancelOrder (new)
// - CancelOrderByExternalID (new)
// - MassCancel (new)
//
// Split into multiple files (e.g., orders_cancel.go) as code grows



