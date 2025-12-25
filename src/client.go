package sdk

import (
	"time"

	"github.com/extended-protocol/extended-sdk-golang/src/account"
	"github.com/extended-protocol/extended-sdk-golang/src/markets"
	"github.com/extended-protocol/extended-sdk-golang/src/models"
	"github.com/extended-protocol/extended-sdk-golang/src/orders"
)

// Client provides REST API functionality for perpetual trading.
// It embeds BaseModule to reuse common functionality like HTTP client, auth, etc.
// It provides access to domain-specific services through the Account, Orders, and Markets fields.
type Client struct {
	*BaseModule
	Account *account.Service
	Orders  *orders.Service
	Markets *markets.Service
}

// NewClient creates a new Client instance with all services initialized.
// It takes an endpoint configuration, a Stark perpetual account, and a client timeout.
func NewClient(
	cfg models.EndpointConfig,
	starkAccount *StarkPerpetualAccount,
	clientTimeout time.Duration,
) *Client {
	baseModule := NewBaseModule(cfg, starkAccount.APIKey(), starkAccount, nil, clientTimeout)
	client := &Client{
		BaseModule: baseModule,
	}

	// Initialize services with reference to the client
	client.Account = &account.Service{Client: client}
	client.Orders = &orders.Service{Client: client}
	client.Markets = &markets.Service{Client: client}

	return client
}

// Close closes the HTTP client and cleans up resources.
// It delegates to BaseModule's Close method.
func (c *Client) Close() {
	c.BaseModule.Close()
}

