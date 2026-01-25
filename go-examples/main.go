package main

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/matijamarjanovic/extended-sdk-golang/x10"
)

func main() {
	godotenv.Load()

	// read account credentials from environment
	vaultStr := os.Getenv("MAINNET_VAULT_ID")
	vault, _ := strconv.ParseUint(vaultStr, 10, 64)
	privateKey := os.Getenv("MAINNET_PRIVATE_KEY")
	publicKey := os.Getenv("MAINNET_PUBLIC_KEY")
	apiKey := os.Getenv("MAINNET_API_KEY")

	// create starknet account
	account, err := x10.NewStarkPerpetualAccount(vault, privateKey, publicKey, apiKey)
	if err != nil {
		panic(err)
	}

	// create client with mainnet configuration
	cfg := x10.STARKNET_MAINNET_CONFIG
	client := x10.NewClient(cfg, account, 30*time.Second)
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
