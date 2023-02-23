package rule_engine

import (
	"encoding/json"
	"fmt"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strconv"
)

var knowledgeLibrary = *ast.NewKnowledgeLibrary()
var ruleBuilder *builder.RuleBuilder

const KNOWLEDGE_BASE_NAME = "Test Knowledge Base"
const KNOWLEDGE_BASE_VERSION = "0.0."

var versionCount = 0

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

func Test(rule interface{}) (err error) {
	rule1, err := json.Marshal(rule)
	if err != nil {
		return
	}
	underlying := pkg.NewBytesResource([]byte(rule1))
	gruleJson := pkg.NewJSONResourceFromResource(underlying)
	err = ruleBuilder.BuildRuleFromResource("fee", "0.0.2", gruleJson)
	return
}

func buildRuleEngine(ruleType string) (err error) {
	ruleBuilder = builder.NewRuleBuilder(&knowledgeLibrary)

	return reloadRuleEngine(ruleType, 0)
}

func reloadRuleEngine(ruleType string, version int) (err error) {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	path := "/rule_engine/rules/" + ruleType + ".grl"
	path = filepath.Join(dir, path)
	ruleFile := pkg.NewFileResource(path)
	versionString := KNOWLEDGE_BASE_VERSION + strconv.Itoa(versionCount)
	err = ruleBuilder.BuildRuleFromResource(KNOWLEDGE_BASE_NAME, versionString, ruleFile)
	versionCount = versionCount + version
	if err != nil {
		return
	}
	return
}

func (svc *RuleEngineSvc) Execute(ruleConf RuleConfig) error {
	// get KnowledgeBase instance to execute particular rule
	knowledgeBase := knowledgeLibrary.NewKnowledgeBaseInstance(KNOWLEDGE_BASE_NAME, KNOWLEDGE_BASE_VERSION)

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

func ParseRule(rule *GruleJSON) (r string, err error) {
	// Convert the "when" and "then" interfaces to JSON strings
	whenJSON, err := json.Marshal(rule.When)
	if err != nil {
		return "", fmt.Errorf("failed to marshal 'when' condition: %v", err)
	}
	thenJSON, err := json.Marshal(rule.Then)
	if err != nil {
		return "", fmt.Errorf("failed to marshal 'then' action: %v", err)
	}

	// Create the Grule rule string
	ruleString := fmt.Sprintf("rule \"%s\"\nwhen\n%s\nthen\n%s\nend", rule.Name, whenJSON, thenJSON)

	return ruleString, nil
}

func GetKnowledgeBase() *ast.KnowledgeBase {
	return knowledgeLibrary.GetKnowledgeBase(KNOWLEDGE_BASE_NAME, KNOWLEDGE_BASE_VERSION)
}

func AddRule(ruleType string, ruleJSON string) (err error) {
	rule, err := pkg.ParseJSONRule([]byte(ruleJSON))
	if err != nil {
		logrus.Errorf("Error: %s", err.Error())
		return
	}
	err = addRuleToGRL(ruleType, rule)
	if err != nil {
		return
	}
	err = reloadRuleEngine("fee", 1)
	return
}

func addRuleToGRL(ruleType string, rule string) (err error) {
	dir, err := os.Getwd()
	path := "/rule_engine/rules/" + ruleType + ".grl"
	path = filepath.Join(dir, path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	n, err := f.Write([]byte(rule))
	if err != nil {
		return
	}
	fmt.Println("wrote %d bytes", n)
	return
}

//func AddRuleToKnowledgeBase(ruleType string) {
//	rule := `
//
//rule AgeNameCheck "test" {
//    when
//      Pogo.GetStringLength("9999") > 0  &&
//      Pogo.Result == ""
//    then
//      Pogo.Result = "String len above 0";
//}
//`
//	dir, err := os.Getwd()
//	path := "/rule_engine/rules/" + ruleType + ".grl"
//	path = filepath.Join(dir, path)
//	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
//	if err != nil {
//		panic(err)
//	}
//	defer f.Close()
//
//	n, err := f.Write([]byte(rule))
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("wrote %d bytes", n)
//	return
//}

//func UpdateRuleInKnowledgeBase(ruleType string, ruleName string, updatedRule string) {
//	kb := knowledgeLibrary.GetKnowledgeBase(KNOWLEDGE_BASE_NAME, KNOWLEDGE_BASE_VERSION)
//	wm := kb.WorkingMemory
//	logrus.Infof("working memory %v", wm)
//	kb.RemoveRuleEntry(ruleName)
//	AddRuleToKnowledgeBase(ruleType)
//}

//func DeleteRuleInKnowledgeBase(ruleType string, ruleName string) {
//	kb := knowledgeLibrary.GetKnowledgeBase(KNOWLEDGE_BASE_NAME, KNOWLEDGE_BASE_VERSION)
//	kb.RemoveRuleEntry(ruleName)
//
//	dir, err := os.Getwd()
//	if err != nil {
//		fmt.Errorf("Error, %v", err)
//	}
//	path := "/rule_engine/rules/" + ruleType + ".grl"
//	path = filepath.Join(dir, path)
//	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
//	err = knowledgeLibrary.StoreKnowledgeBaseToWriter(f, KNOWLEDGE_BASE_NAME, KNOWLEDGE_BASE_VERSION)
//	if err != nil {
//		fmt.Errorf("Error, %v", err)
//	}
//}
