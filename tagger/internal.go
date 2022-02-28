package tagger

import (
	"fmt"
	"reflect"
	"strings"
)

// setFieldInfos adds to the fieldsInfo the information extracteds by all the taggers
func (rf *Tagger) setFieldInfos(
	data interface{},
	fieldName string,
	fieldsInfo *FieldsInfo,
	includePaths []string,
	excludePaths []string,
) (err error) {
	t := reflect.TypeOf(data)

	val := reflect.ValueOf(data)

	if !val.CanInterface() {
		return
	}

	switch val.Kind() {
	case reflect.String:
		if !isValidateFieldPath(fieldName, includePaths, excludePaths) {
			return
		}

		extractorInfoByTaggerName, err := rf.handleStringTaggers(val.String())
		if err != nil {
			return err
		}
		fieldInfo := &FieldInfo{Name: fieldName, Taggers: extractorInfoByTaggerName}
		*fieldsInfo = append(*fieldsInfo, fieldInfo)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !isValidateFieldPath(fieldName, includePaths, excludePaths) {
			return
		}

		extractorInfoByTaggerName, err := rf.handleIntTaggers(val.Int())
		if err != nil {
			return err
		}
		fieldInfo := &FieldInfo{Name: fieldName, Taggers: extractorInfoByTaggerName}
		*fieldsInfo = append(*fieldsInfo, fieldInfo)

	case reflect.Float32, reflect.Float64:
		if !isValidateFieldPath(fieldName, includePaths, excludePaths) {
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
			if !val.Field(i).CanInterface() {
				continue
			}
			err := rf.setFieldInfos(val.Field(i).Interface(), fn, fieldsInfo, includePaths, excludePaths)
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
			if !v.CanInterface() {
				continue
			}
			err := rf.setFieldInfos(v.Interface(), fn, fieldsInfo, includePaths, excludePaths)
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
			if !val.Index(i).CanInterface() {
				continue
			}
			err := rf.setFieldInfos(val.Index(i).Interface(), fn, fieldsInfo, includePaths, excludePaths)
			if err != nil {
				return err
			}
		}
	}

	return
}

// handleFloatTaggers hander of taggers of the type float
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

// handleFloatTaggers hander of taggers of the type int
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

// handleFloatTaggers hander of taggers of the type string
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

// isValidateFieldPath returns true if the field path is valid for tagging
func isValidateFieldPath(fieldPath string, includePaths []string, excludePaths []string) bool {
	if len(excludePaths) > 0 {
		for _, excP := range excludePaths {
			if strings.HasPrefix(fieldPath, excP) {
				return false
			}
		}
	}

	if len(includePaths) > 0 {
		for _, incP := range includePaths {
			if strings.HasPrefix(fieldPath, incP) {
				return true
			}
		}
		return false
	}

	return true
}
