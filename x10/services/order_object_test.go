package services

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/matijamarjanovic/extended-sdk-golang/x10/client"
	"github.com/matijamarjanovic/extended-sdk-golang/x10/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

const (
	TestPrivateKeyHex = "0x7a7ff6fd3cab02ccdcd4a572563f5976f8976899b03a39773795a3c486d4986"
	TestPublicKeyHex  = "0x61c5e7e8339b7d56f197f54ea91b776776690e3232313de0f2ecbd0ef76f466"
	TestVaultID       = 10002
	TestAPIKey        = "test-api-key"
	TestNonce         = 1473459052
)

func createTestAccount() (*client.StarkPerpetualAccount, error) {
	return client.NewStarkPerpetualAccount(TestVaultID, TestPrivateKeyHex, TestPublicKeyHex, TestAPIKey)
}

func createTestBTCUSDMarket() models.MarketModel {
	return models.MarketModel{
		Name:                     "BTC-USD",
		AssetName:                "BTC",
		AssetPrecision:           8,
		CollateralAssetName:      "USD",
		CollateralAssetPrecision: 6,
		Active:                   true,
		L2Config: models.L2ConfigModel{
			Type:                 "perpetual",
			CollateralID:         "0x31857064564ed0ff978e687456963cba09c2c6985d8f9300a1de4962fafa054",
			CollateralResolution: 1000000,
			SyntheticID:          "0x4254432d3600000000000000000000",
			SyntheticResolution:  1000000,
		},
	}
}

func createTestStarknetDomain() models.StarknetDomain {
	return models.StarknetDomain{
		Name:     "Perpetuals",
		Version:  "v0",
		ChainID:  "SN_SEPOLIA",
		Revision: "1",
	}
}

func createTestFrozenTime() time.Time {
	return time.Date(2024, 1, 5, 1, 8, 57, 0, time.UTC)
}

type OrdersTestSuite struct {
	suite.Suite
	account        *client.StarkPerpetualAccount
	market         models.MarketModel
	starknetDomain models.StarknetDomain
	frozenTime     time.Time
	nonce          int
}

func (suite *OrdersTestSuite) SetupTest() {
	var err error
	suite.account, err = createTestAccount()
	suite.Require().NoError(err)

	suite.market = createTestBTCUSDMarket()
	suite.starknetDomain = createTestStarknetDomain()
	suite.frozenTime = createTestFrozenTime()
	suite.nonce = TestNonce
}

