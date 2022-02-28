package tagger

import (
	"fmt"
	"testing"

	"github.com/pedroegsilva/gotagthem/dsl"
	"github.com/stretchr/testify/assert"
)

func TestNewTagger(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		stringTaggers  []StringTagger
		intTaggers     []IntTagger
		floatTaggers   []FloatTagger
		expectedTagger *Tagger
		message        string
	}{
		{
			stringTaggers: []StringTagger{},
			intTaggers:    []IntTagger{},
			floatTaggers:  []FloatTagger{},
			expectedTagger: &Tagger{
				stringTaggers:               []StringTagger{},
				intTaggers:                  []IntTagger{},
				floatTaggers:                []FloatTagger{},
				expressionWrapperByExprName: make(map[string][]ExpressionWrapper),
				fields:                      make(map[string]struct{}),
				tags:                        make(map[string]struct{}),
			},
			message: "empty tagger",
		},
		{
			stringTaggers: []StringTagger{&emptyStrTagger{}},
			intTaggers:    []IntTagger{&emptyIntTagger{}},
			floatTaggers:  []FloatTagger{&emptyFloatTagger{}},
			expectedTagger: &Tagger{
				stringTaggers:               []StringTagger{&emptyStrTagger{}},
				intTaggers:                  []IntTagger{&emptyIntTagger{}},
				floatTaggers:                []FloatTagger{&emptyFloatTagger{}},
				expressionWrapperByExprName: make(map[string][]ExpressionWrapper),
				fields:                      make(map[string]struct{}),
				tags:                        make(map[string]struct{}),
			},
			message: "tagger with empty taggers",
		},
	}

	for _, tc := range tests {
		tagger := NewTagger(tc.stringTaggers, tc.intTaggers, tc.floatTaggers)
		assert.Equal(tc.expectedTagger, tagger, tc.message)
	}
}

func TestNewTaggerWithRules(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		rulesByName    map[string][]string
		expectedTagger *Tagger
		expectedErr    error
		message        string
	}{
		{
			rulesByName: map[string][]string{
				"rule1": {
					`"tag1"`,
					`"tag2:field1"`,
				},
				"rule2": {
					`"tag3:field2.field3"`,
					`"tag4"`,
				},
			},
			expectedTagger: &Tagger{
				expressionWrapperByExprName: map[string][]ExpressionWrapper{
					"rule1": {
						{
							ExpressionString: `"tag1"`,
							Expression: &dsl.Expression{
								Type: dsl.UNIT_EXPR,
								Tag:  dsl.TagInfo{Name: "tag1", FieldPath: ""},
							},
						},
						{
							ExpressionString: `"tag2:field1"`,
							Expression: &dsl.Expression{
								Type: dsl.UNIT_EXPR,
								Tag:  dsl.TagInfo{Name: "tag2", FieldPath: "field1"},
							},
						},
					},
					"rule2": {
						{
							ExpressionString: `"tag3:field2.field3"`,
							Expression: &dsl.Expression{
								Type: dsl.UNIT_EXPR,
								Tag:  dsl.TagInfo{Name: "tag3", FieldPath: "field2.field3"},
							},
						},
						{
							ExpressionString: `"tag4"`,
							Expression: &dsl.Expression{
								Type: dsl.UNIT_EXPR,
								Tag:  dsl.TagInfo{Name: "tag4", FieldPath: ""},
							},
						},
					},
				},
				fields: map[string]struct{}{
					"field1":        {},
					"field2.field3": {},
				},
				tags: map[string]struct{}{
					"tag1": {},
					"tag2": {},
					"tag3": {},
					"tag4": {},
				},
			},
			expectedErr: nil,
			message:     "new tagger with valid rules",
		},
		{
			rulesByName: map[string][]string{
				"rule1": {`"tag1`},
			},
			expectedTagger: &Tagger{
				expressionWrapperByExprName: map[string][]ExpressionWrapper{},
				fields:                      map[string]struct{}{},
				tags:                        map[string]struct{}{},
			},
			expectedErr: fmt.Errorf("fail to scan tag: expected ':' but found EOF"),
			message:     "new tagger with invalid rules",
		},
	}

	for _, tc := range tests {
		tagger, err := NewTaggerWithRules(nil, nil, nil, tc.rulesByName)
		assert.Equal(tc.expectedErr, err, tc.message)
		assert.Equal(tc.expectedTagger, tagger, tc.message)
	}
}

