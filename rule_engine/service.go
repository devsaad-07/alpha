package rule_engine

import (
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"os"
	"path/filepath"
)

var knowledgeLibrary = *ast.NewKnowledgeLibrary()

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

func buildRuleEngine(ruleType string) {
	ruleBuilder := builder.NewRuleBuilder(&knowledgeLibrary)

	// Read rule from file and build rules
	dir, err := os.Getwd()
	path := "/rule_engine/rules/" + ruleType + ".grl"
	path = filepath.Join(dir, path)
	ruleFile := pkg.NewFileResource(path)
	err = ruleBuilder.BuildRuleFromResource("Rules", "0.0.1", ruleFile)
	if err != nil {
		panic(err)
	}
}

func (svc *RuleEngineSvc) Execute(ruleConf RuleConfig) error {
	// get KnowledgeBase instance to execute particular rule
	knowledgeBase := knowledgeLibrary.NewKnowledgeBaseInstance("Rules", "0.0.1")

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
