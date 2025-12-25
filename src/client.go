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
// It embeds BaseClient to reuse common functionality like HTTP client, auth, etc.
// It provides access to domain-specific services through the Account, Orders, and Markets fields.
type Client struct {
	*client.BaseClient
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
	baseClient := client.NewBaseClient(cfg, starkAccount.APIKey(), starkAccount, nil, clientTimeout)
	sdkClient := &Client{
		BaseClient: baseClient,
	}

	// Initialize services with reference to BaseClient
	sdkClient.Account = &account.Service{Base: baseClient}
	sdkClient.Orders = &orders.Service{Base: baseClient}
	sdkClient.Markets = &markets.Service{Base: baseClient}

	return sdkClient
}

// Close closes the HTTP client and cleans up resources.
// It delegates to BaseClient's Close method.
func (c *Client) Close() {
	c.BaseClient.Close()
}

