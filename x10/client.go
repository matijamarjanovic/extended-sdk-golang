package sdk

import (
	"time"

	"github.com/extended-protocol/extended-sdk-golang/x10/client"
	"github.com/extended-protocol/extended-sdk-golang/x10/models"
	"github.com/extended-protocol/extended-sdk-golang/x10/services"
)

// Client provides REST API functionality for perpetual trading.
// It embeds BaseClient to reuse common functionality like HTTP client, auth, etc.
// It provides access to domain-specific services through the Account, Orders, Markets, and Streaming fields.
type Client struct {
	*client.BaseClient
	Account   *services.AccountService
	Orders    *services.OrdersService
	Markets   *services.MarketsService
	Streaming *services.StreamingService
}

// NewClient creates a new Client instance with all services initialized.
// It takes an endpoint configuration, a Stark perpetual account, and a client timeout.
func NewClient(
	cfg models.EndpointConfig,
	starkAccount *StarkPerpetualAccount,
	clientTimeout time.Duration,
) *Client {
	acc := (*client.StarkPerpetualAccount)(starkAccount)
	baseClient := client.NewBaseClient(cfg, acc.APIKey(), acc, nil, clientTimeout)
	sdkClient := &Client{
		BaseClient: baseClient,
	}

	sdkClient.Account = &services.AccountService{Base: baseClient}
	sdkClient.Orders = &services.OrdersService{Base: baseClient}
	sdkClient.Markets = &services.MarketsService{Base: baseClient}
	sdkClient.Streaming = &services.StreamingService{Base: baseClient}

	return sdkClient
}

// Close closes the HTTP client and cleans up resources.
// It delegates to BaseClient's Close method.
func (c *Client) Close() {
	c.BaseClient.Close()
}
