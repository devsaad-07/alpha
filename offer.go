package main

type UserOfferContext struct {
	UserOfferInput  *UserOfferInput
	UserOfferOutput *UserOfferOutput
}

func (uoc *UserOfferContext) RuleName() string {
	return "user_offers"
}

func (uoc *UserOfferContext) RuleInput() RuleInput {
	return uoc.UserOfferInput
}

func (uoc *UserOfferContext) RuleOutput() RuleOutput {
	return uoc.UserOfferOutput
}

// User data attributes
type UserOfferInput struct {
	Name              string  `json:"name"`
	Username          string  `json:"username"`
	Email             string  `json:"email"`
	Age               int     `json:"age"`
	Gender            string  `json:"gender"`
	TotalOrders       int     `json:"total_orders"`
	AverageOrderValue float64 `json:"average_order_value"`
}

func (u *UserOfferInput) DataKey() string {
	return "InputData"
}

// Offer output object
type UserOfferOutput struct {
	IsOfferApplicable bool `json:"is_offer_applicable"`
}

func (u *UserOfferOutput) DataKey() string {
	return "OutputData"
}

func NewUserOfferContext() *UserOfferContext {
	return &UserOfferContext{
		UserOfferInput:  &UserOfferInput{},
		UserOfferOutput: &UserOfferOutput{},
	}
}
