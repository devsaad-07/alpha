package main

import (
	"fmt"
	"github.com/hyperjumptech/grule-rule-engine/logger"
	"time"
)

// can be part of user serice and a separate directory
type User struct {
	Name              string  `json:"name"`
	Username          string  `json:"username"`
	Email             string  `json:"email"`
	Age               int     `json:"age"`
	Gender            string  `json:"gender"`
	TotalOrders       int     `json:"total_orders"`
	AverageOrderValue float64 `json:"average_order_value"`
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

// can be moved to offer directory
type OfferService interface {
	CheckOfferForUser(user User) bool
}

type FeeService interface {
	GetFeeForUser(feeRequest FeeRequest) FeeResponse
}

type OfferServiceClient struct {
	ruleEngineSvc *RuleEngineSvc
}

type FeeServiceClient struct {
	ruleEngineSvc *RuleEngineSvc
}

//func NewOfferService(ruleEngineSvc *RuleEngineSvc) OfferService {
//	return &OfferServiceClient{
//		ruleEngineSvc: ruleEngineSvc,
//	}
//}

func NewFeeService(ruleEngineSvc *RuleEngineSvc) FeeService {
	return &FeeServiceClient{
		ruleEngineSvc: ruleEngineSvc,
	}
}

//func (svc OfferServiceClient) CheckOfferForUser(user User) bool {
//	offerCard := NewUserOfferContext()
//	offerCard.UserOfferInput = &UserOfferInput{
//		Name:              user.Name,
//		Username:          user.Username,
//		Email:             user.Email,
//		Gender:            user.Gender,
//		Age:               user.Age,
//		TotalOrders:       user.TotalOrders,
//		AverageOrderValue: user.AverageOrderValue,
//	}
//
//	err := svc.ruleEngineSvc.Execute(offerCard)
//	if err != nil {
//		logger.Log.Error("get user offer rule engine failed", err)
//	}
//
//	return offerCard.UserOfferOutput.IsOfferApplicable
//}

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

func main() {
	ruleEngineSvc := NewRuleEngineSvc()
	feeService := NewFeeService(ruleEngineSvc)

	feeCardA := FeeRequest{
		TradeCount: 4,
	}
	feeCardB := FeeRequest{
		TradeCount:            6,
		TradeRequestAmountINR: 999,
	}

	fmt.Println("offer validity for user A: ", feeService.GetFeeForUser(feeCardA))

	fmt.Println("offer validity for user B: ", feeService.GetFeeForUser(feeCardB))
}
