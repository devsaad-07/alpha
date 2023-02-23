package main

import "time"

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
