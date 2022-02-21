package tagger

import (
	"encoding/json"
	"strings"

	"github.com/pedroegsilva/gotagthem/dsl"
)

type FieldInfo struct {
	Name    string
	Taggers map[string]TaggerInfo
}

type TaggerInfo struct {
	Tags    []string
	RunData interface{}
}

type FieldsInfo []*FieldInfo

type ExpressionWrapper struct {
	ExpressionString string
	Expression       *dsl.Expression
}

type Tagger struct {
	stringTaggers               []StringTagger
	intTaggers                  []IntTagger
	floatTaggers                []FloatTagger
	expressionWrapperByExprName map[string][]ExpressionWrapper
	fields                      map[string]struct{}
	tags                        map[string]struct{}
}

func NewTagger(
	stringTaggers []StringTagger,
	intTaggers []IntTagger,
	floatTaggers []FloatTagger,
) *Tagger {
	return &Tagger{
		stringTaggers:               stringTaggers,
		intTaggers:                  intTaggers,
		floatTaggers:                floatTaggers,
		expressionWrapperByExprName: make(map[string][]ExpressionWrapper),
		fields:                      make(map[string]struct{}),
		tags:                        make(map[string]struct{}),
	}
}

func NewTaggerWithRules(
	stringTaggers []StringTagger,
	intTaggers []IntTagger,
	floatTaggers []FloatTagger,
	rulesByName map[string][]string,
) (tagger *Tagger, err error) {
	tagger = NewTagger(stringTaggers, intTaggers, floatTaggers)
	err = tagger.AddRules(rulesByName)
	return
}

func (rf *Tagger) AddRule(ruleName string, expressions []string) error {
	for _, rawExpr := range expressions {
		p := dsl.NewParser(strings.NewReader(rawExpr))
		exp, err := p.Parse()
		if err != nil {
			return err
		}
		expWrapper := ExpressionWrapper{
			ExpressionString: rawExpr,
			Expression:       exp,
		}
		rf.expressionWrapperByExprName[ruleName] = append(rf.expressionWrapperByExprName[ruleName], expWrapper)
		for _, tag := range p.GetTags() {
			rf.tags[tag] = struct{}{}
		}
		for _, field := range p.GetFields() {
			rf.fields[field] = struct{}{}
		}
	}
	return nil
}

func (rf *Tagger) AddRules(rulesByName map[string][]string) error {
	for key, exprs := range rulesByName {
		err := rf.AddRule(key, exprs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rf *Tagger) GetFieldNames() (fields []string) {
	for field := range rf.fields {
		fields = append(fields, field)
	}

	return
}

func (rf *Tagger) TagJson(
	data string,
	includePaths []string,
	excludePaths []string,
) (fieldsInfo FieldsInfo, err error) {
	var genericObj interface{}
	err = json.Unmarshal([]byte(data), &genericObj)
	if err != nil {
		return
	}
	err = rf.extractTagsObject(genericObj, "", &fieldsInfo, includePaths, excludePaths)
	return
}

func (rf *Tagger) TagObject(
	data interface{},
	includePaths []string,
	excludePaths []string,
) (fieldsInfo FieldsInfo, err error) {
	err = rf.extractTagsObject(data, "", &fieldsInfo, includePaths, excludePaths)
	return
}

func (rf *Tagger) TagText(
	data string,
) (extractorInfoByTaggerName map[string]TaggerInfo, err error) {
	fieldsInfo, err := rf.TagObject(data, nil, nil)
	return fieldsInfo[0].Taggers, err
}

func (rf *Tagger) EvaluateRules(
	fieldsByTag map[string][]string,
) (expressionsByRule map[string][]string, err error) {
	expressionsByRule = make(map[string][]string)
	for name, exprWrappers := range rf.expressionWrapperByExprName {
		for _, ew := range exprWrappers {
			eval, err := ew.Expression.Solve(fieldsByTag)
			if err != nil {
				return nil, err
			}
			if eval {
				expressionsByRule[name] = append(expressionsByRule[name], ew.ExpressionString)
			}
		}
	}
	return
}

func (rf *Tagger) ProcessJson(
	rawJson string,
	includePaths []string,
	excludePaths []string,
) (expressionsByRule map[string][]string, err error) {
	fieldsInfo, err := rf.TagJson(rawJson, includePaths, excludePaths)
	if err != nil {
		return nil, err
	}

	return rf.EvaluateRules(fieldsInfo.GetFieldsByTag())
}

func (rf *Tagger) ProcessObject(
	obj interface{},
	includePaths []string,
	excludePaths []string,
) (expressionsByRule map[string][]string, err error) {
	fieldsInfo, err := rf.TagObject(obj, includePaths, excludePaths)
	if err != nil {
		return nil, err
	}

	return rf.EvaluateRules(fieldsInfo.GetFieldsByTag())
}

func (rf *Tagger) ProcessText(
	data string,
) (expressionsByRule map[string][]string, err error) {
	extractorInfoByTaggerName, err := rf.TagText(data)
	if err != nil {
		return nil, err
	}

	fieldsByTag := make(map[string][]string)
	for _, extractorInfo := range extractorInfoByTaggerName {
		for _, tag := range extractorInfo.Tags {
			fieldsByTag[tag] = nil
		}
	}
	return rf.EvaluateRules(fieldsByTag)
}

func (fieldsInfo FieldsInfo) GetFieldsByTag() (fieldsByTag map[string][]string) {
	fieldsByTag = make(map[string][]string)
	for _, fieldInfo := range fieldsInfo {
		for _, extractorInfo := range fieldInfo.Taggers {
			for _, tag := range extractorInfo.Tags {
				fieldsByTag[tag] = append(fieldsByTag[tag], fieldInfo.Name)
			}
		}
	}
	return
}
