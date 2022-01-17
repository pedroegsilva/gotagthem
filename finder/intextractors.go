package finder

type IntExtractor interface {
	IsValid(data int64) bool
	ExtractTags(data int64) (tags []string, runData interface{}, err error)
	GetName() string
}

type UselessIntExtractor struct{}

func NewUselessIntExtractor() (*UselessIntExtractor, error) {
	return &UselessIntExtractor{}, nil
}

func (uie *UselessIntExtractor) IsValid(data int64) bool {
	return data >= 0
}

func (uie *UselessIntExtractor) ExtractTags(data int64) (tags []string, runData interface{}, err error) {
	if data == 42 {
		tags = append(tags, "right")
	} else {
		tags = append(tags, "wrong")
	}

	return tags, nil, nil
}

func (uie *UselessIntExtractor) GetName() string {
	return "uselessInt"
}
