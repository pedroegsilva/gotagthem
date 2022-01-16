package finder

type FloatExtractor interface {
	IsValid(data float64) bool
	ExtractTags(data float64) (tags []string, err error)
	GetName() string
}
