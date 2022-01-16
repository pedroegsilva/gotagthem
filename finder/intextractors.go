package finder

type IntExtractor interface {
	// BuildEngine receive the unique terms that need
	// to be searched to create the engine support structures
	IsValid(data int64) bool
	// FindRegexes receive the text and searchs for the feeded
	// regexes
	ExtractTags(data int64) (tags []string, err error)
	ConvertInt(data int) (newData int64, err error)
}
