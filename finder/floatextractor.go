package finder

type FloatExtractor interface {
	// BuildEngine receive the unique terms that need
	// to be searched to create the engine support structures
	IsValid(data float64) bool
	// FindRegexes receive the text and searchs for the feeded
	// regexes
	ExtractTags(data float64) (tags []string, err error)
	ConvertFloat32(data float32) (newData float64, err error)
	GetName() string
}
