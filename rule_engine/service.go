package rule_engine

import (
	"alpha/db"
	"encoding/json"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

var knowledgeLibrary = *ast.NewKnowledgeLibrary()
var ruleBuilder *builder.RuleBuilder

const KNOWLEDGE_BASE_VERSION = "0.0.0"

// Rule input object
type RuleInput interface {
	DataKey() string
}

// Rule output object
type RuleOutput interface {
	DataKey() string
}

// configs associated with each rule
type RuleConfig interface {
	RuleName() string
	RuleInput() RuleInput
	RuleOutput() RuleOutput
}

type RuleEngineSvc struct {
}

func NewRuleEngineSvc(ruleType string) *RuleEngineSvc {
	// you could add your cloud provider here instead of keeping rule file in your code.
	buildRuleEngine(ruleType)
	return &RuleEngineSvc{}
}

func AddNewRuleWithRuleType(rule interface{}, ruleType string, name string) (err error) {
	ruleBytes, err := json.Marshal(rule)
	if err != nil {
		return
	}
	//rs, _ := pkg.ParseJSONRule(ruleBytes)
	//fmt.Printf("rs: %v", rs)
	underlying := pkg.NewBytesResource(ruleBytes)
	gruleJson := pkg.NewJSONResourceFromResource(underlying)
	err = ruleBuilder.BuildRuleFromResources(ruleType, KNOWLEDGE_BASE_VERSION, []pkg.Resource{gruleJson})
	if err != nil {
		return
	}
	dbRule := db.Rules{
		Type:     ruleType,
		Rule:     string(ruleBytes),
		IsActive: true,
		Name:     name,
	}
	err = db.SaveRule(dbRule)
	return
}

func InjectRulesInEngine(ruleType string) (err error) {
	rules := db.GetAllRules(ruleType)
	grulJSONArray := []pkg.Resource{}

	for _, item := range rules {
		ruleBytes := []byte(item.Rule)
		underlying := pkg.NewBytesResource(ruleBytes)
		gruleJson := pkg.NewJSONResourceFromResource(underlying)
		grulJSONArray = append(grulJSONArray, gruleJson)
	}
	err = ruleBuilder.BuildRuleFromResources(ruleType, KNOWLEDGE_BASE_VERSION, grulJSONArray)
	return
}

func RemoveFromKnowledgeBase(ruleName string, ruleType string) {
	knowledgeLibrary.RemoveRuleEntry(ruleName, ruleType, KNOWLEDGE_BASE_VERSION)
}

func buildRuleEngine(ruleType string) (err error) {
	ruleBuilder = builder.NewRuleBuilder(&knowledgeLibrary)
	InjectRulesInEngine(ruleType)
	return
}

func (svc *RuleEngineSvc) Execute(ruleType string, ruleConf RuleConfig) error {
	// get KnowledgeBase instance to execute particular rule
	knowledgeBase := knowledgeLibrary.NewKnowledgeBaseInstance(ruleType, KNOWLEDGE_BASE_VERSION)

	dataCtx := ast.NewDataContext()
	// add input data context
	err := dataCtx.Add(ruleConf.RuleInput().DataKey(), ruleConf.RuleInput())
	if err != nil {
		return err
	}

	// add output data context
	err = dataCtx.Add(ruleConf.RuleOutput().DataKey(), ruleConf.RuleOutput())
	if err != nil {
		return err
	}

	// create rule engine and execute on provided data and knowledge base
	ruleEngine := engine.NewGruleEngine()
	err = ruleEngine.Execute(dataCtx, knowledgeBase)
	if err != nil {
		return err
	}
	return nil
}
