package x10

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/extended-protocol/extended-sdk-golang/x10/client"
)

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
		panic(err)
	}

	return NewClient(STARKNET_MAINNET_CONFIG, account, 30*time.Second)
}

var GetOrderHash = client.GetOrderHash
var SignMessage = client.SignMessage
