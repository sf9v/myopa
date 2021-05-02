package myopa_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sf9v/myopa"
)

type M map[string]interface{}

func TestEngine(t *testing.T) {
	defaultQuery := "data.example.allow = true"
	unknowns := []string{
		"data.pages",
		"data.page_managers",
	}
	tests := []struct {
		name      string
		query     string
		input     M
		defined   bool
		exprCount int
	}{
		{
			name:  "anon user is not allowed to create page",
			query: defaultQuery,
			input: M{
				"action": "create",
				"object": M{
					"type": "page",
				},
				"user": "anon",
			},
			defined: false,
		},
		{
			name:  "user is allowed create page",
			query: defaultQuery,
			input: M{
				"action": "create",
				"object": M{
					"type": "page",
				},
				"user": "user-1234",
			},
			defined: true,
		},
		{
			name:  "all user is allowed read page",
			query: defaultQuery,
			input: M{
				"action": "read",
				"object": M{
					"type": "page",
					"id":   "page-1234",
				},
				"user": "user-1234",
			},
			defined: true,
		},
		{
			name:  "any page manager is allowed to update page",
			query: defaultQuery,
			input: M{
				"action": "update",
				"object": M{
					"type": "page",
					"id":   "page-1234",
				},
				"user": "user-1234",
			},
			defined:   true,
			exprCount: 6,
		},
		{
			name:  "only a page admin is allowed to delete page",
			query: defaultQuery,
			input: M{
				"action": "delete",
				"object": M{
					"type": "page",
					"id":   "page-1234",
				},
				"user": "user-1234",
			},
			defined:   true,
			exprCount: 3,
		},
	}

	policyFile := "example.rego"
	b, err := ioutil.ReadFile(policyFile)
	require.NoError(t, err)

	e := myopa.New(policyFile, b)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := e.Compile(context.TODO(), defaultQuery, unknowns, tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.defined, result.Defined)
			assert.Len(t, result.Exprs, tc.exprCount)
		})
	}
}

func BenchmarkEngine(b *testing.B) {
	policyFile := "example.rego"
	module, err := ioutil.ReadFile(policyFile)
	require.NoError(b, err)

	e := myopa.New(policyFile, module)

	defaultQuery := "data.example.authz.allow = true"
	unknowns := []string{
		"data.pages",
		"data.page_managers",
	}
	input := M{
		"action": "create",
		"object": M{
			"type": "page",
		},
		"user": "anon",
	}
	for i := 0; i < b.N; i++ {
		_, err = e.Compile(context.TODO(), defaultQuery, unknowns, input)
		require.NoError(b, err)
	}
}