func TestAddRule(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		ruleName       string
		expressions    []string
		expectedTagger *Tagger
		expectedErr    error
		message        string
	}{
		{
			ruleName: "rule1",
			expressions: []string{
				`"tag1"`,
				`"tag2:field1"`,
			},
			expectedTagger: &Tagger{
				expressionWrapperByExprName: map[string][]ExpressionWrapper{
					"rule1": {
						{
							ExpressionString: `"tag1"`,
							Expression: &dsl.Expression{
								Type: dsl.UNIT_EXPR,
								Tag:  dsl.TagInfo{Name: "tag1", FieldPath: ""},
							},
						},
						{
							ExpressionString: `"tag2:field1"`,
							Expression: &dsl.Expression{
								Type: dsl.UNIT_EXPR,
								Tag:  dsl.TagInfo{Name: "tag2", FieldPath: "field1"},
							},
						},
					},
				},
				fields: map[string]struct{}{
					"field1": {},
				},
				tags: map[string]struct{}{
					"tag1": {},
					"tag2": {},
				},
			},
			expectedErr: nil,
			message:     "add valid expressions",
		},
		{
			ruleName: "rule1",
			expressions: []string{
				`"tag1`,
			},
			expectedTagger: &Tagger{
				expressionWrapperByExprName: map[string][]ExpressionWrapper{},
				fields:                      map[string]struct{}{},
				tags:                        map[string]struct{}{},
			},
			expectedErr: fmt.Errorf("fail to scan tag: expected ':' but found EOF"),
			message:     "add invalid expression",
		},
	}

	for _, tc := range tests {
		tagger := NewTagger(nil, nil, nil)
		err := tagger.AddRule(tc.ruleName, tc.expressions)
		assert.Equal(tc.expectedErr, err, tc.message)
		assert.Equal(tc.expectedTagger, tagger, tc.message)
	}
}

func TestAddRules(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		rulesByName    map[string][]string
		expectedTagger *Tagger
		expectedErr    error
		message        string
	}{
		{
			rulesByName: map[string][]string{
				"rule1": {
					`"tag1"`,
					`"tag2:field1"`,
				},
				"rule2": {
					`"tag3:field2.field3"`,
					`"tag4"`,
				},
			},
			expectedTagger: &Tagger{
				expressionWrapperByExprName: map[string][]ExpressionWrapper{
					"rule1": {
						{
							ExpressionString: `"tag1"`,
							Expression: &dsl.Expression{
								Type: dsl.UNIT_EXPR,
								Tag:  dsl.TagInfo{Name: "tag1", FieldPath: ""},
							},
						},
						{
							ExpressionString: `"tag2:field1"`,
							Expression: &dsl.Expression{
								Type: dsl.UNIT_EXPR,
								Tag:  dsl.TagInfo{Name: "tag2", FieldPath: "field1"},
							},
						},
					},
					"rule2": {
						{
							ExpressionString: `"tag3:field2.field3"`,
							Expression: &dsl.Expression{
								Type: dsl.UNIT_EXPR,
								Tag:  dsl.TagInfo{Name: "tag3", FieldPath: "field2.field3"},
							},
						},
						{
							ExpressionString: `"tag4"`,
							Expression: &dsl.Expression{
								Type: dsl.UNIT_EXPR,
								Tag:  dsl.TagInfo{Name: "tag4", FieldPath: ""},
							},
						},
					},
				},
				fields: map[string]struct{}{
					"field1":        {},
					"field2.field3": {},
				},
				tags: map[string]struct{}{
					"tag1": {},
					"tag2": {},
					"tag3": {},
					"tag4": {},
				},
			},
			message: "add rules with valid rules",
		},
		{
			rulesByName: map[string][]string{
				"rule1": {`"tag1`},
			},
			expectedTagger: &Tagger{
				expressionWrapperByExprName: map[string][]ExpressionWrapper{},
				fields:                      map[string]struct{}{},
				tags:                        map[string]struct{}{},
			},
			expectedErr: fmt.Errorf("fail to scan tag: expected ':' but found EOF"),
			message:     "add rules with invalid rule",
		},
	}

	for _, tc := range tests {
		tagger := NewTagger(nil, nil, nil)
		err := tagger.AddRules(tc.rulesByName)
		assert.Equal(tc.expectedErr, err, tc.message)
		assert.Equal(tc.expectedTagger, tagger, tc.message)
	}
}

