package finder

import (
	"reflect"
	"strings"

	"github.com/pedroegsilva/gofindrules/dsl"
)

type RuleFinder struct {
	stringExtractors       []StringExtractor
	intExtractors          []IntExtractor
	floatExtractors        []FloatExtractor
	expressionsByName      map[string][]*dsl.Expression
	rawExpressionByPointer map[*dsl.Expression]string
}

type TagsByExtractorName map[string][]string

func NewRuleFinder(
	stringExtractors []StringExtractor,
) *RuleFinder {
	return &RuleFinder{
		stringExtractors:       stringExtractors,
		expressionsByName:      make(map[string][]*dsl.Expression),
		rawExpressionByPointer: make(map[*dsl.Expression]string),
	}
}

func NewRuleFinderWithRules(
	stringExtractors []StringExtractor,
	rulesByName map[string][]string,
) (finder *RuleFinder, err error) {
	finder = &RuleFinder{
		stringExtractors:       stringExtractors,
		expressionsByName:      make(map[string][]*dsl.Expression),
		rawExpressionByPointer: make(map[*dsl.Expression]string),
	}
	err = finder.AddRules(rulesByName)
	return
}

func (rf *RuleFinder) AddRule(ruleName string, expressions []string) error {
	for _, expr := range expressions {
		p := dsl.NewParser(strings.NewReader(expr))
		exp, err := p.Parse()
		if err != nil {
			return err
		}
		rf.expressionsByName[ruleName] = append(rf.expressionsByName[ruleName], exp)
		rf.rawExpressionByPointer[exp] = expr
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

func (rf *RuleFinder) ExtractTagsObject(
	data interface{},
	fieldsPath []string,
) (tagsByExtractorByField map[string]TagsByExtractorName, err error) {
	tagsByExtractorByField = make(map[string]TagsByExtractorName)
	t := reflect.TypeOf(data)
	val := reflect.ValueOf(data)

	switch val.Kind() {
	case reflect.String:
		tagsByExtractor, err := rf.handleStringExtractors(val.String())
		if err != nil {
			return nil, err
		}
		tagsByExtractorByField[""] = tagsByExtractor

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		tagsByExtractor, err := rf.handleIntExtractors(val.Int())
		if err != nil {
			return nil, err
		}
		tagsByExtractorByField[""] = tagsByExtractor

	case reflect.Float32, reflect.Float64:
		tagsByExtractor, err := rf.handleFloatExtractors(val.Float())
		if err != nil {
			return nil, err
		}
		tagsByExtractorByField[""] = tagsByExtractor

	case reflect.Struct:
		numField := t.NumField()

		for i := 0; i < numField; i++ {
			structField := t.Field(i)

			tagsByExtractorByField, err := rf.ExtractTagsObject(val.Field(i).Interface(), fieldsPath)
			if err != nil {
				return nil, err
			}

			for fieldName, tagsByExtractor := range tagsByExtractorByField {
				key := structField.Name
				if fieldName != "" {
					key += "." + fieldName
				}
				tagsByExtractorByField[key] = tagsByExtractor
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

			tagsByExtractorByField, err := rf.ExtractTagsObject(v.Interface(), fieldsPath)
			if err != nil {
				return nil, err
			}

			for fieldName, tagsByExtractor := range tagsByExtractorByField {
				key := k.String()
				if fieldName != "" {
					key += "." + fieldName
				}
				tagsByExtractorByField[key] = tagsByExtractor
			}
		}
	}

	return
}

func (rf *RuleFinder) ExtractTagsText(
	data string,
	fieldsPath []string,
) (tagsByExtractor map[string][]string, err error) {
	tagsByExtractorByField, err := rf.ExtractTagsObject(data, fieldsPath)
	return tagsByExtractorByField[""], err
}

func (rf *RuleFinder) SolveRules(
	fieldsByTag map[string][]string,
) (expressionsByRule map[string][]string, err error) {
	expressionsByRule = make(map[string][]string)
	for name, exprs := range rf.expressionsByName {
		for _, exp := range exprs {
			eval, err := exp.Solve(fieldsByTag)
			if err != nil {
				return nil, err
			}
			if eval {
				rawExp := rf.rawExpressionByPointer[exp]
				expressionsByRule[name] = append(expressionsByRule[name], rawExp)
			}
		}
	}
	return
}

func (rf *RuleFinder) ProcessObject(
	obj interface{},
) (expressionsByRule map[string][]string, err error) {
	tagsByExtractorByField, err := rf.ExtractTagsObject(obj, nil)
	if err != nil {
		return nil, err
	}

	fieldsByTag := make(map[string][]string)
	for field, tagsByExtractor := range tagsByExtractorByField {
		for _, tags := range tagsByExtractor {
			for _, tag := range tags {
				fieldsByTag[tag] = append(fieldsByTag[tag], field)
			}
		}
	}

	return rf.SolveRules(fieldsByTag)
}

func (rf *RuleFinder) ProcessText(
	data string,
) (expressionsByRule map[string][]string, err error) {
	tagsByExtractor, err := rf.ExtractTagsText(data, nil)
	if err != nil {
		return nil, err
	}

	fieldsByTag := make(map[string][]string)
	for _, tags := range tagsByExtractor {
		for _, tag := range tags {
			fieldsByTag[tag] = nil
		}
	}
	return rf.SolveRules(fieldsByTag)
}
