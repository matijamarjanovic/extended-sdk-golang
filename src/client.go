package sdk

import (
	"time"

	"github.com/extended-protocol/extended-sdk-golang/src/account"
	"github.com/extended-protocol/extended-sdk-golang/src/client"
	"github.com/extended-protocol/extended-sdk-golang/src/markets"
	"github.com/extended-protocol/extended-sdk-golang/src/models"
	"github.com/extended-protocol/extended-sdk-golang/src/orders"
)

// Client provides REST API functionality for perpetual trading.
// It embeds BaseModule to reuse common functionality like HTTP client, auth, etc.
// It provides access to domain-specific services through the Account, Orders, and Markets fields.
type Client struct {
	*client.BaseModule
	Account *account.Service
	Orders  *orders.Service
	Markets *markets.Service
}

// NewClient creates a new Client instance with all services initialized.
// It takes an endpoint configuration, a Stark perpetual account, and a client timeout.
func NewClient(
	cfg models.EndpointConfig,
	starkAccount *client.StarkPerpetualAccount,
	clientTimeout time.Duration,
) *Client {
	baseModule := client.NewBaseModule(cfg, starkAccount.APIKey(), starkAccount, nil, clientTimeout)
	sdkClient := &Client{
		BaseModule: baseModule,
	}

	// Initialize services with reference to BaseModule
	sdkClient.Account = &account.Service{Base: baseModule}
	sdkClient.Orders = &orders.Service{Base: baseModule}
	sdkClient.Markets = &markets.Service{Base: baseModule}

	return sdkClient
}

// Close closes the HTTP client and cleans up resources.
// It delegates to BaseModule's Close method.
func (c *Client) Close() {
	c.BaseModule.Close()
}