func (suite *OrdersTestSuite) TestCreateSellOrderWithDefaultExpiration() {
	expireTime := suite.frozenTime.Add(1 * time.Hour)
	params := createOrderObjectParams{
		Market:                   suite.market,
		Account:                  suite.account,
		SyntheticAmount:          decimal.RequireFromString("0.00100000"),
		Price:                    decimal.RequireFromString("43445.11680000"),
		Side:                     models.OrderSideSell,
		Type:                     models.OrderTypeLimit,
		StarknetDomain:           suite.starknetDomain,
		ExpireTime:               expireTime,
		PostOnly:                 false,
		PreviousOrderExternalID:  nil,
		OrderExternalID:          nil,
		TimeInForce:              models.TimeInForceGTT,
		SelfTradeProtectionLevel: models.SelfTradeProtectionAccount,
		Nonce:                    &suite.nonce,
		BuilderFee:               nil,
		BuilderID:                nil,
	}

	order, err := createOrderObject(params)
	suite.Require().NoError(err)
	suite.Require().NotNil(order)

	orderJSON, err := json.Marshal(order)
	suite.Require().NoError(err)

	var actualOrder map[string]interface{}
	err = json.Unmarshal(orderJSON, &actualOrder)
	suite.Require().NoError(err)

	// expected JSON structure
	expectedOrder := map[string]interface{}{
		"id":                       "529621978301228831750156704671293558063128025271079340676658105549022202327",
		"market":                   "BTC-USD",
		"type":                     "LIMIT",
		"side":                     "SELL",
		"qty":                      "0.001",
		"price":                    "43445.1168",
		"reduceOnly":               false,
		"postOnly":                 false,
		"timeInForce":              "GTT",
		"fee":                      "0.0005",
		"nonce":                    "1473459052",
		"selfTradeProtectionLevel": "ACCOUNT",
		"cancelId":                 nil,
		"trigger":                  nil,
		"tpSlType":                 nil,
		"takeProfit":               nil,
		"stopLoss":                 nil,
		"builderFee":               nil,
		"builderId":                nil,
	}

	suite.Equal(expectedOrder["market"], actualOrder["market"])
	suite.Equal(expectedOrder["type"], actualOrder["type"])
	suite.Equal(expectedOrder["side"], actualOrder["side"])
	suite.Equal(expectedOrder["qty"], actualOrder["qty"])
	suite.Equal(expectedOrder["price"], actualOrder["price"])
	suite.Equal(expectedOrder["reduceOnly"], actualOrder["reduceOnly"])
	suite.Equal(expectedOrder["postOnly"], actualOrder["postOnly"])
	suite.Equal(expectedOrder["timeInForce"], actualOrder["timeInForce"])
	suite.Equal(expectedOrder["fee"], actualOrder["fee"])
	suite.Equal(expectedOrder["nonce"], actualOrder["nonce"])
	suite.Equal(expectedOrder["selfTradeProtectionLevel"], actualOrder["selfTradeProtectionLevel"])
	suite.Equal(expectedOrder["cancelId"], actualOrder["cancelId"])
	suite.Equal(expectedOrder["trigger"], actualOrder["trigger"])
	suite.Equal(expectedOrder["tpSlType"], actualOrder["tpSlType"])
	suite.Equal(expectedOrder["takeProfit"], actualOrder["takeProfit"])
	suite.Equal(expectedOrder["stopLoss"], actualOrder["stopLoss"])
	suite.Equal(expectedOrder["builderFee"], actualOrder["builderFee"])
	suite.Equal(expectedOrder["builderId"], actualOrder["builderId"])

	suite.NotEmpty(actualOrder["id"])
}

