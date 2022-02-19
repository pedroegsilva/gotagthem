package main

import (
	"fmt"

	"github.com/pedroegsilva/gotagthem/finder"
)

func main() {
	gofindthemRules := map[string][]string{
		"tag1": {
			`"string1"`,
			`"string2"`,
		},
		"tag2": {
			`"string3"`,
			`"string4"`,
		},
		"tag3": {
			`"string5"`,
			`"string6"`,
		},
		"tag4": {
			`"string7"`,
			`"string8"`,
		},
	}

	rules := map[string][]string{
		"rule1": {`"tag1" or "tag2"`},
		"rule2": {`"tag3:Field3.SomeField1" or "tag4"`},
		"rule3": {`"tag3:Field3" or "tag4"`},
	}

	gfte, err := finder.NewGoFindThemExtractor(gofindthemRules)
	if err != nil {
		panic(err)
	}
	de, err := finder.NewDummyExtractor()
	if err != nil {
		panic(err)
	}
	u, _ := finder.NewUselessIntExtractor()
	stringExtractors := []finder.StringExtractor{gfte, de}
	intExtractors := []finder.IntExtractor{u}
	floatExtractors := []finder.FloatExtractor{}

	finder, err := finder.NewRuleFinderWithRules(stringExtractors, intExtractors, floatExtractors, rules)
	if err != nil {
		panic(err)
	}

	someObject := struct {
		Field1 string
		Field2 int
		Field3 struct {
			SomeField1 string
			SomeField2 []string
		}
	}{
		Field1: "some pretty text with string1",
		Field2: 42,
		Field3: struct {
			SomeField1 string
			SomeField2 []string
		}{
			SomeField1: "some pretty text with string5",
			SomeField2: []string{"some pretty text with string5", "some pretty text with string2", "some pretty text with string3"},
		},
	}

	fmt.Println("finder.GetFieldNames()", finder.GetFieldNames())
	fieldInfos, err := finder.ExtractTagsObject(someObject, finder.GetFieldNames(), nil)
	if err != nil {
		panic(err)
	}
	for _, fieldInfo := range fieldInfos {
		fmt.Println(fieldInfo.Name)
		for extractorName, info := range fieldInfo.Extractors {
			fmt.Println("    ", extractorName)
			fmt.Println("        tags: ", info.Tags)
			fmt.Println("        statistics: ", info.RunData)
		}
	}

	res, err := finder.ProcessObject(someObject, finder.GetFieldNames(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	fmt.Println("-----------------------------")
	arr := []struct {
		FieldN string
		FieldX string
	}{
		{FieldN: "some pretty text with string5"},
		{FieldN: "some pretty text with string2"},
		{FieldN: "some pretty text with string3"},
	}
	fmt.Println("finder.GetFieldNames()", finder.GetFieldNames())
	fieldInfos2, err := finder.ExtractTagsObject(arr, finder.GetFieldNames(), nil)
	if err != nil {
		panic(err)
	}
	for _, fieldInfo := range fieldInfos2 {
		fmt.Println(fieldInfo.Name)
		for extractorName, info := range fieldInfo.Extractors {
			fmt.Println("    ", extractorName)
			fmt.Println("        tags: ", info.Tags)
			fmt.Println("        statistics: ", info.RunData)
		}
	}

	res2, err := finder.ProcessObject(arr, finder.GetFieldNames(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(res2)
}
