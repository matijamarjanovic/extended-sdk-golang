package sdk

import (
	"time"

	"github.com/extended-protocol/extended-sdk-golang/src/client"
	"github.com/extended-protocol/extended-sdk-golang/src/models"
	"github.com/extended-protocol/extended-sdk-golang/src/services"
)

// Client provides REST API functionality for perpetual trading.
// It embeds BaseClient to reuse common functionality like HTTP client, auth, etc.
// It provides access to domain-specific services through the Account, Orders, and Markets fields.
type Client struct {
	*client.BaseClient
	Account *services.AccountService
	Orders  *services.OrdersService
	Markets *services.MarketsService
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
	sdkClient.Account = &services.AccountService{Base: baseClient}
	sdkClient.Orders = &services.OrdersService{Base: baseClient}
	sdkClient.Markets = &services.MarketsService{Base: baseClient}

	return sdkClient
}

// Close closes the HTTP client and cleans up resources.
// It delegates to BaseClient's Close method.
func (c *Client) Close() {
	c.BaseClient.Close()
}

