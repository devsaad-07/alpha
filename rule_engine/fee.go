package rule_engine

import (
	"github.com/hyperjumptech/grule-rule-engine/logger"
	"time"
)

type FeeContext struct {
	FeeInput  *FeeInput
	FeeOutput *FeeOutput
}

func (uoc *FeeContext) RuleName() string {
	return "user_offers"
}

func (uoc *FeeContext) RuleInput() RuleInput {
	return uoc.FeeInput
}

func (uoc *FeeContext) RuleOutput() RuleOutput {
	return uoc.FeeOutput
}

// User data attributes
type FeeInput struct {
	AssetType             string    `json:"assetType"`
	AssetPair             string    `json:"assetPair"`
	OrderSource           string    `json:"orderSource"` // client calling this ex: CSK, CS Pro
	UserId                int64     `json:"userId"`
	TradeVolumeInr        float64   `json:"tradeVolumeInr"`
	TradeCount            int64     `json:"tradeCount"`
	TradeRequestAmountINR float64   `json:"tradeRequestAmountINR"`
	UserCreatedAt         time.Time `json:"userCreatedAt"`
	UserLastLogin         time.Time `json:"userLastLogin"`
	UserLastTradeTime     time.Time `json:"userLastTradeTime"`
	UserDOB               time.Time `json:"userDOB"`
	TradeExchange         string    `json:"tradeExchange"`
	TradeType             string    `json:"tradeType"`
	OrderType             string    `json:"orderType"` // INSTANT, LIMIT, MARKET etc
	HopType               string    `json:"hopType"`
	Occasion              string    `json:"occasion"`
}

func (u *FeeInput) DataKey() string {
	return "InputData"
}

type FeeOutput struct {
	Fee  float64 `json:"fee" default:"0.5"`
	Type string  `json:"type" default:"percent"` // inr, percentage
}

func (u *FeeOutput) DataKey() string {
	return "OutputData"
}

func NewFeeContext() *FeeContext {
	return &FeeContext{
		FeeInput:  &FeeInput{},
		FeeOutput: &FeeOutput{},
	}
}

func NewFeeService(ruleEngineSvc *RuleEngineSvc) FeeService {
	return &FeeServiceClient{
		ruleEngineSvc: ruleEngineSvc,
	}
}

type FeeRequest struct {
	AssetType             string    `json:"assetType"`
	AssetPair             string    `json:"assetPair"`
	OrderSource           string    `json:"orderSource"` // client calling this ex: CSK, CS Pro
	UserId                int64     `json:"userId"`
	TradeVolumeInr        float64   `json:"tradeVolumeInr"`
	TradeCount            int64     `json:"tradeCount"`
	TradeRequestAmountINR float64   `json:"tradeRequestAmountINR"`
	UserCreatedAt         time.Time `json:"userCreatedAt"`
	UserLastLogin         time.Time `json:"userLastLogin"`
	UserLastTradeTime     time.Time `json:"userLastTradeTime"`
	UserDOB               time.Time `json:"userDOB"`
	TradeExchange         string    `json:"tradeExchange"`
	TradeType             string    `json:"tradeType"`
	OrderType             string    `json:"orderType"` // INSTANT, LIMIT, MARKET etc
	HopType               string    `json:"hopType"`
	Occasion              string    `json:"occasion"`
}

type FeeResponse struct {
	Fee  float64 `json:"fee"`
	Type string  `json:"type"` // inr, percentage
}

type FeeService interface {
	GetFeeForUser(feeRequest FeeRequest) FeeResponse
}

type FeeServiceClient struct {
	ruleEngineSvc *RuleEngineSvc
}

func (svc FeeServiceClient) GetFeeForUser(feeRequest FeeRequest) (feeResponse FeeResponse) {
	feeCard := NewFeeContext()
	feeCard.FeeInput = &FeeInput{
		AssetType:             feeRequest.AssetType,
		AssetPair:             feeRequest.AssetPair,
		OrderSource:           feeRequest.OrderSource,
		UserId:                feeRequest.UserId,
		TradeVolumeInr:        feeRequest.TradeVolumeInr,
		TradeCount:            feeRequest.TradeCount,
		TradeRequestAmountINR: feeRequest.TradeRequestAmountINR,
		UserCreatedAt:         feeRequest.UserCreatedAt,
		UserLastLogin:         feeRequest.UserLastLogin,
		UserLastTradeTime:     feeRequest.UserLastTradeTime,
		UserDOB:               feeRequest.UserDOB,
		TradeExchange:         feeRequest.TradeExchange,
		TradeType:             feeRequest.TradeType,
		OrderType:             feeRequest.OrderType,
		HopType:               feeRequest.HopType,
		Occasion:              feeRequest.Occasion,
	}

	err := svc.ruleEngineSvc.Execute(feeCard)
	if err != nil {
		logger.Log.Error("get user fee rule engine failed", err)
	}

	feeResponse = FeeResponse{
		Fee:  feeCard.FeeOutput.Fee,
		Type: feeCard.FeeOutput.Type,
	}
	return
}