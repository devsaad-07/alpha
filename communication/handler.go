package communication

import (
	feeDiscount "alpha/db"
	"alpha/rule_engine"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Alpha struct {
	ruleEngineMap map[string]*rule_engine.RuleEngineSvc
}

var alpha Alpha

func init() {
	ruleEngineMap := map[string]*rule_engine.RuleEngineSvc{
		"fee": rule_engine.NewRuleEngineSvc("fee"),
	}
	alpha = Alpha{
		ruleEngineMap: ruleEngineMap,
	}
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func HandleFee(c *gin.Context) {
	ruleEngineSvc := alpha.ruleEngineMap["fee"]
	feeService := rule_engine.NewFeeService(ruleEngineSvc)

	var feeCard rule_engine.FeeRequest
	if err := c.ShouldBindJSON(&feeCard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	feeCard.UserMetrics = feeDiscount.GetUser(feeCard.UserId)
	feeCardResponse := feeService.GetFeeForUser(feeCard)
	c.JSON(http.StatusOK, gin.H{
		"payload": feeCardResponse,
	})
}

type RuleRequest struct {
	Type string      `json:"type"`
	Rule interface{} `json:"rule"`
}

func HandleAddRule(c *gin.Context) {
	var ruleRequest RuleRequest
	if err := c.ShouldBindJSON(&ruleRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	err := rule_engine.AddNewRuleWithRuleType(ruleRequest.Rule, ruleRequest.Type)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"payload": "Rule Added",
	})
}
