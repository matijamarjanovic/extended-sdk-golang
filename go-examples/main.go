package main

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	sdk "github.com/matijamarjanovic/extended-sdk-golang/x10"
)

func main() {
	godotenv.Load()

	// read account credentials from environment
	vaultStr := os.Getenv("TESTNET_VAULT_ID")
	vault, _ := strconv.ParseUint(vaultStr, 10, 64)
	privateKey := os.Getenv("TESTNET_PRIVATE_KEY")
	publicKey := os.Getenv("TESTNET_PUBLIC_KEY")
	apiKey := os.Getenv("TESTNET_API_KEY")

	// create starknet account
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
