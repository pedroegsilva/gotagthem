package finder

type IntExtractor interface {
	IsValid(data int64) bool
	ExtractTags(data int64) (tags []string, err error)
	GetName() string
}