func TestTagJson(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		rawJsonStr         string
		expectedFieldsInfo FieldsInfo
		expectedErr        error
		message            string
	}{
		{
			rawJsonStr: `{"strField": "some string", "intField": 42, "floatField": 42.42}`,
			expectedFieldsInfo: FieldsInfo{
				&FieldInfo{
					Name: "strField",
					Taggers: map[string]TaggerInfo{
						"emptyStrTagger": {
							Tags:    []string{"strTag"},
							RunData: nil,
						},
					},
				},
				// At the moment golang json unmarshal, when provided a interface{} as the target object
				// consider all numbers as float 64.
				&FieldInfo{
					Name: "intField",
					Taggers: map[string]TaggerInfo{
						"emptyFloatTagger": {
							Tags:    []string{"floatTag"},
							RunData: nil,
						},
					},
				},
				&FieldInfo{
					Name: "floatField",
					Taggers: map[string]TaggerInfo{
						"emptyFloatTagger": {
							Tags:    []string{"floatTag"},
							RunData: nil,
						},
					},
				},
			},
			expectedErr: nil,
			message:     "tag json",
		},
	}

	for _, tc := range tests {
		tagger := NewTagger(
			[]StringTagger{&emptyStrTagger{}},
			[]IntTagger{&emptyIntTagger{}},
			[]FloatTagger{&emptyFloatTagger{}},
		)
		fieldsInfo, err := tagger.TagJson(tc.rawJsonStr, nil, nil)
		assert.Equal(tc.expectedErr, err, tc.message+" expected error")
		assert.Equal(len(tc.expectedFieldsInfo), len(fieldsInfo), tc.message+" result length")
		eqCount := 0
		for _, info := range fieldsInfo {
			for _, expInfo := range tc.expectedFieldsInfo {
				if expInfo.Name == info.Name {
					res := assert.Equal(expInfo, info, tc.message+" expected field element")
					if res {
						eqCount++
					}
				}
			}
		}
		assert.Equal(len(tc.expectedFieldsInfo), eqCount, tc.message+" expected equal elements found")
	}

}

