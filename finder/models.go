package finder

type FieldInfo struct {
	Name       string
	Extractors map[string]ExtractorInfo
}

type ExtractorInfo struct {
	Tags    []string
	RunData interface{}
}

type FieldsInfo []*FieldInfo

func (fieldsInfo FieldsInfo) GetFieldsByTag() (fieldsByTag map[string][]string) {
	fieldsByTag = make(map[string][]string)
	for _, fieldInfo := range fieldsInfo {
		for _, extractorInfo := range fieldInfo.Extractors {
			for _, tag := range extractorInfo.Tags {
				fieldsByTag[tag] = append(fieldsByTag[tag], fieldInfo.Name)
			}
		}
	}
	return
}
