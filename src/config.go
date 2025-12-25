package sdk

import "github.com/extended-protocol/extended-sdk-golang/src/models"

// STARKNET_TESTNET_CONFIG is the pre-configured endpoint configuration for Starknet testnet
var STARKNET_TESTNET_CONFIG = models.EndpointConfig{
	ChainRPCURL:              "https://rpc.sepolia.org",
	APIBaseURL:               "https://api.starknet.sepolia.extended.exchange/api/v1",
	StreamURL:                "wss://api.starknet.sepolia.extended.exchange/stream.extended.exchange/v1",
	OnboardingURL:            "https://api.starknet.sepolia.extended.exchange",
	SigningDomain:            "starknet.sepolia.extended.exchange",
	CollateralAssetContract:  "0x31857064564ed0ff978e687456963cba09c2c6985d8f9300a1de4962fafa054",
	AssetOperationsContract:  "",
	CollateralAssetOnChainID: "0x1",
	CollateralDecimals:       6,
	StarknetDomain: models.StarknetDomain{
		Name:     "Perpetuals",
		Version:  "v0",
		ChainID:  "SN_SEPOLIA",
		Revision: "1",
	},
	CollateralAssetID: "0x1",
}

// STARKNET_MAINNET_CONFIG is the pre-configured endpoint configuration for Starknet mainnet
var STARKNET_MAINNET_CONFIG = models.EndpointConfig{
	ChainRPCURL:              "",
	APIBaseURL:               "https://api.starknet.extended.exchange/api/v1",
	StreamURL:                "wss://api.starknet.extended.exchange/stream.extended.exchange/v1",
	OnboardingURL:            "https://api.starknet.extended.exchange",
	SigningDomain:            "extended.exchange",
	CollateralAssetContract:  "",
	AssetOperationsContract:  "",
	CollateralAssetOnChainID: "0x1",
	CollateralDecimals:       6,
	StarknetDomain: models.StarknetDomain{
		Name:     "Perpetuals",
		Version:  "v0",
		ChainID:  "SN_MAIN",
		Revision: "1",
	},
	CollateralAssetID: "0x1",
}

