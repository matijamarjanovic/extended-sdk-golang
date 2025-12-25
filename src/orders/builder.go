package orders

import (
	"fmt"
	"math/big"
	"time"

	"github.com/extended-protocol/extended-sdk-golang/src/client"
	"github.com/extended-protocol/extended-sdk-golang/src/models"
	"github.com/shopspring/decimal"
)

// CreateOrderObjectParams represents the parameters for creating an order object
type CreateOrderObjectParams struct {
	Market                   models.MarketModel
	Account                  *client.StarkPerpetualAccount
	SyntheticAmount          decimal.Decimal
	Price                    decimal.Decimal
	Side                     models.OrderSide
	Signer                   func(string) (*big.Int, *big.Int, error) // Function that takes string and returns two values
	StarknetDomain           models.StarknetDomain
	ExpireTime               *time.Time
	PostOnly                 bool
	PreviousOrderExternalID  *string
	OrderExternalID          *string
	TimeInForce              models.TimeInForce
	SelfTradeProtectionLevel models.SelfTradeProtectionLevel
	Nonce                    *int
	BuilderFee               *decimal.Decimal
	BuilderID                *int
}

// CreateOrderObject creates a PerpetualOrderModel with the given parameters
func CreateOrderObject(params CreateOrderObjectParams) (*models.PerpetualOrderModel, error) {
	market := params.Market

	if params.ExpireTime == nil {
		cur := time.Now().Add(1 * time.Hour)
		params.ExpireTime = &cur
	}

	// Error if nonce is nil, we keep the input as a pointer so that
	// it is the same as the input to the function
	if params.Nonce == nil {
		return nil, fmt.Errorf("nonce must be provided")
	}

	// If we are buying, then we round up, otherwise we round down
	is_buying_synthetic := params.Side == models.OrderSideBuy
	collateral_amount := params.SyntheticAmount.Mul(params.Price)

	// For now we only use the default fee type
	// TODO: Allow users to add different fee types
	fees := models.DefaultFees

	total_fee := fees.TakerFeeRate
	if params.BuilderFee != nil {
		total_fee = total_fee.Add(*params.BuilderFee)
	}

	fee_amount := total_fee.Mul(collateral_amount)

	stark_collateral_amount_dec := collateral_amount.Mul(decimal.NewFromInt(market.L2Config.CollateralResolution))
	stark_synthetic_amount_dec := params.SyntheticAmount.Mul(decimal.NewFromInt(market.L2Config.SyntheticResolution))

	// Round accordingly
	if is_buying_synthetic {
		stark_collateral_amount_dec = stark_collateral_amount_dec.Ceil()
		stark_synthetic_amount_dec = stark_synthetic_amount_dec.Ceil()
	} else {
		stark_collateral_amount_dec = stark_collateral_amount_dec.Floor()
		stark_synthetic_amount_dec = stark_synthetic_amount_dec.Floor()
	}

	stark_collateral_amount := stark_collateral_amount_dec.IntPart()
	stark_synthetic_amount := stark_synthetic_amount_dec.IntPart()
	stark_fee_part := fee_amount.Mul(decimal.NewFromInt(market.L2Config.CollateralResolution)).Ceil().IntPart()

	if is_buying_synthetic {
		stark_collateral_amount = -stark_collateral_amount
	} else {
		stark_synthetic_amount = -stark_synthetic_amount
	}

	order_hash, err := HashOrder(HashOrderParams{
		AmountSynthetic:     stark_synthetic_amount,
		SyntheticAssetID:    market.L2Config.SyntheticID,
		AmountCollateral:    stark_collateral_amount,
		CollateralAssetID:   market.L2Config.CollateralID,
		MaxFee:              stark_fee_part,
		Nonce:               *params.Nonce,
		PositionID:          int(params.Account.Vault()),
		ExpirationTimestamp: *params.ExpireTime,
		PublicKey:           params.Account.PublicKey(),
		StarknetDomain:      params.StarknetDomain,
	})

	if err != nil {
		return nil, fmt.Errorf("hashing order failed: %w", err)
	}

	sig_r, sig_s, err := params.Signer(order_hash)
	if err != nil {
		return nil, fmt.Errorf("signer function failed: %w", err)
	}

	settlement := models.Settlement{
		Signature: models.Signature{
			R: fmt.Sprintf("0x%x", sig_r),
			S: fmt.Sprintf("0x%x", sig_s),
		},
		StarkKey:           params.Account.PublicKey(),
		CollateralPosition: fmt.Sprintf("%d", params.Account.Vault()),
	}

	if params.OrderExternalID == nil {
		defaultID := order_hash
		params.OrderExternalID = &defaultID
	}

	var fee_builder_str *string
	if params.BuilderFee != nil {
		builderFeeStr := params.BuilderFee.String()
		fee_builder_str = &builderFeeStr
	}

	// Convert expire time to epoch milliseconds
	expiryEpochMillis := params.ExpireTime.UnixNano() / int64(time.Millisecond)

	order := &models.PerpetualOrderModel{
		ID:                       *params.OrderExternalID,
		Market:                   params.Market.Name,
		Type:                     models.OrderTypeLimit,
		Side:                     params.Side,
		Qty:                      params.SyntheticAmount.String(),
		Price:                    params.Price.String(),
		PostOnly:                 params.PostOnly,
		TimeInForce:              params.TimeInForce,
		ExpiryEpochMillis:        expiryEpochMillis,
		Fee:                      fees.TakerFeeRate.String(),
		SelfTradeProtectionLevel: params.SelfTradeProtectionLevel,
		Nonce:                    fmt.Sprintf("%d", *params.Nonce),
		CancelID:                 params.PreviousOrderExternalID,
		Settlement:               settlement,
		BuilderFee:               fee_builder_str,
		BuilderID:                params.BuilderID,
	}

	return order, nil
}

