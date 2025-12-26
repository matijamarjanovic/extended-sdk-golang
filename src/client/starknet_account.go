package client

import (
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/extended-protocol/extended-sdk-golang/src/models"
)

// StarkPerpetualAccount represents a Stark perpetual trading account.
type StarkPerpetualAccount struct {
	vault      uint64
	privateKey string
	publicKey  string
	apiKey     string
	
	// tradingFee caches market-specific trading fees retrieved from the API via GET /api/v1/user/fees?market={market}
	// The map key is the market name (e.g., "BTC-USD")
	// Fees are determined by the platform for each sub-account and cannot be set by users.
	// This cache should be populated by calling AccountService.GetMarketFee() or AccountService.GetFees()
	// and then updating the cache using SetTradingFee() or SetTradingFees().
	// If a market is not found in this cache, DefaultFees will be used as fallback.
	tradingFee map[string]models.TradingFeeModel
	tradingFeeMu sync.RWMutex // Protects tradingFee map
}

// NewStarkPerpetualAccount constructs the account, validating hex inputs.
func NewStarkPerpetualAccount(vault uint64, privateKeyHex, publicKeyHex, apiKey string) (*StarkPerpetualAccount, error) {
	if err := isHexString(privateKeyHex); err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}
	if err := isHexString(publicKeyHex); err != nil {
		return nil, fmt.Errorf("invalid public key: %w", err)
	}

	// Ensure that private key and public key have 0x prefix
	if len(privateKeyHex) < 2 || privateKeyHex[:2] != "0x" {
		return nil, fmt.Errorf("private key must start with 0x")
	}
	if len(publicKeyHex) < 2 || publicKeyHex[:2] != "0x" {
		return nil, fmt.Errorf("public key must start with 0x")
	}

	// Check that API key does not start with 0x
	if len(apiKey) >= 2 && apiKey[:2] == "0x" {
		return nil, fmt.Errorf("api key should not start with 0x")
	}

	return &StarkPerpetualAccount{
		vault:      vault,
		privateKey: privateKeyHex,
		publicKey:  publicKeyHex,
		apiKey:     apiKey,
		tradingFee: make(map[string]models.TradingFeeModel),
	}, nil
}

// Vault returns the vault id.
func (s *StarkPerpetualAccount) Vault() uint64 { return s.vault }

// PublicKey returns the public key as a string.
func (s *StarkPerpetualAccount) PublicKey() string { return s.publicKey }

// APIKey returns the API key string.
func (s *StarkPerpetualAccount) APIKey() string { return s.apiKey }

// Sign delegates to SignFunc, returning (r,s).
func (stark *StarkPerpetualAccount) Sign(msgHash string) (*big.Int, *big.Int, error) {
	if msgHash == "" {
		return big.NewInt(0), big.NewInt(0), errors.New("msgHash is empty")
	}

	sig, err := SignMessage(msgHash, stark.privateKey)
	if err != nil {
		return big.NewInt(0), big.NewInt(0), err
	}

	// Extract r, s from the signature string.
	// Signature is in the format of {r}{s}{v}, where r, s and v are 64 chars each (192 hex chars).
	r, isGoodR := big.NewInt(0).SetString(sig[:64], 16)
	s, isGoodS := big.NewInt(0).SetString(sig[64:128], 16)

	if !isGoodR || !isGoodS {
		return big.NewInt(0), big.NewInt(0), errors.New("big int setting failed")
	}

	return r, s, nil
}

// GetTradingFee returns the trading fee for the given market.
// If the market is not found in the account's trading_fee cache, it returns DefaultFees as fallback.
// Fees are determined by the platform via GET /api/v1/user/fees?market={market} and cannot be set by users.
// To populate the cache, call AccountService.GetMarketFee() or AccountService.GetFees()
// and then SetTradingFee() to cache the API response.
func (s *StarkPerpetualAccount) GetTradingFee(marketName string) models.TradingFeeModel {
	s.tradingFeeMu.RLock()
	defer s.tradingFeeMu.RUnlock()
	
	if fee, ok := s.tradingFee[marketName]; ok {
		return fee
	}
	return models.DefaultFees
}

// SetTradingFee caches the trading fee for a specific market from an API response.
// This should be called after retrieving fees via AccountService.GetMarketFee() or AccountService.GetFees()
// to cache the platform-determined fees. Users cannot set arbitrary fees - fees are determined by the platform.
func (s *StarkPerpetualAccount) SetTradingFee(marketName string, fee models.TradingFeeModel) {
	s.tradingFeeMu.Lock()
	defer s.tradingFeeMu.Unlock()
	
	s.tradingFee[marketName] = fee
}

// SetTradingFees caches multiple trading fees at once from an API response.
// This is useful when retrieving fees for multiple markets via AccountService.GetFees()
// to cache the platform-determined fees. Users cannot set arbitrary fees - fees are determined by the platform.
func (s *StarkPerpetualAccount) SetTradingFees(fees map[string]models.TradingFeeModel) {
	s.tradingFeeMu.Lock()
	defer s.tradingFeeMu.Unlock()
	
	for market, fee := range fees {
		s.tradingFee[market] = fee
	}
}