func (suite *OrdersTestSuite) TestCreateSellOrder() {
	// set expiry time (1 hour from frozen time = 1704420537000 milliseconds)
	expiryTime := suite.frozenTime.Add(1 * time.Hour)

	params := createOrderObjectParams{
		Market:                   suite.market,
		Account:                  suite.account,
		SyntheticAmount:          decimal.RequireFromString("0.00100000"),
		Price:                    decimal.RequireFromString("43445.11680000"),
		Side:                     models.OrderSideSell,
		Type:                     models.OrderTypeLimit,
		StarknetDomain:           suite.starknetDomain,
		ExpireTime:               expiryTime,
		PostOnly:                 false,
		PreviousOrderExternalID:  nil,
		OrderExternalID:          nil,
		TimeInForce:              models.TimeInForceGTT,
		SelfTradeProtectionLevel: models.SelfTradeProtectionAccount,
		Nonce:                    &suite.nonce,
		BuilderFee:               nil,
		BuilderID:                nil,
	}

	order, err := createOrderObject(params)
	suite.Require().NoError(err)
	suite.Require().NotNil(order)

	orderJSON, err := json.Marshal(order)
	suite.Require().NoError(err)

	var actualOrder map[string]interface{}
	err = json.Unmarshal(orderJSON, &actualOrder)
	suite.Require().NoError(err)

	expectedOrder := map[string]interface{}{
		"id":                       "529621978301228831750156704671293558063128025271079340676658105549022202327",
		"market":                   "BTC-USD",
		"type":                     "LIMIT",
		"side":                     "SELL",
		"qty":                      "0.001",
		"price":                    "43445.1168",
		"reduceOnly":               false,
		"postOnly":                 false,
		"timeInForce":              "GTT",
		"expiryEpochMillis":        float64(1704420537000),
		"fee":                      "0.0005",
		"nonce":                    "1473459052",
		"selfTradeProtectionLevel": "ACCOUNT",
		"cancelId":                 nil,
		"settlement": map[string]interface{}{
			"signature": map[string]interface{}{
				"r": "0x3d17d8b9652e5f60d40d079653cfa92b1065ea8cf159609a3c390070dcd44f7",
				"s": "0x76a6deccbc84ac324f695cfbde80e0ed62443e95f5dcd8722d12650ccc122e5",
			},
			"starkKey":           TestPublicKeyHex,
			"collateralPosition": "10002",
		},
		"trigger":    nil,
		"tpSlType":   nil,
		"takeProfit": nil,
		"stopLoss":   nil,
		"builderFee": nil,
		"builderId":  nil,
	}

	suite.Equal(expectedOrder["market"], actualOrder["market"])
	suite.Equal(expectedOrder["type"], actualOrder["type"])
	suite.Equal(expectedOrder["side"], actualOrder["side"])
	suite.Equal(expectedOrder["qty"], actualOrder["qty"])
	suite.Equal(expectedOrder["price"], actualOrder["price"])
	suite.Equal(expectedOrder["reduceOnly"], actualOrder["reduceOnly"])
	suite.Equal(expectedOrder["postOnly"], actualOrder["postOnly"])
	suite.Equal(expectedOrder["timeInForce"], actualOrder["timeInForce"])
	suite.Equal(expectedOrder["expiryEpochMillis"], actualOrder["expiryEpochMillis"])
	suite.Equal(expectedOrder["fee"], actualOrder["fee"])
	suite.Equal(expectedOrder["nonce"], actualOrder["nonce"])
	suite.Equal(expectedOrder["selfTradeProtectionLevel"], actualOrder["selfTradeProtectionLevel"])
	suite.Equal(expectedOrder["cancelId"], actualOrder["cancelId"])
	suite.Equal(expectedOrder["settlement"], actualOrder["settlement"])
	suite.Equal(expectedOrder["trigger"], actualOrder["trigger"])
	suite.Equal(expectedOrder["tpSlType"], actualOrder["tpSlType"])
	suite.Equal(expectedOrder["takeProfit"], actualOrder["takeProfit"])
	suite.Equal(expectedOrder["stopLoss"], actualOrder["stopLoss"])
	suite.Equal(expectedOrder["builderFee"], actualOrder["builderFee"])
	suite.Equal(expectedOrder["builderId"], actualOrder["builderId"])

	suite.NotEmpty(actualOrder["id"])
}