func TestTagObject(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		object             interface{}
		expectedFieldsInfo FieldsInfo
		expectedErr        error
		message            string
	}{
		{
			object: `some random string`,
			expectedFieldsInfo: FieldsInfo{
				&FieldInfo{
					Name: "",
					Taggers: map[string]TaggerInfo{
						"emptyStrTagger": {
							Tags:    []string{"strTag"},
							RunData: nil,
						},
					},
				},
			},
			expectedErr: nil,
			message:     "tag object raw string",
		},
		{
			object: []string{`some random string`, `some random string`},
			expectedFieldsInfo: FieldsInfo{
				&FieldInfo{
					Name: "index(0)",
					Taggers: map[string]TaggerInfo{
						"emptyStrTagger": {
							Tags:    []string{"strTag"},
							RunData: nil,
						},
					},
				},
				&FieldInfo{
					Name: "index(1)",
					Taggers: map[string]TaggerInfo{
						"emptyStrTagger": {
							Tags:    []string{"strTag"},
							RunData: nil,
						},
					},
				},
			},
			expectedErr: nil,
			message:     "tag object array of string",
		},
		{
			object: struct {
				StrField   string
				StrArray   []string
				AnotherObj struct {
					Field1        int
					Field2        float32
					internalField int
				}
				internalStr string
				internalArr []string
				internalObj struct {
					Field3 float64
				}
			}{
				StrField: "some random string",
				StrArray: []string{
					"some random string 1",
					"some random string 2",
				},
				AnotherObj: struct {
					Field1        int
					Field2        float32
					internalField int
				}{
					Field1:        42,
					Field2:        42.42,
					internalField: 0,
				},
				internalStr: "some internal value",
				internalArr: []string{"some internal value 0"},
				internalObj: struct {
					Field3 float64
				}{
					Field3: 0.0,
				},
			},
			expectedFieldsInfo: FieldsInfo{
				&FieldInfo{
					Name: "StrField",
					Taggers: map[string]TaggerInfo{
						"emptyStrTagger": {
							Tags:    []string{"strTag"},
							RunData: nil,
						},
					},
				},
				&FieldInfo{
					Name: "StrArray.index(0)",
					Taggers: map[string]TaggerInfo{
						"emptyStrTagger": {
							Tags:    []string{"strTag"},
							RunData: nil,
						},
					},
				},
				&FieldInfo{
					Name: "StrArray.index(1)",
					Taggers: map[string]TaggerInfo{
						"emptyStrTagger": {
							Tags:    []string{"strTag"},
							RunData: nil,
						},
					},
				},
				&FieldInfo{
					Name: "AnotherObj.Field1",
					Taggers: map[string]TaggerInfo{
						"emptyIntTagger": {
							Tags:    []string{"intTag"},
							RunData: nil,
						},
					},
				},
				&FieldInfo{
					Name: "AnotherObj.Field2",
					Taggers: map[string]TaggerInfo{
						"emptyFloatTagger": {
							Tags:    []string{"floatTag"},
							RunData: nil,
						},
					},
				},
			},
			expectedErr: nil,
			message:     "tag object struct with internal fields",
		},
	}

	for _, tc := range tests {
		tagger := NewTagger(
			[]StringTagger{&emptyStrTagger{}},
			[]IntTagger{&emptyIntTagger{}},
			[]FloatTagger{&emptyFloatTagger{}},
		)
		fieldsInfo, err := tagger.TagObject(tc.object, nil, nil)
		assert.Equal(tc.expectedErr, err, tc.message+" expected error")
		assert.Equal(len(tc.expectedFieldsInfo), len(fieldsInfo), tc.message+" result length")
		eqCount := 0
		for _, info := range fieldsInfo {
			for _, expInfo := range tc.expectedFieldsInfo {
				if expInfo.Name == info.Name {
					res := assert.Equal(expInfo, info, tc.message+" expected field element")
					if res {
						eqCount++
					}
				}
			}
		}
		assert.Equal(len(tc.expectedFieldsInfo), eqCount, tc.message+" expected equal elements found")
	}

}

func TestTagText(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		text                              string
		expecterExtractorInfoByTaggerName map[string]TaggerInfo
		expectedErr                       error
		message                           string
	}{
		{
			text: `some random string`,
			expecterExtractorInfoByTaggerName: map[string]TaggerInfo{
				"emptyStrTagger": {
					Tags:    []string{"strTag"},
					RunData: nil,
				},
			},
			expectedErr: nil,
			message:     "tag text",
		},
	}
	for _, tc := range tests {
		tagger := NewTagger(
			[]StringTagger{&emptyStrTagger{}},
			nil,
			nil,
		)
		extractorInfoByTaggerName, err := tagger.TagText(tc.text)
		assert.Equal(tc.expectedErr, err, tc.message+" expected error")
		assert.Equal(tc.expecterExtractorInfoByTaggerName, extractorInfoByTaggerName, tc.message+" result")
	}
}

