package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/extended-protocol/extended-sdk-golang/src/models"
)

var (
	ErrAPIKeyNotSet       = errors.New("api key is not set")
	ErrStarkAccountNotSet = errors.New("stark account is not set")
)

// BaseModule provides common functionality for API modules.
type BaseModule struct {
	endpointConfig models.EndpointConfig
	apiKey         string
	starkAccount   *StarkPerpetualAccount
	httpClient     *http.Client
	clientTimeout  time.Duration
}

// NewBaseModule constructs a BaseModule with all fields explicitly provided.
// Pass nil for httpClient to allow lazy creation. Pass nil for starkAccount if intentionally absent.
func NewBaseModule(
	cfg models.EndpointConfig,
	apiKey string,
	starkAccount *StarkPerpetualAccount,
	httpClient *http.Client,
	clientTimeout time.Duration,
) *BaseModule {
	return &BaseModule{
		endpointConfig: cfg,
		apiKey:         apiKey,
		starkAccount:   starkAccount,
		httpClient:     httpClient,
		clientTimeout:  clientTimeout,
	}
}

func (m *BaseModule) EndpointConfig() models.EndpointConfig {
	return m.endpointConfig
}

func (m *BaseModule) APIKey() (string, error) {
	if m.apiKey == "" {
		return "", ErrAPIKeyNotSet
	}
	return m.apiKey, nil
}

func (m *BaseModule) StarkAccount() (*StarkPerpetualAccount, error) {
	if m.starkAccount == nil {
		return nil, ErrStarkAccountNotSet
	}
	return m.starkAccount, nil
}

func (m *BaseModule) HTTPClient() *http.Client {
	if m.httpClient == nil {
		m.httpClient = &http.Client{
			Timeout: m.clientTimeout,
		}
	}
	return m.httpClient
}

// Close analogous to closing aiohttp session.
func (m *BaseModule) Close() {
	if m.httpClient != nil {
		m.httpClient.CloseIdleConnections()
		m.httpClient = nil
	}
}

// GetURL builds a full URL with optional query params.
func (m *BaseModule) GetURL(path string, query map[string]string) (string, error) {
	full := m.endpointConfig.APIBaseURL + path
	u, err := url.Parse(full)
	if err != nil {
		return "", err
	}
	if len(query) > 0 {
		q := u.Query()
		for k, v := range query {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
	}
	return u.String(), nil
}

// DoRequest performs an HTTP request and unmarshals the JSON response into the provided object
// This function deduplicates common HTTP request logic across the SDK
func (m *BaseModule) DoRequest(ctx context.Context, method, url string, body io.Reader, result interface{}) error {
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "ExtendedSDKGolang/0.1.0")
	req.Header.Set("Content-Type", "application/json")
	
	if m.apiKey != "" {
		req.Header.Set("X-Api-Key", m.apiKey)
	}

	// Execute request
	client := m.HTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	// Parse JSON response into the provided result object
	if err := json.Unmarshal(responseBody, result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}