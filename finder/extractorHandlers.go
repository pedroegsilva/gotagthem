package finder

func (rf *RuleFinder) handleFloatExtractors(
	data float64,
) (extractorInfoByExtractorName map[string]ExtractorInfo, err error) {
	extractorInfoByExtractorName = make(map[string]ExtractorInfo)
	for _, extractor := range rf.floatExtractors {
		if extractor.IsValid(data) {
			tags, runData, err := extractor.ExtractTags(data)
			if err != nil {
				return nil, err
			}
			extractorInfoByExtractorName[extractor.GetName()] = ExtractorInfo{
				Tags:    tags,
				RunData: runData,
			}
		}
	}
	return
}

func (rf *RuleFinder) handleIntExtractors(
	data int64,
) (extractorInfoByExtractorName map[string]ExtractorInfo, err error) {
	extractorInfoByExtractorName = make(map[string]ExtractorInfo)

	for _, extractor := range rf.intExtractors {
		if extractor.IsValid(data) {
			tags, runData, err := extractor.ExtractTags(data)
			if err != nil {
				return nil, err
			}
			extractorInfoByExtractorName[extractor.GetName()] = ExtractorInfo{
				Tags:    tags,
				RunData: runData,
			}
		}
	}
	return
}

func (rf *RuleFinder) handleStringExtractors(
	data string,
) (extractorInfoByExtractorName map[string]ExtractorInfo, err error) {
	extractorInfoByExtractorName = make(map[string]ExtractorInfo)

	for _, extractor := range rf.stringExtractors {
		if extractor.IsValid(data) {
			tags, runData, err := extractor.ExtractTags(data)
			if err != nil {
				return nil, err
			}
			extractorInfoByExtractorName[extractor.GetName()] = ExtractorInfo{
				Tags:    tags,
				RunData: runData,
			}
		}
	}
	return
}
