package main

import (
	"fmt"

	"github.com/pedroegsilva/gofindrules/finder"
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

	stringExtractors := []finder.StringExtractor{gfte}
	finder, err := finder.NewRuleFinderWithRules(stringExtractors, rules)
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

	res, err := finder.ExtractTagsObject(someObject, nil)
	if err != nil {
		panic(err)
	}
	for _, r := range res {
		fmt.Println(r)
	}
}
