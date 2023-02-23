package rule_engine

import (
	feediscount "alpha/db"
	"time"

	"github.com/hyperjumptech/grule-rule-engine/logger"
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
	UserMetrics           feediscount.UserMetrics
	AssetType             string    `json:"assetType"`
	AssetPair             string    `json:"assetPair"`
	OrderSource           string    `json:"orderSource"` // client calling this ex: CSK, CS Pro
	UserId                int64     `json:"userId"`
	UserDOB               time.Time `json:"userDOB"`
	TradeRequestAmountINR float64   `json:"tradeRequestAmountINR"`
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
	UserMetrics           feediscount.UserMetrics
	AssetType             string    `json:"assetType"`
	AssetPair             string    `json:"assetPair"`
	OrderSource           string    `json:"orderSource"` // client calling this ex: CSK, CS Pro
	UserId                int64     `json:"userId"`
	TradeVolumeInr        float64   `json:"tradeVolumeInr"`
	TradeRequestAmountINR float64   `json:"tradeRequestAmountINR"`
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
		UserMetrics:           feeRequest.UserMetrics,
		AssetType:             feeRequest.AssetType,
		AssetPair:             feeRequest.AssetPair,
		OrderSource:           feeRequest.OrderSource,
		UserId:                feeRequest.UserId,
		TradeRequestAmountINR: feeRequest.TradeRequestAmountINR,
		UserDOB:               feeRequest.UserDOB,
		TradeExchange:         feeRequest.TradeExchange,
		TradeType:             feeRequest.TradeType,
		OrderType:             feeRequest.OrderType,
		HopType:               feeRequest.HopType,
		Occasion:              feeRequest.Occasion,
	}
	err := svc.ruleEngineSvc.Execute("fee", feeCard)
	if err != nil {
		logger.Log.Error("get user fee rule engine failed", err)
	}

	feeResponse = FeeResponse{
		Fee:  feeCard.FeeOutput.Fee,
		Type: feeCard.FeeOutput.Type,
	}
	return
}
