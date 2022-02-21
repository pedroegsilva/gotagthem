package tagger

import (
	"fmt"
	"reflect"
	"strings"
)

func (rf *Tagger) extractTagsObject(
	data interface{},
	fieldName string,
	fieldsInfo *FieldsInfo,
	includePaths []string,
	excludePaths []string,
) (err error) {
	t := reflect.TypeOf(data)
	val := reflect.ValueOf(data)
	switch val.Kind() {
	case reflect.String:
		if !validatePath(fieldName, includePaths, excludePaths) {
			return
		}

		extractorInfoByTaggerName, err := rf.handleStringTaggers(val.String())
		if err != nil {
			return err
		}
		fieldInfo := &FieldInfo{Name: fieldName, Taggers: extractorInfoByTaggerName}
		*fieldsInfo = append(*fieldsInfo, fieldInfo)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !validatePath(fieldName, includePaths, excludePaths) {
			return
		}

		extractorInfoByTaggerName, err := rf.handleIntTaggers(val.Int())
		if err != nil {
			return err
		}
		fieldInfo := &FieldInfo{Name: fieldName, Taggers: extractorInfoByTaggerName}
		*fieldsInfo = append(*fieldsInfo, fieldInfo)

	case reflect.Float32, reflect.Float64:
		if !validatePath(fieldName, includePaths, excludePaths) {
			return
		}

		extractorInfoByTaggerName, err := rf.handleFloatTaggers(val.Float())
		if err != nil {
			return err
		}

		fieldInfo := &FieldInfo{Name: fieldName, Taggers: extractorInfoByTaggerName}
		*fieldsInfo = append(*fieldsInfo, fieldInfo)

	case reflect.Struct:
		numField := t.NumField()

		for i := 0; i < numField; i++ {
			structField := t.Field(i)
			fn := structField.Name
			if fieldName != "" {
				fn = fieldName + "." + fn
			}
			err := rf.extractTagsObject(val.Field(i).Interface(), fn, fieldsInfo, includePaths, excludePaths)
			if err != nil {
				return err
			}
		}

	case reflect.Map:
		iter := val.MapRange()
		for iter.Next() {
			k := iter.Key()
			if k.Type().Kind() != reflect.String {
				break
			}

			v := iter.Value()
			fn := k.String()
			if fieldName != "" {
				fn = fieldName + "." + fn
			}
			err := rf.extractTagsObject(v.Interface(), fn, fieldsInfo, includePaths, excludePaths)
			if err != nil {
				return err
			}
		}

	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			fn := fmt.Sprintf("index(%d)", i)
			if fieldName != "" {
				fn = fieldName + "." + fn
			}
			err := rf.extractTagsObject(val.Index(i).Interface(), fn, fieldsInfo, includePaths, excludePaths)
			if err != nil {
				return err
			}
		}
	}

	return
}

func (rf *Tagger) handleFloatTaggers(
	data float64,
) (extractorInfoByTaggerName map[string]TaggerInfo, err error) {
	extractorInfoByTaggerName = make(map[string]TaggerInfo)
	for _, extractor := range rf.floatTaggers {
		if extractor.IsValid(data) {
			tags, runData, err := extractor.GetTags(data)
			if err != nil {
				return nil, err
			}
			extractorInfoByTaggerName[extractor.GetName()] = TaggerInfo{
				Tags:    tags,
				RunData: runData,
			}
		}
	}
	return
}

func (rf *Tagger) handleIntTaggers(
	data int64,
) (extractorInfoByTaggerName map[string]TaggerInfo, err error) {
	extractorInfoByTaggerName = make(map[string]TaggerInfo)

	for _, extractor := range rf.intTaggers {
		if extractor.IsValid(data) {
			tags, runData, err := extractor.GetTags(data)
			if err != nil {
				return nil, err
			}
			extractorInfoByTaggerName[extractor.GetName()] = TaggerInfo{
				Tags:    tags,
				RunData: runData,
			}
		}
	}
	return
}

func (rf *Tagger) handleStringTaggers(
	data string,
) (extractorInfoByTaggerName map[string]TaggerInfo, err error) {
	extractorInfoByTaggerName = make(map[string]TaggerInfo)

	for _, extractor := range rf.stringTaggers {
		if extractor.IsValid(data) {
			tags, runData, err := extractor.GetTags(data)
			if err != nil {
				return nil, err
			}
			extractorInfoByTaggerName[extractor.GetName()] = TaggerInfo{
				Tags:    tags,
				RunData: runData,
			}
		}
	}
	return
}

func validatePath(path string, includePaths []string, excludePaths []string) bool {
	if len(excludePaths) > 0 {
		for _, excP := range excludePaths {
			if strings.HasPrefix(path, excP) {
				return false
			}
		}
	}

	if len(includePaths) > 0 {
		for _, incP := range includePaths {
			if strings.HasPrefix(path, incP) {
				return true
			}
		}
		return false
	}

	return true
}
