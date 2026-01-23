package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/matijamarjanovic/extended-sdk-golang/x10/client"
	"github.com/matijamarjanovic/extended-sdk-golang/x10/models"
	"github.com/gorilla/websocket"
)

type StreamingService struct {
	Base *client.BaseClient
}

// StreamConnection represents an active WebSocket connection to a stream.
// It provides methods to receive messages and manage the connection lifecycle.
type StreamConnection struct {
	conn      *websocket.Conn
	closed    bool
	msgsCount int64
}

func (sc *StreamConnection) Close() error {
	if sc.closed {
		return nil
	}
	sc.closed = true
	return sc.conn.Close()
}

func (sc *StreamConnection) IsClosed() bool {
	return sc.closed
}

func (sc *StreamConnection) MessagesCount() int64 {
	return sc.msgsCount
}

// Recv receives a message from the stream and unmarshals it into the provided model.
// The model should match the expected type for the stream (e.g., OrderbookUpdateModel, StreamPublicTradeModel, etc.).
func (sc *StreamConnection) Recv(ctx context.Context, result interface{}) error {
	if sc.closed {
		return fmt.Errorf("connection is closed")
	}

	_, message, err := sc.conn.ReadMessage()
	if err != nil {
		sc.closed = true
		return fmt.Errorf("failed to read message: %w", err)
	}

	sc.msgsCount++

	var wrapped models.WrappedStreamResponse
	if err := json.Unmarshal(message, &wrapped); err != nil {
		return fmt.Errorf("failed to unmarshal wrapped response: %w", err)
	}

	if wrapped.Error != nil {
		return fmt.Errorf("stream error: %s", *wrapped.Error)
	}

	if wrapped.Data != nil {
		dataBytes, err := json.Marshal(wrapped.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal data: %w", err)
		}
		if err := json.Unmarshal(dataBytes, result); err != nil {
			return fmt.Errorf("failed to unmarshal data: %w", err)
		}
	}

	return nil
}

// SubscribeToOrderbooks subscribes to orderbook updates.
// If marketName is empty, subscribes to all markets.
// If depth is nil, uses default depth.
func (s *StreamingService) SubscribeToOrderbooks(ctx context.Context, marketName string, depth *int) (*StreamConnection, error) {
	path := "/orderbooks"
	if marketName != "" {
		path = "/orderbooks/" + marketName
	}

	query := make(map[string]string)
	if depth != nil {
		query["depth"] = strconv.Itoa(*depth)
	}

	streamURL, err := s.buildStreamURL(path, query)
	if err != nil {
		return nil, err
	}

	return s.connect(ctx, streamURL)
}

// SubscribeToPublicTrades subscribes to public trade updates.
// If marketName is empty, subscribes to all markets.
func (s *StreamingService) SubscribeToPublicTrades(ctx context.Context, marketName string) (*StreamConnection, error) {
	path := "/publicTrades"
	if marketName != "" {
		path = "/publicTrades/" + marketName
	}

	streamURL, err := s.buildStreamURL(path, nil)
	if err != nil {
		return nil, err
	}

	return s.connect(ctx, streamURL)
}

// SubscribeToFundingRates subscribes to funding rate updates.
// If marketName is empty, subscribes to all markets.
func (s *StreamingService) SubscribeToFundingRates(ctx context.Context, marketName string) (*StreamConnection, error) {
	path := "/funding"
	if marketName != "" {
		path = "/funding/" + marketName
	}

	streamURL, err := s.buildStreamURL(path, nil)
	if err != nil {
		return nil, err
	}

	return s.connect(ctx, streamURL)
}

// SubscribeToCandles subscribes to candle updates for a specific market, candle type, and interval.
func (s *StreamingService) SubscribeToCandles(ctx context.Context, marketName string, candleType models.CandleType, interval models.CandleInterval) (*StreamConnection, error) {
	path := fmt.Sprintf("/candles/%s/%s", marketName, candleType)

	query := map[string]string{
		"interval": string(interval),
	}

	streamURL, err := s.buildStreamURL(path, query)
	if err != nil {
		return nil, err
	}

	return s.connect(ctx, streamURL)
}

// SubscribeToAccountUpdates subscribes to account updates (orders, positions, trades, balance).
// Requires an API key to be set in the BaseClient.
func (s *StreamingService) SubscribeToAccountUpdates(ctx context.Context) (*StreamConnection, error) {
	streamURL, err := s.buildStreamURL("/account", nil)
	if err != nil {
		return nil, err
	}

	return s.connect(ctx, streamURL)
}

// SubscribeToMarkPrices subscribes to mark price updates.
// Mark prices are used to calculate unrealized P&L and serve as the reference for liquidations.
// If marketName is empty, subscribes to all markets.
func (s *StreamingService) SubscribeToMarkPrices(ctx context.Context, marketName string) (*StreamConnection, error) {
	path := "/prices/mark"
	if marketName != "" {
		path = "/prices/mark/" + marketName
	}

	streamURL, err := s.buildStreamURL(path, nil)
	if err != nil {
		return nil, err
	}

	return s.connect(ctx, streamURL)
}

// SubscribeToIndexPrices subscribes to index price updates.
// Index prices are composite spot prices sourced from multiple external providers,
// used as the reference for funding-rate calculations.
// If marketName is empty, subscribes to all markets.
func (s *StreamingService) SubscribeToIndexPrices(ctx context.Context, marketName string) (*StreamConnection, error) {
	path := "/prices/index"
	if marketName != "" {
		path = "/prices/index/" + marketName
	}

	streamURL, err := s.buildStreamURL(path, nil)
	if err != nil {
		return nil, err
	}

	return s.connect(ctx, streamURL)
}

// buildStreamURL builds a WebSocket URL for the given path and parameters.
func (s *StreamingService) buildStreamURL(path string, query map[string]string) (string, error) {
	streamBaseURL := s.Base.EndpointConfig().StreamURL

	streamBaseURL = strings.TrimSuffix(streamBaseURL, "/")

	path = strings.TrimPrefix(path, "/")

	fullURL := streamBaseURL + "/" + path

	if len(query) > 0 {
		u, err := url.Parse(fullURL)
		if err != nil {
			return "", fmt.Errorf("failed to parse URL: %w", err)
		}
		q := u.Query()
		for k, v := range query {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}

	return fullURL, nil
}

// connect establishes a WebSocket connection to the given URL.
func (s *StreamingService) connect(ctx context.Context, streamURL string) (*StreamConnection, error) {
	apiKey, _ := s.Base.APIKey()

	dialer := websocket.Dialer{}

	headers := make(map[string][]string)
	headers["User-Agent"] = []string{"ExtendedSDKGolang/0.1.0"}
	if apiKey != "" {
		headers["X-Api-Key"] = []string{apiKey}
	}

	conn, _, err := dialer.DialContext(ctx, streamURL, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to stream: %w", err)
	}

	return &StreamConnection{
		conn:      conn,
		closed:    false,
		msgsCount: 0,
	}, nil
}