func (suite *OrdersTestSuite) TestCreateBuyOrderWithClientProtection() {
	// set expiry time (1 hour from frozen time)
	// @freeze_time("2024-01-05 01:08:56.860694")
	expiryTime := time.Date(2024, 1, 5, 1, 8, 56, 860694000, time.UTC).Add(14 * 24 * time.Hour)

	params := createOrderObjectParams{
		Market:                   suite.market,
		Account:                  suite.account,
		SyntheticAmount:          decimal.RequireFromString("0.00100000"),
		Price:                    decimal.RequireFromString("43445.11680000"),
		Side:                     models.OrderSideBuy,
		Type:                     models.OrderTypeLimit,
		StarknetDomain:           suite.starknetDomain,
		ExpireTime:               expiryTime,
		PostOnly:                 false,
		PreviousOrderExternalID:  nil,
		OrderExternalID:          nil,
		TimeInForce:              models.TimeInForceGTT,
		SelfTradeProtectionLevel: models.SelfTradeProtectionClient,
		Nonce:                    &suite.nonce,
		BuilderFee:               nil,
		BuilderID:                nil,
	}

	order, err := createOrderObject(params)
	suite.Require().NoError(err)
	suite.Require().NotNil(order)

	orderJSON, err := json.Marshal(order)
	suite.Require().NoError(err)

	var actualOrder map[string]interface{}
	err = json.Unmarshal(orderJSON, &actualOrder)
	suite.Require().NoError(err)

	expectedOrder := map[string]interface{}{
		"market":                   "BTC-USD",
		"type":                     "LIMIT",
		"side":                     "BUY",
		"qty":                      "0.001",
		"price":                    "43445.1168",
		"reduceOnly":               false,
		"postOnly":                 false,
		"timeInForce":              "GTT",
		"expiryEpochMillis":        float64(1705626536861),
		"fee":                      "0.0005",
		"nonce":                    "1473459052",
		"selfTradeProtectionLevel": "CLIENT",
		"cancelId":                 nil,
		"settlement": map[string]interface{}{
			"signature": map[string]interface{}{
				"r": "0xa55625c7d5f1b85bed22556fc805224b8363074979cf918091d9ddb1403e13",
				"s": "0x504caf634d859e643569743642ccf244434322859b2421d76f853af43ae7a46",
			},
			"starkKey":           TestPublicKeyHex,
			"collateralPosition": "10002",
		},
		"trigger":    nil,
		"tpSlType":   nil,
		"takeProfit": nil,
		"stopLoss":   nil,
		"builderFee": nil,
		"builderId":  nil,
	}

	suite.Equal(expectedOrder["market"], actualOrder["market"])
	suite.Equal(expectedOrder["type"], actualOrder["type"])
	suite.Equal(expectedOrder["side"], actualOrder["side"])
	suite.Equal(expectedOrder["qty"], actualOrder["qty"])
	suite.Equal(expectedOrder["price"], actualOrder["price"])
	suite.Equal(expectedOrder["selfTradeProtectionLevel"], actualOrder["selfTradeProtectionLevel"])
	suite.Equal(expectedOrder["settlement"], actualOrder["settlement"])
	suite.NotEmpty(actualOrder["id"])
}

func (suite *OrdersTestSuite) TestCancelPreviousOrder() {
	// set expiry time (1 hour from frozen time)
	expiryTime := suite.frozenTime.Add(1 * time.Hour)
	previousOrderID := "previous_custom_id"

	params := createOrderObjectParams{
		Market:                   suite.market,
		Account:                  suite.account,
		SyntheticAmount:          decimal.RequireFromString("0.00100000"),
		Price:                    decimal.RequireFromString("43445.11680000"),
		Side:                     models.OrderSideBuy,
		Type:                     models.OrderTypeLimit,
		StarknetDomain:           suite.starknetDomain,
		ExpireTime:               expiryTime,
		PostOnly:                 false,
		PreviousOrderExternalID:  &previousOrderID,
		OrderExternalID:          nil,
		TimeInForce:              models.TimeInForceGTT,
		SelfTradeProtectionLevel: models.SelfTradeProtectionAccount,
		Nonce:                    &suite.nonce,
		BuilderFee:               nil,
		BuilderID:                nil,
	}

	order, err := createOrderObject(params)
	suite.Require().NoError(err)
	suite.Require().NotNil(order)

	orderJSON, err := json.Marshal(order)
	suite.Require().NoError(err)

	var actualOrder map[string]interface{}
	err = json.Unmarshal(orderJSON, &actualOrder)
	suite.Require().NoError(err)

	suite.Equal(previousOrderID, actualOrder["cancelId"])
}