func TestEvaluateRules(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		rulesByName               map[string][]string
		fieldsByTag               map[string][]string
		expectedExpressionsByRule map[string][]string
		expectedErr               error
		message                   string
	}{
		{
			rulesByName: map[string][]string{
				"rule1": {
					`"tag1" and "tag2"`,
				},
				"rule2": {
					`"tag3:field1" and "tag4:field2.innerfield1"`,
				},
				"rule3": {
					`"tag5:field3"`,
				},
				"unmatched rule 1": {
					`"tag4:field1"`,
				},
			},
			fieldsByTag: map[string][]string{
				"tag1": {"randomFiled"},
				"tag2": {"randomFiled2"},
				"tag3": {"field1"},
				"tag4": {"field2.innerfield1"},
				"tag5": {"field3.innerfield2"},
			},
			expectedExpressionsByRule: map[string][]string{
				"rule1": {
					`"tag1" and "tag2"`,
				},
				"rule2": {
					`"tag3:field1" and "tag4:field2.innerfield1"`,
				},
				"rule3": {
					`"tag5:field3"`,
				},
			},
			expectedErr: nil,
			message:     "evaluate rules",
		},
	}
	for _, tc := range tests {
		tagger, _ := NewTaggerWithRules(
			nil,
			nil,
			nil,
			tc.rulesByName,
		)
		extractorInfoByTaggerName, err := tagger.EvaluateRules(tc.fieldsByTag)
		assert.Equal(tc.expectedErr, err, tc.message+" expected error")
		assert.Equal(tc.expectedExpressionsByRule, extractorInfoByTaggerName, tc.message+" result")
	}
}

func TestGetFieldsByTag(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		fieldsInfo          FieldsInfo
		expectedFieldsByTag map[string][]string
		message             string
	}{
		{
			fieldsInfo: FieldsInfo{
				&FieldInfo{
					Name: "StrField",
					Taggers: map[string]TaggerInfo{
						"emptyStrTagger": {
							Tags:    []string{"strTag"},
							RunData: nil,
						},
					},
				},
				&FieldInfo{
					Name: "StrArray.index(0)",
					Taggers: map[string]TaggerInfo{
						"emptyStrTagger": {
							Tags:    []string{"strTag"},
							RunData: nil,
						},
					},
				},
				&FieldInfo{
					Name: "StrArray.index(1)",
					Taggers: map[string]TaggerInfo{
						"emptyStrTagger": {
							Tags:    []string{"strTag"},
							RunData: nil,
						},
					},
				},
				&FieldInfo{
					Name: "AnotherObj.Field1",
					Taggers: map[string]TaggerInfo{
						"emptyIntTagger": {
							Tags:    []string{"intTag"},
							RunData: nil,
						},
					},
				},
				&FieldInfo{
					Name: "AnotherObj.Field2",
					Taggers: map[string]TaggerInfo{
						"emptyFloatTagger": {
							Tags:    []string{"floatTag"},
							RunData: nil,
						},
					},
				},
			},

			expectedFieldsByTag: map[string][]string{
				"strTag":   {"StrField", "StrArray.index(0)", "StrArray.index(1)"},
				"intTag":   {"AnotherObj.Field1"},
				"floatTag": {"AnotherObj.Field2"},
			},

			message: "get fields by tag",
		},
	}
	for _, tc := range tests {
		fieldsByTag := tc.fieldsInfo.GetFieldsByTag()
		assert.Equal(tc.expectedFieldsByTag, fieldsByTag, tc.message+" result")
	}
}

type emptyStrTagger struct{}

func (est *emptyStrTagger) IsValid(data string) bool {
	return true
}

func (est *emptyStrTagger) GetTags(data string) (tags []string, runData interface{}, err error) {
	tags = append(tags, "strTag")
	return
}

func (est *emptyStrTagger) GetName() string {
	return "emptyStrTagger"
}

type emptyIntTagger struct{}

func (est *emptyIntTagger) IsValid(data int64) bool {
	return true
}

func (est *emptyIntTagger) GetTags(data int64) (tags []string, runData interface{}, err error) {
	tags = append(tags, "intTag")
	return
}

func (est *emptyIntTagger) GetName() string {
	return "emptyIntTagger"
}

type emptyFloatTagger struct{}

func (est *emptyFloatTagger) IsValid(data float64) bool {
	return true
}

func (est *emptyFloatTagger) GetTags(data float64) (tags []string, runData interface{}, err error) {
	tags = append(tags, "floatTag")
	return
}

func (est *emptyFloatTagger) GetName() string {
	return "emptyFloatTagger"
}
