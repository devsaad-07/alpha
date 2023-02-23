package rule_engine

import (
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

func AddNewRuleWithRuleType(rule interface{}, ruleType string) (err error) {
	ruleBytes, err := json.Marshal(rule)
	if err != nil {
		return
	}
	underlying := pkg.NewBytesResource(ruleBytes)
	gruleJson := pkg.NewJSONResourceFromResource(underlying)
	err = ruleBuilder.BuildRuleFromResources(ruleType, KNOWLEDGE_BASE_VERSION, []pkg.Resource{gruleJson})
	// TODO: save rule in DB when no error
	return
}

func InjectRulesInEngine(ruleType string) (err error) {
	// TODO: get rules from DB
	return
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

type GruleJSON struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	When        interface{} `json:"when"`
	Then        interface{} `json:"then"`
}