func (suite *OrdersTestSuite) TestExternalOrderID() {
	expiryTime := suite.frozenTime.Add(1 * time.Hour)
	customOrderID := "custom_id"

	params := createOrderObjectParams{
		Market:                   suite.market,
		Account:                  suite.account,
		SyntheticAmount:          decimal.RequireFromString("0.00100000"),
		Price:                    decimal.RequireFromString("43445.11680000"),
		Side:                     models.OrderSideBuy,
		Type:                     models.OrderTypeLimit,
		StarknetDomain:           suite.starknetDomain,
		ExpireTime:               expiryTime,
		PostOnly:                 false,
		PreviousOrderExternalID:  nil,
		OrderExternalID:          &customOrderID,
		TimeInForce:              models.TimeInForceGTT,
		SelfTradeProtectionLevel: models.SelfTradeProtectionAccount,
		Nonce:                    &suite.nonce,
		BuilderFee:               nil,
		BuilderID:                nil,
	}

	order, err := createOrderObject(params)
	suite.Require().NoError(err)
	suite.Require().NotNil(order)

	orderJSON, err := json.Marshal(order)
	suite.Require().NoError(err)

	var actualOrder map[string]interface{}
	err = json.Unmarshal(orderJSON, &actualOrder)
	suite.Require().NoError(err)

	suite.Equal(customOrderID, actualOrder["id"])
}

func (suite *OrdersTestSuite) TestCreateBuyOrderWithTakeProfit() {
	expiryTime := suite.frozenTime.Add(14 * 24 * time.Hour)

	tpSlType := models.TpSlTypeOrder
	takeProfit := &models.TpSlTriggerParam{
		TriggerPrice:     decimal.RequireFromString("49000"),
		TriggerPriceType: models.TriggerPriceTypeMark,
		Price:            decimal.RequireFromString("50000"),
		PriceType:        models.ExecutionPriceTypeLimit,
	}

	params := createOrderObjectParams{
		Market:                   suite.market,
		Account:                  suite.account,
		SyntheticAmount:          decimal.RequireFromString("0.00100000"),
		Price:                    decimal.RequireFromString("43445.11680000"),
		Side:                     models.OrderSideBuy,
		Type:                     models.OrderTypeLimit,
		StarknetDomain:           suite.starknetDomain,
		ExpireTime:               expiryTime,
		PostOnly:                 false,
		PreviousOrderExternalID:  nil,
		OrderExternalID:          nil,
		TimeInForce:              models.TimeInForceGTT,
		SelfTradeProtectionLevel: models.SelfTradeProtectionClient,
		Nonce:                    &suite.nonce,
		BuilderFee:               nil,
		BuilderID:                nil,
		TpSlType:                 &tpSlType,
		TakeProfit:               takeProfit,
		StopLoss:                 nil,
	}

	order, err := createOrderObject(params)
	suite.Require().NoError(err)
	suite.Require().NotNil(order)

	suite.NotNil(order.TakeProfit, "take profit should be set")
	suite.Nil(order.StopLoss, "stop loss should be nil")
	suite.Equal(&tpSlType, order.TpSlType, "TPSL type should be set")

	tp := order.TakeProfit
	suite.Equal("49000", tp.TriggerPrice, "trigger price should match")
	suite.Equal(models.TriggerPriceTypeMark, tp.TriggerPriceType, "trigger price type should match")
	suite.Equal("50000", tp.Price, "price should match")
	suite.Equal(models.ExecutionPriceTypeLimit, tp.PriceType, "price type should match")
	suite.NotEmpty(tp.Settlement.Signature.R, "settlement signature R should be set")
	suite.NotEmpty(tp.Settlement.Signature.S, "settlement signature S should be set")
	suite.Equal(TestPublicKeyHex, tp.Settlement.StarkKey, "settlement stark key should match")
	suite.Equal("10002", tp.Settlement.CollateralPosition, "settlement collateral position should match")
}

