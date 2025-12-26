package services

import (
	"fmt"
	"time"

	"github.com/extended-protocol/extended-sdk-golang/src/client"
	"github.com/extended-protocol/extended-sdk-golang/src/models"
	"github.com/shopspring/decimal"
)


// createOrderObjectParams represents the parameters for creating an order object
// All fields are required. For optional fields, pass nil or empty string.
type createOrderObjectParams struct {
	Market                   models.MarketModel
	Account                  *client.StarkPerpetualAccount
	SyntheticAmount          decimal.Decimal
	Price                    decimal.Decimal
	Side                     models.OrderSide
	Type                     models.OrderType
	StarknetDomain           models.StarknetDomain
	ExpireTime               time.Time
	PostOnly                 bool
	ReduceOnly               bool
	PreviousOrderExternalID  *string // Optional: pass nil if not canceling a previous order
	OrderExternalID          *string // Optional: pass nil to use order hash as ID
	TimeInForce              models.TimeInForce
	SelfTradeProtectionLevel models.SelfTradeProtectionLevel
	Nonce                    int
	BuilderFee               *decimal.Decimal // Optional: pass nil if no builder fee
	BuilderID                *int             // Optional: pass nil if no builder ID
	TpSlType                 *models.TpSlType        // Optional: TPSL type (ORDER or POSITION)
	TakeProfit               *models.TpSlTriggerParam // Optional: take profit trigger parameters
	StopLoss                 *models.TpSlTriggerParam // Optional: stop loss trigger parameters
}

// createOrderObject creates a PerpetualOrderModel with the given parameters
func createOrderObject(params createOrderObjectParams) (*models.PerpetualOrderModel, error) {
	// Validate side (must be BUY or SELL)
	if params.Side != models.OrderSideBuy && params.Side != models.OrderSideSell {
		return nil, fmt.Errorf("unexpected order side value: %s", params.Side)
	}

	// Validate time_in_force (must be GTT or IOC, not FOK)
	if params.TimeInForce == models.TimeInForceFOK {
		return nil, fmt.Errorf("unexpected time in force value: FOK is not supported")
	}
	if params.TimeInForce != models.TimeInForceGTT && params.TimeInForce != models.TimeInForceIOC {
		return nil, fmt.Errorf("unexpected time in force value: %s", params.TimeInForce)
	}

	// Validate expire_time is not zero
	if params.ExpireTime.IsZero() {
		return nil, fmt.Errorf("expire_time must be provided")
	}

	market := params.Market

	// If we are buying, then we round up, otherwise we round down
	is_buying_synthetic := params.Side == models.OrderSideBuy
	collateral_amount := params.SyntheticAmount.Mul(params.Price)

	// Get trading fees for the market. First check account's trading_fee cache,
	// then fall back to DefaultFees if not found.
	// Note: Fees are determined by the platform via GET /api/v1/user/fees?market={market}
	// and cannot be set by users. The team reserves the right to update the fee schedule. Currently:
	// - Taker: 0.025% (0.0005 in decimal)
	// - Maker: 0.000% (0.0000 in decimal)
	// To cache platform-determined fees, call AccountService.GetMarketFee() or AccountService.GetFees()
	// and then update the account's fee cache using account.SetTradingFee().
	fees := params.Account.GetTradingFee(params.Market.Name)

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
		Nonce:               params.Nonce,
		PositionID:          int(params.Account.Vault()),
		ExpirationTimestamp: params.ExpireTime,
		PublicKey:           params.Account.PublicKey(),
		StarknetDomain:      params.StarknetDomain,
	})

	if err != nil {
		return nil, fmt.Errorf("hashing order failed: %w", err)
	}

	sig_r, sig_s, err := params.Account.Sign(order_hash)
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

	// Use order hash as ID if OrderExternalID is not provided
	orderID := order_hash
	if params.OrderExternalID != nil {
		orderID = *params.OrderExternalID
	}

	order := &models.PerpetualOrderModel{
		ID:                       orderID,
		Market:                   params.Market.Name,
		Type:                     params.Type,
		Side:                     params.Side,
		Qty:                      params.SyntheticAmount.String(),
		Price:                    params.Price.String(),
		PostOnly:                 params.PostOnly,
		ReduceOnly:               params.ReduceOnly,
		TimeInForce:              params.TimeInForce,
		ExpiryEpochMillis:        expiryEpochMillis,
		Fee:                      fees.TakerFeeRate.String(),
		SelfTradeProtectionLevel: params.SelfTradeProtectionLevel,
		Nonce:                    fmt.Sprintf("%d", params.Nonce),
		CancelID:                 params.PreviousOrderExternalID,
		Settlement:               settlement,
		BuilderFee:               fee_builder_str,
		BuilderID:                params.BuilderID,
		// TPSL fields are set to nil for now - full implementation would require settlement data with opposite side
		TpSlType:   params.TpSlType,
		TakeProfit: nil,
		StopLoss:   nil,
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
// This function remains exported in case someone needs/wants to make their own implementation of the SDK but
// doesn't want to go through the trouble of implementing the hashing.
// It follows the same logic as the Python SDK, adding a 14 day buffer to the expiration timestamp.
func HashOrder(params HashOrderParams) (string, error) {
	expireTimeWithBuffer := params.ExpirationTimestamp.Add(14 * 24 * time.Hour)

	expireTimeRounded := expireTimeWithBuffer.Truncate(time.Second)
	if expireTimeWithBuffer.After(expireTimeRounded) {
		expireTimeRounded = expireTimeRounded.Add(time.Second)
	}

	expireTimeAsSeconds := expireTimeRounded.Unix()

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
		params.StarknetDomain.Revision,               // domain_revision
	)

	if err != nil {
		return "", fmt.Errorf("failed to compute order hash: %w", err)
	}

	return hash, nil
}

