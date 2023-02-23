package main

import (
	"fmt"
	"github.com/hyperjumptech/grule-rule-engine/logger"
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

// can be moved to offer directory
type OfferService interface {
	CheckOfferForUser(user User) bool
}

type OfferServiceClient struct {
	ruleEngineSvc *RuleEngineSvc
}

func NewOfferService(ruleEngineSvc *RuleEngineSvc) OfferService {
	return &OfferServiceClient{
		ruleEngineSvc: ruleEngineSvc,
	}
}

func (svc OfferServiceClient) CheckOfferForUser(user User) bool {
	offerCard := NewUserOfferContext()
	offerCard.UserOfferInput = &UserOfferInput{
		Name:              user.Name,
		Username:          user.Username,
		Email:             user.Email,
		Gender:            user.Gender,
		Age:               user.Age,
		TotalOrders:       user.TotalOrders,
		AverageOrderValue: user.AverageOrderValue,
	}

	err := svc.ruleEngineSvc.Execute(offerCard)
	if err != nil {
		logger.Log.Error("get user offer rule engine failed", err)
	}

	return offerCard.UserOfferOutput.IsOfferApplicable
}

func main() {
	ruleEngineSvc := NewRuleEngineSvc()
	offerSvc := NewOfferService(ruleEngineSvc)

	userA := User{
		Name:              "Mohit Khare",
		Username:          "mkfeuhrer",
		Email:             "me@mohitkhare.com",
		Gender:            "Male",
		Age:               99,
		TotalOrders:       50,
		AverageOrderValue: 225,
	}

	fmt.Println("offer validity for user A: ", offerSvc.CheckOfferForUser(userA))

	userB := User{
		Name:              "Pranjal Sharma",
		Username:          "pj",
		Email:             "pj@abc.com",
		Gender:            "Male",
		Age:               25,
		TotalOrders:       10,
		AverageOrderValue: 80,
	}

	fmt.Println("offer validity for user B: ", offerSvc.CheckOfferForUser(userB))
}
