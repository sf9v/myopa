package myopa

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/pkg/errors"
)

// Engine is an OPA engine
type Engine struct {
	policyFile string
	policy     []byte
}

// Result is a compilation result
type Result struct {
	Defined bool
	Exprs   []*Expr
}

// Expr is an expression
type Expr struct {
	// Op is the operator
	Op Op
	// L is the l-value
	L Val
	// R is the r-value
	R Val
}

// Op is an operator
type Op int

func (op Op) String() string {
	return [...]string{
		"invalid",
		"equal",
	}[op]
}

// List of operators
const (
	OpEq Op = iota + 1
)

// Val is a value
type Val struct {
	T VT
	V interface{}
}

// VT is a value type
type VT int

func (vt VT) String() string {
	return [...]string{
		"invalid",
		"constant",
		"key-value",
	}[vt]
}

// List of value types
const (
	VTConstant VT = iota + 1
	VTKeyValue
)

// M is an alias to map of interfaces
type M map[string]interface{}

// New reads a policy file and returns a new engine
func New(policyFile string) (*Engine, error) {
	policy, err := ioutil.ReadFile(policyFile)
	if err != nil {
		return nil, errors.Wrap(err, "read policy file")
	}

	return &Engine{
		policyFile: policyFile,
		policy:     policy,
	}, nil
}

// Compile compiles the query
func (e *Engine) Compile(ctx context.Context, query string,
	unknowns []string, input M) (Result, error) {
	rg := rego.New(
		rego.Query(query),
		rego.Module(e.policyFile, string(e.policy)),
		rego.Input(input),
		rego.Unknowns(unknowns),
	)

	pq, err := rg.Partial(ctx)
	if err != nil {
		return Result{}, err
	}

	if len(pq.Queries) == 0 {
		// always deny
		return Result{Defined: false}, nil
	}

	return processQuery(pq.Queries)
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

			expectOps := 2
			gotOps := len(astExpr.Operands())
			if gotOps != expectOps {
				return Result{}, errors.Errorf("invalid expression: expecting %d operands but got %d", expectOps, gotOps)
			}

			// operator
			var op Op
			switch astExpr.Operator().String() {
			default:
				op = OpEq
			}

			expr := &Expr{Op: op}
			for i, term := range astExpr.Operands() {
				var val Val
				if ast.IsConstant(term.Value) {
					v, err := ast.JSON(term.Value)
					if err != nil {
						return Result{}, errors.Wrap(err, "convert term value to json")
					}

					val = Val{T: VTConstant, V: v}
				} else {
					processedTerm := processTerm(term.String())
					if processedTerm == nil {
						return Result{}, nil
					}

					val = Val{T: VTKeyValue, V: processedTerm}
				}

				// we only expect two operands
				if i == 0 {
					expr.L = val
				} else if i == 1 {
					expr.R = val
				}
			}

			exprs = append(exprs, expr)
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

	if result == nil {
		return nil
	}

	indexName := result[1]
	fieldName := result[2]
	if len(result) > 2 {
		fieldName = strings.Join(result[2:], ".")
	}

	return []string{indexName, fieldName}
}

func removeOpenBrace(input string) string {
	return strings.Split(input, "[")[0]
}
