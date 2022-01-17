package finder

type FieldInfo struct {
	Name       string
	Extractors map[string]ExtractorInfo
}

type ExtractorInfo struct {
	Tags    []string
	RunData interface{}
}
