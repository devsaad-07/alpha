package communication

import (
	"alpha/rule_engine"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Alpha struct {
	rule_engine_map map[string]*rule_engine.RuleEngineSvc
}

var alpha Alpha

func init() {
	ruleEngineMap := map[string]*rule_engine.RuleEngineSvc{
		"fee": rule_engine.NewRuleEngineSvc("fee"),
	}
	alpha = Alpha{
		rule_engine_map: ruleEngineMap,
	}
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func HandleFee(c *gin.Context) {
	ruleEngineSvc := alpha.rule_engine_map["fee"]
	feeService := rule_engine.NewFeeService(ruleEngineSvc)

	var feeCard rule_engine.FeeRequest
	if err := c.ShouldBindJSON(&feeCard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	feeCardResponse := feeService.GetFeeForUser(feeCard)
	c.JSON(http.StatusOK, gin.H{
		"payload": feeCardResponse,
	})
}
