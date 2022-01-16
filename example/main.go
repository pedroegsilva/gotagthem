package main

import (
	"fmt"
	"strings"

	"github.com/pedroegsilva/gofindrules/dsl"
	"github.com/pedroegsilva/gofindrules/finder"
)

func main() {
	p := dsl.NewParser(strings.NewReader(`"foo" and ("bar" or "some:field1")`))
	expression, err := p.Parse()

	fmt.Println("Error: ", err)
	fmt.Println("expression: ", expression.PrettyFormat())

	eval, err := expression.Solve(map[string][]string{"foo": {"fieldn"}, "some": {"field1"}})
	fmt.Println("Error solve: ", err)
	fmt.Println("eval: ", eval)
	fmt.Println("-------------------------------------------------------------------")

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
		"rule2": {`"tag3" or "tag4"`},
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

	res, err := finder.ProcessObject("some pretty text with string1")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
