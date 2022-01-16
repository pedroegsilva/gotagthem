package finder

type RuleFinder struct {
	stringExtractors []StringExtractor
	intExtractors    []IntExtractor
	floatExtractors  []FloatExtractor
	rulesByName      map[string][]string
}

type FoundTags struct {
	TagsByExtractorName map[string][]string
	FieldType           string
}

func NewRuleFinder(
	stringExtractors []StringExtractor,
	intExtractors []IntExtractor,
	floatExtractors []FloatExtractor,
	rulesByName map[string][]string,

) *RuleFinder {
	return &RuleFinder{
		stringExtractors: stringExtractors,
		intExtractors:    intExtractors,
		floatExtractors:  floatExtractors,
		rulesByName:      rulesByName,
	}
}

func NewRuleFinderWithRules(
	stringExtractors []StringExtractor,
	intExtractors []IntExtractor,
	floatExtractors []FloatExtractor,
) *RuleFinder {
	return &RuleFinder{
		stringExtractors: stringExtractors,
		intExtractors:    intExtractors,
		floatExtractors:  floatExtractors,
	}
}

func (rf *RuleFinder) AddRule(ruleName string, expressions []string) (err error) {
	return
}

func (rf *RuleFinder) AddRules(rulesByName map[string][]string) (err error) {
	return
}

func (rf *RuleFinder) ExtractTagsObject(
	obj interface{},
	fieldsPath []string,
) (foundTagsByFieldName map[string]FoundTags, err error) {
	return
}

func (rf *RuleFinder) ExtractTagsJson(
	jsonData string,
	fieldsPath []string,
) (foundTagsByFieldName map[string]FoundTags, err error) {
	return
}

func (rf *RuleFinder) ExtractTagsText(
	data string,
	fieldsPath []string,
) (foundTagsByFieldName map[string]FoundTags, err error) {
	return
}

func (rf *RuleFinder) SolveRules(
	foundTagsByFieldName map[string]FoundTags,
	expressionsByRuleName map[string][]string,
) (matchedExpressionsByRuleName map[string][]string, err error) {
	return
}

func (rf *RuleFinder) ProcessObject(
	obj interface{},
	expressionsByRuleName map[string][]string,
) (matchedExpressionsByRuleName map[string][]string, foundTagsByFieldName map[string]FoundTags, err error) {

	return
}

func (rf *RuleFinder) ProcessJson(
	jsonData string,
	expressionsByRuleName map[string][]string,
) (matchedExpressionsByRuleName map[string][]string, foundTagsByFieldName map[string]FoundTags, err error) {

	return
}

func (rf *RuleFinder) ProcessText(
	data string,
	expressionsByRuleName map[string][]string,
) (matchedExpressionsByRuleName map[string][]string, foundTagsByFieldName map[string]FoundTags, err error) {

	return
}