func (suite *OrdersTestSuite) TestCreateBuyOrderWithStopLoss() {
	expiryTime := suite.frozenTime.Add(14 * 24 * time.Hour)

	tpSlType := models.TpSlTypeOrder
	stopLoss := &models.TpSlTriggerParam{
		TriggerPrice:     decimal.RequireFromString("40000"),
		TriggerPriceType: models.TriggerPriceTypeMark,
		Price:            decimal.RequireFromString("39000"),
		PriceType:        models.ExecutionPriceTypeLimit,
	}

	params := createOrderObjectParams{
		Market:                   suite.market,
		Account:                  suite.account,
		SyntheticAmount:          decimal.RequireFromString("0.00100000"),
		Price:                    decimal.RequireFromString("43445.11680000"),
		Side:                     models.OrderSideBuy,
		Type:                     models.OrderTypeLimit,
		StarknetDomain:           suite.starknetDomain,
		ExpireTime:               expiryTime,
		PostOnly:                 false,
		PreviousOrderExternalID:  nil,
		OrderExternalID:          nil,
		TimeInForce:              models.TimeInForceGTT,
		SelfTradeProtectionLevel: models.SelfTradeProtectionClient,
		Nonce:                    &suite.nonce,
		BuilderFee:               nil,
		BuilderID:                nil,
		TpSlType:                 &tpSlType,
		TakeProfit:               nil,
		StopLoss:                 stopLoss,
	}

	order, err := createOrderObject(params)
	suite.Require().NoError(err)
	suite.Require().NotNil(order)

	suite.Nil(order.TakeProfit, "take profit should be nil")
	suite.NotNil(order.StopLoss, "stop loss should be set")
	suite.Equal(&tpSlType, order.TpSlType, "TPSL type should be set")

	sl := order.StopLoss
	suite.Equal("40000", sl.TriggerPrice, "trigger price should match")
	suite.Equal(models.TriggerPriceTypeMark, sl.TriggerPriceType, "trigger price type should match")
	suite.Equal("39000", sl.Price, "price should match")
	suite.Equal(models.ExecutionPriceTypeLimit, sl.PriceType, "price type should match")
	suite.NotEmpty(sl.Settlement.Signature.R, "settlement signature R should be set")
	suite.NotEmpty(sl.Settlement.Signature.S, "settlement signature S should be set")
	suite.Equal(TestPublicKeyHex, sl.Settlement.StarkKey, "settlement stark key should match")
	suite.Equal("10002", sl.Settlement.CollateralPosition, "settlement collateral position should match")
}

func (suite *OrdersTestSuite) TestCreateBuyOrderWithBothTPSL() {
	expiryTime := suite.frozenTime.Add(14 * 24 * time.Hour)

	tpSlType := models.TpSlTypeOrder
	takeProfit := &models.TpSlTriggerParam{
		TriggerPrice:     decimal.RequireFromString("49000"),
		TriggerPriceType: models.TriggerPriceTypeMark,
		Price:            decimal.RequireFromString("50000"),
		PriceType:        models.ExecutionPriceTypeLimit,
	}
	stopLoss := &models.TpSlTriggerParam{
		TriggerPrice:     decimal.RequireFromString("40000"),
		TriggerPriceType: models.TriggerPriceTypeMark,
		Price:            decimal.RequireFromString("39000"),
		PriceType:        models.ExecutionPriceTypeLimit,
	}

	params := createOrderObjectParams{
		Market:                   suite.market,
		Account:                  suite.account,
		SyntheticAmount:          decimal.RequireFromString("0.00100000"),
		Price:                    decimal.RequireFromString("43445.11680000"),
		Side:                     models.OrderSideBuy,
		Type:                     models.OrderTypeLimit,
		StarknetDomain:           suite.starknetDomain,
		ExpireTime:               expiryTime,
		PostOnly:                 false,
		PreviousOrderExternalID:  nil,
		OrderExternalID:          nil,
		TimeInForce:              models.TimeInForceGTT,
		SelfTradeProtectionLevel: models.SelfTradeProtectionClient,
		Nonce:                    &suite.nonce,
		BuilderFee:               nil,
		BuilderID:                nil,
		TpSlType:                 &tpSlType,
		TakeProfit:               takeProfit,
		StopLoss:                 stopLoss,
	}

	order, err := createOrderObject(params)
	suite.Require().NoError(err)
	suite.Require().NotNil(order)

	suite.NotNil(order.TakeProfit, "take profit should be set")
	suite.NotNil(order.StopLoss, "stop loss should be set")
	suite.Equal(&tpSlType, order.TpSlType, "TPSL type should be set")

	tp := order.TakeProfit
	suite.Equal("49000", tp.TriggerPrice)
	suite.Equal("50000", tp.Price)
	suite.NotEmpty(tp.Settlement.Signature.R)
	suite.NotEmpty(tp.Settlement.Signature.S)

	sl := order.StopLoss
	suite.Equal("40000", sl.TriggerPrice)
	suite.Equal("39000", sl.Price)
	suite.NotEmpty(sl.Settlement.Signature.R)
	suite.NotEmpty(sl.Settlement.Signature.S)

	suite.NotEqual(tp.Settlement.Signature.R, sl.Settlement.Signature.R, "take profit and stop loss should have different signatures")
}

