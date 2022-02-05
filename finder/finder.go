package finder

import (
	"reflect"
	"strings"

	"github.com/pedroegsilva/gotagthem/dsl"
)

type RuleFinder struct {
	stringExtractors       []StringExtractor
	intExtractors          []IntExtractor
	floatExtractors        []FloatExtractor
	solverOrderByExprName  map[string][]*dsl.SolverOrder
	rawExpressionByPointer map[*dsl.SolverOrder]string
}

func NewRuleFinder(
	stringExtractors []StringExtractor,
	intExtractors []IntExtractor,
	floatExtractors []FloatExtractor,
) *RuleFinder {
	return &RuleFinder{
		stringExtractors:       stringExtractors,
		intExtractors:          intExtractors,
		floatExtractors:        floatExtractors,
		solverOrderByExprName:  make(map[string][]*dsl.SolverOrder),
		rawExpressionByPointer: make(map[*dsl.SolverOrder]string),
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
		so := exp.CreateSolverOrder()
		rf.solverOrderByExprName[ruleName] = append(rf.solverOrderByExprName[ruleName], &so)
		rf.rawExpressionByPointer[&so] = rawExpr
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
) (fieldsInfo []*FieldInfo, err error) {
	t := reflect.TypeOf(data)
	val := reflect.ValueOf(data)

	switch val.Kind() {
	case reflect.String:
		extractorInfoByExtractorName, err := rf.handleStringExtractors(val.String())
		if err != nil {
			return nil, err
		}
		fieldInfo := &FieldInfo{Name: "", Extractors: extractorInfoByExtractorName}
		fieldsInfo = append(fieldsInfo, fieldInfo)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		extractorInfoByExtractorName, err := rf.handleIntExtractors(val.Int())
		if err != nil {
			return nil, err
		}
		fieldInfo := &FieldInfo{Name: "", Extractors: extractorInfoByExtractorName}
		fieldsInfo = append(fieldsInfo, fieldInfo)

	case reflect.Float32, reflect.Float64:
		extractorInfoByExtractorName, err := rf.handleFloatExtractors(val.Float())
		if err != nil {
			return nil, err
		}
		fieldInfo := &FieldInfo{Name: "", Extractors: extractorInfoByExtractorName}
		fieldsInfo = append(fieldsInfo, fieldInfo)

	case reflect.Struct:
		numField := t.NumField()

		for i := 0; i < numField; i++ {
			structField := t.Field(i)
			res, err := rf.ExtractTagsObject(val.Field(i).Interface(), fieldsPath)
			if err != nil {
				return nil, err
			}

			for _, fieldInfo := range res {
				newName := structField.Name
				if fieldInfo.Name != "" {
					newName += "." + fieldInfo.Name
				}
				fieldInfo.Name = newName
				fieldsInfo = append(fieldsInfo, fieldInfo)
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
			res, err := rf.ExtractTagsObject(v.Interface(), fieldsPath)
			if err != nil {
				return nil, err
			}

			for _, fieldInfo := range res {
				newName := k.String()
				if fieldInfo.Name != "" {
					newName += "." + fieldInfo.Name
				}
				fieldInfo.Name = newName
				fieldInfo.Name = newName
				fieldsInfo = append(fieldsInfo, fieldInfo)
			}
		}
	}

	return
}

func (rf *RuleFinder) ExtractTagsText(
	data string,
	fieldsPath []string,
) (extractorInfoByExtractorName map[string]ExtractorInfo, err error) {
	fieldsInfo, err := rf.ExtractTagsObject(data, fieldsPath)
	return fieldsInfo[0].Extractors, err
}

func (rf *RuleFinder) SolveRules(
	fieldsByTag map[string][]string,
) (expressionsByRule map[string][]string, err error) {
	expressionsByRule = make(map[string][]string)
	for name, exprs := range rf.solverOrderByExprName {
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
	fieldsInfo, err := rf.ExtractTagsObject(obj, nil)
	if err != nil {
		return nil, err
	}

	fieldsByTag := make(map[string][]string)
	for _, fieldInfo := range fieldsInfo {
		for _, extractorInfo := range fieldInfo.Extractors {
			for _, tag := range extractorInfo.Tags {
				fieldsByTag[tag] = append(fieldsByTag[tag], fieldInfo.Name)
			}
		}
	}

	return rf.SolveRules(fieldsByTag)
}

func (rf *RuleFinder) ProcessText(
	data string,
) (expressionsByRule map[string][]string, err error) {
	extractorInfoByExtractorName, err := rf.ExtractTagsText(data, nil)
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
