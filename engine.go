package myopa

import (
	"context"
	"io/ioutil"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/pkg/errors"
)

type Engine struct {
	policyFile string
	policy     []byte
}

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

type Result struct {
	Op Op
	L  Val
	R  Val
}

func (e *Engine) Compile(ctx context.Context, query string,
	unknowns []string, input M) ([]*Result, error) {
	rg := rego.New(
		rego.Query(query),
		rego.Module(e.policyFile, string(e.policy)),
		rego.Input(input),
		rego.Unknowns(unknowns),
	)

	pq, err := rg.Partial(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "partial evaluation")
	}

	if len(pq.Queries) == 0 {
		rs, err := rg.Eval(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "eval")
		}

		if len(rs) != 1 || len(rs[0].Expressions) != 1 {
			return nil, errors.New("empty result set")
		}
	}

	results := []*Result{}
	for _, query := range pq.Queries {
		for _, expr := range query {
			if !expr.IsCall() {
				continue
			}

			expectNumOps := 2
			gotNumOps := len(expr.Operands())
			if gotNumOps != expectNumOps {
				return nil, errors.Errorf("invalid expression: expecting %d operands but got %d", expectNumOps, gotNumOps)
			}

			// operator
			var op Op
			switch expr.Operator().String() {
			default:
				op = OpEq
			}

			result := &Result{Op: op}
			for i, term := range expr.Operands() {
				var val Val
				if ast.IsConstant(term.Value) {
					v, err := ast.JSON(term.Value)
					if err != nil {
						return nil, errors.Wrap(err, "convert term value to json")
					}

					val = Val{
						T: VTConstant,
						V: v,
					}
				} else {
					processedTerm := processTerm(term.String())
					if processedTerm == nil {
						return results, nil
					}

					val = Val{
						T: VTKeyValue,
						V: processedTerm,
					}
				}

				// we only expect two operands
				if i == 0 {
					result.L = val
				} else if i == 1 {
					result.R = val
				}
			}

			results = append(results, result)
		}
	}

	return results, nil
}
