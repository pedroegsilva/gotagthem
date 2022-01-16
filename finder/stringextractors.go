package finder

type StringExtractor interface {
	// BuildEngine receive the unique terms that need
	// to be searched to create the engine support structures
	IsValid(data string) bool
	// FindRegexes receive the text and searchs for the feeded
	// regexes
	ExtractTags(data string) (tags []string, err error)
}
