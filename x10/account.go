package x10

import "github.com/matijamarjanovic/extended-sdk-golang/x10/client"

// NewStarkPerpetualAccount creates a new Stark perpetual trading account.
// It wraps client.NewStarkPerpetualAccount to provide a root-level API.
// Panics if account creation fails.
func NewStarkPerpetualAccount(vault uint64, privateKeyHex, publicKeyHex, apiKey string) *StarkPerpetualAccount {
	acc, err := client.NewStarkPerpetualAccount(vault, privateKeyHex, publicKeyHex, apiKey)
	if err != nil {
		panic(err)
	}
	return (*StarkPerpetualAccount)(acc)
}