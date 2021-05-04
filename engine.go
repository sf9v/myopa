package myopa

import (
	"context"
	"strings"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/pkg/errors"
)

// Engine is an OPA engine
type Engine struct {
	policyFile string
	module     []byte
}

// New reads a policy file and returns a new engine
func New(policyFile string, module []byte) *Engine {
	return &Engine{
		policyFile: string(policyFile),
		module:     module,
	}
}

// Result is a compilation result
type Result struct {
	Defined bool
	Exprs   []*Expr
}

// Expr is an expression
type Expr struct {
	// Operator is the operator
	Operator Operator
	Left     *Operand
	Right    *Operand
}

type Operator int

const (
	OperatorEq Operator = iota + 1
)

func (op Operator) String() string {
	return [...]string{
		"invalid",
		"equal",
	}[op]
}

func strToOp(s string) Operator {
	switch s {
	case "eq":
		return OperatorEq
	}

	return Operator(0)
}

// Val is an expression value
type Operand struct {
	// T is the operand type
	T OperandType
	// V is the operand value
	V interface{}
}

// OperandType is a value type
type OperandType int

// List of operand types
const (
	OperandTypeConstant OperandType = iota + 1
	OperandTypeIndexField
)

func (vt OperandType) String() string {
	return [...]string{
		"invalid",
		"constant",
		"index-field",
	}[vt]
}

// Compile compiles the query
func (e *Engine) Compile(ctx context.Context, query string, input interface{}, unknowns ...string) (Result, error) {
	rg := rego.New(
		rego.Query(query),
		rego.Module(e.policyFile, string(e.module)),
		rego.Input(input),
		rego.Unknowns(unknowns),
	)

	pqs, err := rg.Partial(ctx)
	if err != nil {
		return Result{}, errors.Wrap(err, "partial eval")
	}

	if len(pqs.Queries) == 0 {
		// always deny
		return Result{}, nil
	}

	return processQuery(pqs.Queries)
}

func processQuery(queries []ast.Body) (Result, error) {
	exprs := []*Expr{}
	for _, query := range queries {
		if len(query) == 0 {
			// always allow
			return Result{Defined: true}, nil
		}

		for _, astExpr := range query {
			if !astExpr.IsCall() {
				continue
			}

			gotOps := len(astExpr.Operands())
			if gotOps != 2 {
				return Result{}, errors.Errorf("invalid expression: expecting 2 operands but got %d", gotOps)
			}

			var left, right *Operand
			for i, term := range astExpr.Operands() {
				var operand *Operand
				if ast.IsConstant(term.Value) {
					v, err := ast.JSON(term.Value)
					if err != nil {
						return Result{}, errors.Wrap(err, "convert term value to json")
					}
					operand = &Operand{T: OperandTypeConstant, V: v}
				} else {
					processedTerm := processTerm(term.String())
					if processedTerm == nil {
						return Result{}, nil
					}
					operand = &Operand{T: OperandTypeIndexField, V: processedTerm}
				}

				if i == 1 {
					left = operand
				} else {
					right = operand
				}
			}

			exprs = append(exprs, &Expr{
				Operator: strToOp(astExpr.Operator().String()),
				Left:     left,
				Right:    right,
			})
		}
	}

	return Result{
		Defined: true,
		Exprs:   exprs,
	}, nil
}

func processTerm(query string) []string {
	splitQ := strings.Split(query, ".")
	var result []string
	for _, term := range splitQ {
		result = append(result, removeOpenBrace(term))
	}

	indexName, fieldName := result[1], result[2]
	if len(result) > 2 {
		fieldName = strings.Join(result[2:], ".")
	}

	return []string{indexName, fieldName}
}

func removeOpenBrace(input string) string {
	return strings.Split(input, "[")[0]
}
