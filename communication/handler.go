package communication

import (
	"alpha/db"
	feeDiscount "alpha/db"
	"alpha/rule_engine"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	Name string      `json:"name"`
}

func HandleAddRule(c *gin.Context) {
	var ruleRequest RuleRequest
	if err := c.ShouldBindJSON(&ruleRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	err := rule_engine.AddNewRuleWithRuleType(ruleRequest.Rule, ruleRequest.Type, ruleRequest.Name)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"payload": "Rule Added",
	})
}

func HandleGetAllRules(c *gin.Context) {
	queryMap := c.Request.URL.Query()
	ruleType := queryMap["type"][0]
	rules := db.GetAllRules(ruleType)
	c.JSON(http.StatusOK, gin.H{
		"rules": rules,
	})
}

func DeleteRule(c *gin.Context) {
	queryMap := c.Request.URL.Query()
	id, _ := strconv.Atoi(queryMap["id"][0])
	rule := db.GetRuleById(id)
	// Remove rule from builder
	rule_engine.RemoveFromKnowledgeBase(rule.Name, rule.Type)
	// Mark rule as in_active on builder
	db.UpdateRuleStatus(id, false)
}
