package main

import (
	"fmt"
	"strings"

	"github.com/pedroegsilva/gofindrules/dsl"
)

func main() {
	p := dsl.NewParser(strings.NewReader(`":field2" and not ("bar" or "some:field1[type1]")`))
	expression, err := p.Parse()

	fmt.Println("Error: ", err)
	fmt.Println("expression: ", expression.PrettyFormat())

	eval, err := expression.Solve(map[string]map[string]string{"foo": nil, "some": {"field1": "type"}})
	fmt.Println("Error solve: ", err)
	fmt.Println("eval: ", eval)
}