// HashOrderParams represents the parameters for hashing an order
type HashOrderParams struct {
	AmountSynthetic     int64
	SyntheticAssetID    string // hex string for asset ID
	AmountCollateral    int64
	CollateralAssetID   string // hex string for asset ID
	MaxFee              int64
	Nonce               int
	PositionID          int
	ExpirationTimestamp time.Time
	PublicKey           string
	StarknetDomain      models.StarknetDomain
}

// HashOrder computes the order hash using the provided parameters.
// This mimics the Python hash_order function
func HashOrder(params HashOrderParams) (string, error) {
	// Add 14 days buffer to expiration timestamp
	expireTimeWithBuffer := params.ExpirationTimestamp.Add(14 * 24 * time.Hour)

	// Round UP to the nearest second
	expireTimeRounded := expireTimeWithBuffer.Truncate(time.Second)
	if expireTimeWithBuffer.After(expireTimeRounded) {
		expireTimeRounded = expireTimeRounded.Add(time.Second)
	}

	expireTimeAsSeconds := expireTimeRounded.Unix()

	// Call GetOrderHash from client package
	hash, err := client.GetOrderHash(
		fmt.Sprintf("%d", params.PositionID),       // position_id
		params.SyntheticAssetID,                    // base_asset_id_hex
		fmt.Sprintf("%d", params.AmountSynthetic),  // base_amount
		params.CollateralAssetID,                   // quote_asset_id_hex
		fmt.Sprintf("%d", params.AmountCollateral), // quote_amount
		params.CollateralAssetID,                   // fee_asset_id_hex (same as collateral)
		fmt.Sprintf("%d", params.MaxFee),           // fee_amount
		fmt.Sprintf("%d", expireTimeAsSeconds),     // expiration
		fmt.Sprintf("%d", params.Nonce),            // salt (nonce)
		params.PublicKey,                           // user_public_key_hex
		params.StarknetDomain.Name,                 // domain_name
		params.StarknetDomain.Version,              // domain_version
		params.StarknetDomain.ChainID,              // domain_chain_id
		params.StarknetDomain.Revision,             // domain_revision
	)

	if err != nil {
		return "", fmt.Errorf("failed to compute order hash: %w", err)
	}

	return hash, nil
}
