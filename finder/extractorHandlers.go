package finder

func (rf *RuleFinder) handleFloatExtractors(
	data float64,
) (tagsByExtractor map[string][]string, err error) {
	tagsByExtractor = make(map[string][]string)

	for _, extractor := range rf.floatExtractors {
		if extractor.IsValid(data) {
			tags, err := extractor.ExtractTags(data)
			if err != nil {
				return nil, err
			}
			if tags != nil {
				tagsByExtractor[extractor.GetName()] = tags
			}
		}
	}
	return
}

func (rf *RuleFinder) handleIntExtractors(
	data int64,
) (tagsByExtractor map[string][]string, err error) {
	tagsByExtractor = make(map[string][]string)

	for _, extractor := range rf.intExtractors {
		if extractor.IsValid(data) {
			tags, err := extractor.ExtractTags(data)
			if err != nil {
				return nil, err
			}
			if tags != nil {
				tagsByExtractor[extractor.GetName()] = tags
			}
		}
	}
	return
}

func (rf *RuleFinder) handleStringExtractors(
	data string,
) (tagsByExtractor map[string][]string, err error) {
	tagsByExtractor = make(map[string][]string)

	for _, extractor := range rf.stringExtractors {
		if extractor.IsValid(data) {
			tags, err := extractor.ExtractTags(data)
			if err != nil {
				return nil, err
			}
			if tags != nil {
				tagsByExtractor[extractor.GetName()] = tags
			}
		}
	}
	return
}
