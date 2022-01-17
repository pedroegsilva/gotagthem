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
	}

	gfte, err := finder.NewGoFindThemExtractor(gofindthemRules)
	if err != nil {
		panic(err)
	}
	u, _ := finder.NewUselessIntExtractor()
	stringExtractors := []finder.StringExtractor{gfte}
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
		}
	}{
		Field1: "some pretty text with string1",
		Field2: 42,
		Field3: struct{ SomeField1 string }{
			SomeField1: "some pretty text with string5",
		},
	}

	fieldInfos, err := finder.ExtractTagsObject(someObject, nil)
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
}
