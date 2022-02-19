package finder

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pedroegsilva/gotagthem/dsl"
)

type ExpressionWrapper struct {
	ExpressionString string
	Expression       *dsl.Expression
}

type RuleFinder struct {
	stringExtractors            []StringExtractor
	intExtractors               []IntExtractor
	floatExtractors             []FloatExtractor
	expressionWrapperByExprName map[string][]ExpressionWrapper
	fields                      map[string]struct{}
	tags                        map[string]struct{}
}

func NewRuleFinder(
	stringExtractors []StringExtractor,
	intExtractors []IntExtractor,
	floatExtractors []FloatExtractor,
) *RuleFinder {
	return &RuleFinder{
		stringExtractors:            stringExtractors,
		intExtractors:               intExtractors,
		floatExtractors:             floatExtractors,
		expressionWrapperByExprName: make(map[string][]ExpressionWrapper),
		fields:                      make(map[string]struct{}),
		tags:                        make(map[string]struct{}),
	}
}

func NewRuleFinderWithRules(
	stringExtractors []StringExtractor,
	intExtractors []IntExtractor,
	floatExtractors []FloatExtractor,
	rulesByName map[string][]string,
) (finder *RuleFinder, err error) {
	finder = NewRuleFinder(stringExtractors, intExtractors, floatExtractors)
	err = finder.AddRules(rulesByName)
	return
}

func (rf *RuleFinder) AddRule(ruleName string, expressions []string) error {
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
			if tag != "" {
				rf.tags[tag] = struct{}{}
			}
		}
		for _, field := range p.GetFields() {
			if field != "" {
				rf.fields[field] = struct{}{}
			}
		}
	}
	return nil
}

func (rf *RuleFinder) AddRules(rulesByName map[string][]string) error {
	for key, exprs := range rulesByName {
		err := rf.AddRule(key, exprs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rf *RuleFinder) GetFieldNames() (fields []string) {
	for field := range rf.fields {
		fields = append(fields, field)
	}

	return
}

func (rf *RuleFinder) ExtractTagsObject(
	data interface{},
	includePaths []string,
	excludePaths []string,
) (fieldsInfo FieldsInfo, err error) {
	err = rf.extractTagsObject(data, "", &fieldsInfo, includePaths, excludePaths)
	return
}

func (rf *RuleFinder) extractTagsObject(
	data interface{},
	fieldName string,
	fieldsInfo *FieldsInfo,
	includePaths []string,
	excludePaths []string,
) (err error) {
	t := reflect.TypeOf(data)
	val := reflect.ValueOf(data)

	switch val.Kind() {
	case reflect.String:
		if !validatePath(fieldName, includePaths, excludePaths) {
			return
		}

		extractorInfoByExtractorName, err := rf.handleStringExtractors(val.String())
		if err != nil {
			return err
		}
		fieldInfo := &FieldInfo{Name: fieldName, Extractors: extractorInfoByExtractorName}
		*fieldsInfo = append(*fieldsInfo, fieldInfo)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !validatePath(fieldName, includePaths, excludePaths) {
			return
		}

		extractorInfoByExtractorName, err := rf.handleIntExtractors(val.Int())
		if err != nil {
			return err
		}
		fieldInfo := &FieldInfo{Name: fieldName, Extractors: extractorInfoByExtractorName}
		*fieldsInfo = append(*fieldsInfo, fieldInfo)

	case reflect.Float32, reflect.Float64:
		if !validatePath(fieldName, includePaths, excludePaths) {
			return
		}

		extractorInfoByExtractorName, err := rf.handleFloatExtractors(val.Float())
		if err != nil {
			return err
		}

		fieldInfo := &FieldInfo{Name: fieldName, Extractors: extractorInfoByExtractorName}
		*fieldsInfo = append(*fieldsInfo, fieldInfo)

	case reflect.Struct:
		numField := t.NumField()

		for i := 0; i < numField; i++ {
			structField := t.Field(i)
			fn := structField.Name
			if fieldName != "" {
				fn = fieldName + "." + fn
			}
			err := rf.extractTagsObject(val.Field(i).Interface(), fn, fieldsInfo, includePaths, excludePaths)
			if err != nil {
				return err
			}
		}

	case reflect.Map:
		iter := val.MapRange()
		for iter.Next() {
			k := iter.Key()
			if k.Type().Kind() != reflect.String {
				break
			}

			v := iter.Value()
			fn := k.String()
			if fieldName != "" {
				fn = fieldName + "." + fn
			}
			err := rf.extractTagsObject(v.Interface(), fn, fieldsInfo, includePaths, excludePaths)
			if err != nil {
				return err
			}
		}

	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			fn := fmt.Sprintf("index(%d)", i)
			if fieldName != "" {
				fn = fieldName + "." + fn
			}
			err := rf.extractTagsObject(val.Index(i).Interface(), fn, fieldsInfo, includePaths, excludePaths)
			if err != nil {
				return err
			}
		}
	}

	return
}

func (rf *RuleFinder) ExtractTagsText(
	data string,
) (extractorInfoByExtractorName map[string]ExtractorInfo, err error) {
	fieldsInfo, err := rf.ExtractTagsObject(data, nil, nil)
	return fieldsInfo[0].Extractors, err
}

func (rf *RuleFinder) SolveRules(
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

func (rf *RuleFinder) ProcessObject(
	obj interface{},
	includePaths []string,
	excludePaths []string,
) (expressionsByRule map[string][]string, err error) {
	fieldsInfo, err := rf.ExtractTagsObject(obj, includePaths, excludePaths)
	if err != nil {
		return nil, err
	}

	return rf.SolveRules(fieldsInfo.GetFieldsByTag())
}

func (rf *RuleFinder) ProcessText(
	data string,
) (expressionsByRule map[string][]string, err error) {
	extractorInfoByExtractorName, err := rf.ExtractTagsText(data)
	if err != nil {
		return nil, err
	}

	fieldsByTag := make(map[string][]string)
	for _, extractorInfo := range extractorInfoByExtractorName {
		for _, tag := range extractorInfo.Tags {
			fieldsByTag[tag] = nil
		}
	}
	return rf.SolveRules(fieldsByTag)
}

func validatePath(path string, includePaths []string, excludePaths []string) bool {
	if len(excludePaths) > 0 {
		for _, excP := range excludePaths {
			if strings.HasPrefix(path, excP) {
				return false
			}
		}
	}

	if len(includePaths) > 0 {
		for _, incP := range includePaths {
			if strings.HasPrefix(path, incP) {
				return true
			}
		}
		return false
	}

	return true
}
