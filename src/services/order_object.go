package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
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
	PreviousOrderExternalID  *string // optional: pass nil if not canceling a previous order
	OrderExternalID          *string // optional: pass nil to use order hash as ID
	TimeInForce              models.TimeInForce
	SelfTradeProtectionLevel models.SelfTradeProtectionLevel
	Nonce                    *int                     // optional: pass nil to auto-generate
	BuilderFee               *decimal.Decimal         // optional: pass nil if no builder fee
	BuilderID                *int                     // optional: pass nil if no builder ID
	TpSlType                 *models.TpSlType         // optional: TPSL type (ORDER or POSITION)
	TakeProfit               *models.TpSlTriggerParam // optional: take profit trigger parameters
	StopLoss                 *models.TpSlTriggerParam // optional: stop loss trigger parameters
}

// createOrderObject creates a PerpetualOrderModel with the given parameters
func createOrderObject(params createOrderObjectParams) (*models.PerpetualOrderModel, error) {
	if params.Side != models.OrderSideBuy && params.Side != models.OrderSideSell {
		return nil, fmt.Errorf("unexpected order side value: %s", params.Side)
	}

	if params.TimeInForce == models.TimeInForceFOK {
		return nil, fmt.Errorf("unexpected time in force value: FOK is not supported")
	}
	if params.TimeInForce != models.TimeInForceGTT && params.TimeInForce != models.TimeInForceIOC {
		return nil, fmt.Errorf("unexpected time in force value: %s", params.TimeInForce)
	}

	if params.ExpireTime.IsZero() {
		return nil, fmt.Errorf("expire_time must be provided")
	}

	// auto-generate nonce if not provided
	nonce := params.Nonce
	if nonce == nil {
		// generate random nonce in range [0, 2^32 - 1]
		maxNonce := big.NewInt(1<<32 - 1)
		generatedNonce, err := rand.Int(rand.Reader, maxNonce)
		if err != nil {
			return nil, fmt.Errorf("failed to generate nonce: %w", err)
		}
		nonceValue := int(generatedNonce.Int64())
		nonce = &nonceValue
	}

	market := params.Market

	// if we are buying, then we round up, otherwise we round down
	is_buying_synthetic := params.Side == models.OrderSideBuy
	collateral_amount := params.SyntheticAmount.Mul(params.Price)

	// get trading fees for the market. First check account's trading_fee cache,
	// then fall back to DefaultFees if not found.
	// https://api.docs.extended.exchange/#get-fees
	fees := params.Account.GetTradingFee(params.Market.Name)

	total_fee := fees.TakerFeeRate
	if params.BuilderFee != nil {
		total_fee = total_fee.Add(*params.BuilderFee)
	}

	fee_amount := total_fee.Mul(collateral_amount)

	stark_collateral_amount_dec := collateral_amount.Mul(decimal.NewFromInt(market.L2Config.CollateralResolution))
	stark_synthetic_amount_dec := params.SyntheticAmount.Mul(decimal.NewFromInt(market.L2Config.SyntheticResolution))

	// round accordingly
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
		Nonce:               *nonce,
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

	// convert expire time to epoch milliseconds
	expiryEpochMillis := params.ExpireTime.UnixNano() / int64(time.Millisecond)

	// use order hash as ID if OrderExternalID is not provided
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
		Nonce:                    fmt.Sprintf("%d", *nonce),
		CancelID:                 params.PreviousOrderExternalID,
		Settlement:               settlement,
		BuilderFee:               fee_builder_str,
		BuilderID:                params.BuilderID,
		TpSlType:                 params.TpSlType,
	}

	// create TPSL triggers - if requested, they must succeed or the entire order creation fails
	// use the resolved nonce (same as main order) for TPSL triggers
	if params.TakeProfit != nil {
		takeProfit, err := createTpSlTrigger(params.TakeProfit, params, market, fees, *nonce)
		if err != nil {
			return nil, fmt.Errorf("failed to create take profit trigger: %w", err)
		}
		order.TakeProfit = takeProfit
	}

	if params.StopLoss != nil {
		stopLoss, err := createTpSlTrigger(params.StopLoss, params, market, fees, *nonce)
		if err != nil {
			return nil, fmt.Errorf("failed to create stop loss trigger: %w", err)
		}
		order.StopLoss = stopLoss
	}

	return order, nil
}

