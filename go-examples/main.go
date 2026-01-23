package examples

import (
	"os"
	"strconv"
	"time"

	sdk "github.com/matijamarjanovic/extended-sdk-golang/x10"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// read account credentials from environment
	vaultStr := os.Getenv("VAULT")
	vault, _ := strconv.ParseUint(vaultStr, 10, 64)
	privateKey := os.Getenv("PRIVATE_KEY")
	publicKey := os.Getenv("PUBLIC_KEY")
	apiKey := os.Getenv("API_KEY")

	// create starknet account (panics on error)
	account := sdk.NewStarkPerpetualAccount(vault, privateKey, publicKey, apiKey)

	// create client with testnet configuration
	cfg := sdk.STARKNET_TESTNET_CONFIG
	client := sdk.NewClient(cfg, account, 30*time.Second)
	defer client.Close()

	if len(os.Args) < 2 {
		return
	}

	switch os.Args[1] {
	case "account":
		accountExample(client)
	case "markets":
		marketsExample(client)
	case "orders":
		ordersExample(client)
	case "streaming":
		streamingExample(client)
	}
}
