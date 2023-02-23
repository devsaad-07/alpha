package communication

import (
	"alpha/rule_engine"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Alpha struct {
	rule_engine *rule_engine.RuleEngineSvc
}

var alpha Alpha

func init() {
	alpha = Alpha{
		rule_engine: rule_engine.NewRuleEngineSvc(),
	}
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func HandleFee(c *gin.Context) {
	ruleEngineSvc := alpha.rule_engine
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