func (suite *OrdersTestSuite) TestCreateSellOrderWithTPSL() {
	expiryTime := suite.frozenTime.Add(14 * 24 * time.Hour)

	tpSlType := models.TpSlTypeOrder
	takeProfit := &models.TpSlTriggerParam{
		TriggerPrice:     decimal.RequireFromString("50000"),
		TriggerPriceType: models.TriggerPriceTypeMark,
		Price:            decimal.RequireFromString("51000"),
		PriceType:        models.ExecutionPriceTypeLimit,
	}

	params := createOrderObjectParams{
		Market:                   suite.market,
		Account:                  suite.account,
		SyntheticAmount:          decimal.RequireFromString("0.00100000"),
		Price:                    decimal.RequireFromString("45000"),
		Side:                     models.OrderSideSell,
		Type:                     models.OrderTypeLimit,
		StarknetDomain:           suite.starknetDomain,
		ExpireTime:               expiryTime,
		PostOnly:                 false,
		PreviousOrderExternalID:  nil,
		OrderExternalID:          nil,
		TimeInForce:              models.TimeInForceGTT,
		SelfTradeProtectionLevel: models.SelfTradeProtectionClient,
		Nonce:                    &suite.nonce,
		BuilderFee:               nil,
		BuilderID:                nil,
		TpSlType:                 &tpSlType,
		TakeProfit:               takeProfit,
		StopLoss:                 nil,
	}

	order, err := createOrderObject(params)
	suite.Require().NoError(err)
	suite.Require().NotNil(order)

	suite.NotNil(order.TakeProfit, "take profit should be set for SELL order")
	tp := order.TakeProfit
	suite.Equal("50000", tp.TriggerPrice)
	suite.Equal("51000", tp.Price)
	suite.NotEmpty(tp.Settlement.Signature.R)
	suite.NotEmpty(tp.Settlement.Signature.S)
}

func (suite *OrdersTestSuite) TestTPSLFailsIfHashingFails() {
	expiryTime := suite.frozenTime.Add(14 * 24 * time.Hour)

	params := createOrderObjectParams{
		Market:                   suite.market,
		Account:                  suite.account,
		SyntheticAmount:          decimal.RequireFromString("0.00100000"),
		Price:                    decimal.RequireFromString("43445.11680000"),
		Side:                     models.OrderSideBuy,
		Type:                     models.OrderTypeLimit,
		StarknetDomain:           suite.starknetDomain,
		ExpireTime:               expiryTime,
		PostOnly:                 false,
		PreviousOrderExternalID:  nil,
		OrderExternalID:          nil,
		TimeInForce:              models.TimeInForceGTT,
		SelfTradeProtectionLevel: models.SelfTradeProtectionClient,
		Nonce:                    &suite.nonce,
		BuilderFee:               nil,
		BuilderID:                nil,
		TpSlType:                 nil,
		TakeProfit:               nil,
		StopLoss:                 nil,
	}

	order, err := createOrderObject(params)
	suite.Require().NoError(err)
	suite.Require().NotNil(order)
	suite.Nil(order.TakeProfit)
	suite.Nil(order.StopLoss)
}

func TestOrdersTestSuite(t *testing.T) {
	suite.Run(t, new(OrdersTestSuite))
}
