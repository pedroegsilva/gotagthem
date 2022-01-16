package dsl

import (
	"fmt"
	"strings"
)

// ExprType are Special tokens used to define the expression type
type ExprType int

const (
	UNSET_EXPR ExprType = iota
	AND_EXPR
	OR_EXPR
	NOT_EXPR
	UNIT_EXPR
)

// GetName returns a readable name for the ExprType value
func (exprType ExprType) GetName() string {
	switch exprType {
	case UNSET_EXPR:
		return "UNSET"
	case AND_EXPR:
		return "AND"
	case OR_EXPR:
		return "OR"
	case NOT_EXPR:
		return "NOT"
	case UNIT_EXPR:
		return "UNIT"
	default:
		return "UNEXPECTED"
	}
}

// Expression can be a literal (UNIT) or a function composed by
// one or two other expressions (NOT, AND, OR).
type Literal struct {
	Tag    string
	Fields []FieldInfo
}

type FieldInfo struct {
	Type string
	Path string
}

// Expression can be a literal (UNIT) or a function composed by
// one or two other expressions (NOT, AND, OR).
type Expression struct {
	LExpr   *Expression
	RExpr   *Expression
	Type    ExprType
	Literal Literal

	// used on solveOrder solver
	val bool
}

// PatternResult stores if the patter was matched on
// the text and the positions it was found
type PatternResult struct {
	Val            bool
	SortedMatchPos []int
}

// GetTypeName returns the type of the expression with a readable name
func (exp *Expression) GetTypeName() string {
	return exp.Type.GetName()
}

// Solve solves the expresion recursively. It has the option to use a complete map of
// PatternResult or a incomplete map. If the complete map option is used the map must have
// all the terms needed to solve de expression or it will return an error.
// If the incomplete map is used, missing keys will be considered as a no match on the
// document.
func (exp *Expression) Solve(
	fieldTypeByFieldPathByTag map[string]map[string]string,
) (bool, error) {
	eval, err := exp.solve(fieldTypeByFieldPathByTag)
	return eval, err
}

//solve implements Solve
func (exp *Expression) solve(
	fieldTypeByFieldPathByTag map[string]map[string]string,
) (bool, error) {
	switch exp.Type {
	case UNIT_EXPR:
		if typeByFieldPath, ok := fieldTypeByFieldPathByTag[exp.Literal.Tag]; ok {
			if len(exp.Literal.Fields) == 0 {
				return true, nil
			}

			for _, expField := range exp.Literal.Fields {
				if fieldType, ok := typeByFieldPath[expField.Path]; ok {
					if expField.Type == "" {
						return true, nil
					}

					return fieldType == expField.Type, nil
				}
			}
		}

		return false, nil

	case AND_EXPR:
		if exp.LExpr == nil || exp.RExpr == nil {
			return false, fmt.Errorf("AND statment do not have rigth or left expression: %v", exp)
		}
		lval, err := exp.LExpr.solve(fieldTypeByFieldPathByTag)
		if err != nil {
			return false, err
		}
		rval, err := exp.RExpr.solve(fieldTypeByFieldPathByTag)
		if err != nil {
			return false, err
		}

		return lval && rval, nil
	case OR_EXPR:
		if exp.LExpr == nil || exp.RExpr == nil {
			return false, fmt.Errorf("OR statment do not have rigth or left expression: %v", exp)
		}
		lval, err := exp.LExpr.solve(fieldTypeByFieldPathByTag)
		if err != nil {
			return false, err
		}
		rval, err := exp.RExpr.solve(fieldTypeByFieldPathByTag)
		if err != nil {
			return false, err
		}

		return lval || rval, nil
	case NOT_EXPR:
		if exp.RExpr == nil {
			return false, fmt.Errorf("NOT statement do not have expression: %v", exp)
		}
		rval, err := exp.RExpr.solve(fieldTypeByFieldPathByTag)
		if err != nil {
			return false, err
		}
		return !rval, nil
	default:
		return false, fmt.Errorf("unable to process expression type %d", exp.Type)
	}
}

// PrettyPrint returns the expression formated on a tabbed structure
// Eg: for the expression ("a" and "b") or "c"
//    OR
//        AND
//            a
//            b
//        c
func (exp *Expression) PrettyFormat() string {
	return exp.prettyFormat(0)
}

func (exp *Expression) prettyFormat(lvl int) (pprint string) {
	tabs := "    "
	onLVL := strings.Repeat(tabs, lvl)
	if exp.Type == UNIT_EXPR {
		return fmt.Sprintf("%s%s\n", onLVL, exp.Literal)
	}
	pprint = fmt.Sprintf("%s%s\n", onLVL, exp.GetTypeName())
	if exp.LExpr != nil {
		pprint += exp.LExpr.prettyFormat(lvl + 1)
	}

	if exp.RExpr != nil {
		pprint += exp.RExpr.prettyFormat(lvl + 1)
	}

	return
}

// SolverOrder store the expressions Preorder
type SolverOrder []*Expression

// Solve solves the expresion iteratively. It has the option to use a complete map of
// PatternResult or a incomplete map. If the complete map option is used the map must have
// all the terms needed to solve de expression or it will return an error.
// If the incomplete map is used, missing keys will be considered as a no match on the
// document.
func (so SolverOrder) Solve(fieldTypeByFieldPathByTag map[string]map[string]string) (bool, error) {
	for i := len(so) - 1; i >= 0; i-- {
		exp := so[i]
		if exp == nil {
			return false, fmt.Errorf("malformed solver order - solver order should not have nil values")
		}
	Switch:
		switch exp.Type {
		case UNIT_EXPR:
			exp.val = false
			if typeByFieldPath, ok := fieldTypeByFieldPathByTag[exp.Literal.Tag]; ok {
				if len(exp.Literal.Fields) == 0 {
					exp.val = true
					break Switch
				}

				for _, expField := range exp.Literal.Fields {
					if fieldType, ok := typeByFieldPath[expField.Path]; ok {
						if expField.Type == "" {
							exp.val = true
							break Switch
						}

						exp.val = fieldType == expField.Type
						break Switch
					}
				}
			}
		case AND_EXPR:
			if exp.LExpr == nil || exp.RExpr == nil {
				return false, fmt.Errorf("AND statement do not have right or left expression: %v", exp)
			}
			exp.val = exp.LExpr.val && exp.RExpr.val

		case OR_EXPR:
			if exp.LExpr == nil || exp.RExpr == nil {
				return false, fmt.Errorf("OR statement do not have right or left expression: %v", exp)
			}
			exp.val = exp.LExpr.val || exp.RExpr.val

		case NOT_EXPR:
			if exp.RExpr == nil {
				return false, fmt.Errorf("NOT statement do not have expression: %v", exp)
			}

			exp.val = !exp.RExpr.val

		default:
			return false, fmt.Errorf("unable to process expression type %d", exp.Type)
		}
	}
	return so[0].val, nil
}

// CreateSolverOrder traverses the expression tree in Preorder and
// stores the expressions on SolverOrder
func (exp *Expression) CreateSolverOrder() SolverOrder {
	solverOrder := new(SolverOrder)
	cpExp := exp
	createSolverOrder(cpExp, solverOrder)
	return *solverOrder
}

// createSolverOrder recursion that traverses the expression
// tree in Preorder
func createSolverOrder(exp *Expression, arr *SolverOrder) {
	(*arr) = append((*arr), exp)

	if exp.LExpr != nil {
		createSolverOrder(exp.LExpr, arr)
	}

	if exp.RExpr != nil {
		createSolverOrder(exp.RExpr, arr)
	}
}
