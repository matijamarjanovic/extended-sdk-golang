package sdk

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/extended-protocol/extended-sdk-golang/src/client"
	"github.com/extended-protocol/extended-sdk-golang/src/models"
	"github.com/joho/godotenv"
)

// Test setup functions

func init() { load() }

func load() {
	wd, _ := os.Getwd()
	for {
		p := filepath.Join(wd, ".env")
		if _, err := os.Stat(p); err == nil {
			_ = godotenv.Load(p)
			return
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			return
		}
		wd = parent
	}
}

func createTestClient() *Client {
	apiKey := os.Getenv("TEST_API_KEY")
	vaultStr := os.Getenv("TEST_VAULT")
	vault, _ := strconv.ParseUint(vaultStr, 10, 64)
	publicKey := os.Getenv("TEST_PUBLIC_KEY")
	privateKey := os.Getenv("TEST_PRIVATE_KEY")

	account, err := NewStarkPerpetualAccount(vault, privateKey, publicKey, apiKey)

	if err != nil {
		panic("Failed to create StarkPerpetualAccount: " + err.Error())
	}

	return NewClient(STARKNET_MAINNET_CONFIG, account, 30*time.Second)
}

// Type aliases for commonly used model types in tests
// These allow using shorter names in test files without the models. prefix

type MarketModel = models.MarketModel
type L2ConfigModel = models.L2ConfigModel
type StarknetDomain = models.StarknetDomain

// Type aliases for client package types
type StarkPerpetualAccount = client.StarkPerpetualAccount

// Function aliases for client package functions
var NewStarkPerpetualAccount = client.NewStarkPerpetualAccount
var GetOrderHash = client.GetOrderHash
var SignMessage = client.SignMessage

// Constants for commonly used enum values in tests
const (
	OrderSideBuy  = models.OrderSideBuy
	OrderSideSell = models.OrderSideSell

	TimeInForceGTT = models.TimeInForceGTT

	SelfTradeProtectionDisabled = models.SelfTradeProtectionDisabled
	SelfTradeProtectionAccount  = models.SelfTradeProtectionAccount
	SelfTradeProtectionClient   = models.SelfTradeProtectionClient
)