// getOppositeSide returns the opposite order side
func getOppositeSide(side models.OrderSide) models.OrderSide {
	if side == models.OrderSideBuy {
		return models.OrderSideSell
	}
	return models.OrderSideBuy
}

// createTpSlTrigger creates a TPSL trigger model with settlement data for the opposite side order
// Returns an error if TPSL trigger creation fails - this is critical for trading safety
// nonce must be provided (same as the main order nonce)
func createTpSlTrigger(
	triggerParam *models.TpSlTriggerParam,
	params createOrderObjectParams,
	market models.MarketModel,
	fees models.TradingFeeModel,
	nonce int,
) (*models.TpSlTrigger, error) {
	if triggerParam == nil {
		return nil, nil
	}

	oppositeSide := getOppositeSide(params.Side)
	is_buying_synthetic := oppositeSide == models.OrderSideBuy

	// use the TPSL trigger price (not the main order price)
	tpslPrice := triggerParam.Price
	collateral_amount := params.SyntheticAmount.Mul(tpslPrice)

	// calculate fees
	total_fee := fees.TakerFeeRate
	if params.BuilderFee != nil {
		total_fee = total_fee.Add(*params.BuilderFee)
	}
	fee_amount := total_fee.Mul(collateral_amount)

	// convert to Stark amounts
	stark_collateral_amount_dec := collateral_amount.Mul(decimal.NewFromInt(market.L2Config.CollateralResolution))
	stark_synthetic_amount_dec := params.SyntheticAmount.Mul(decimal.NewFromInt(market.L2Config.SyntheticResolution))

	// round accordingly
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

	// apply sign based on side
	if is_buying_synthetic {
		stark_collateral_amount = -stark_collateral_amount
	} else {
		stark_synthetic_amount = -stark_synthetic_amount
	}

	tpslOrderHash, err := HashOrder(HashOrderParams{
		AmountSynthetic:     stark_synthetic_amount,
		SyntheticAssetID:    market.L2Config.SyntheticID,
		AmountCollateral:    stark_collateral_amount,
		CollateralAssetID:   market.L2Config.CollateralID,
		MaxFee:              stark_fee_part,
		Nonce:               nonce,
		PositionID:          int(params.Account.Vault()),
		ExpirationTimestamp: params.ExpireTime,
		PublicKey:           params.Account.PublicKey(),
		StarknetDomain:      params.StarknetDomain,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to hash TPSL order: %w", err)
	}

	sig_r, sig_s, err := params.Account.Sign(tpslOrderHash)
	if err != nil {
		return nil, fmt.Errorf("failed to sign TPSL order: %w", err)
	}

	tpslSettlement := models.Settlement{
		Signature: models.Signature{
			R: fmt.Sprintf("0x%x", sig_r),
			S: fmt.Sprintf("0x%x", sig_s),
		},
		StarkKey:           params.Account.PublicKey(),
		CollateralPosition: fmt.Sprintf("%d", params.Account.Vault()),
	}

	return &models.TpSlTrigger{
		TriggerPrice:     triggerParam.TriggerPrice.String(),
		TriggerPriceType: triggerParam.TriggerPriceType,
		Price:            triggerParam.Price.String(),
		PriceType:        triggerParam.PriceType,
		Settlement:       tpslSettlement,
	}, nil
}

// HashOrderParams represents the parameters for hashing an order
type HashOrderParams struct {
	AmountSynthetic     int64
	SyntheticAssetID    string // hex string for asset id
	AmountCollateral    int64
	CollateralAssetID   string // hex string for asset id
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
		params.StarknetDomain.Revision,             // domain_revision
	)

	if err != nil {
		return "", fmt.Errorf("failed to compute order hash: %w", err)
	}

	return hash, nil
}
