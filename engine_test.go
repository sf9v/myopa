package myopa_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sf9v/myopa"
)

func TestEngine(t *testing.T) {
	defaultQuery := "data.example.allow = true"
	unknowns := []string{
		"data.pages",
		"data.page_managers",
	}
	tests := []struct {
		name      string
		query     string
		input     myopa.M
		defined   bool
		exprCount int
	}{
		{
			name:  "anon user is not allowed to create page",
			query: defaultQuery,
			input: myopa.M{
				"action": "create",
				"object": myopa.M{
					"type": "page",
				},
				"user": "anon",
			},
			defined: false,
		},
		{
			name:  "an auth user is allowed create a page",
			query: defaultQuery,
			input: myopa.M{
				"action": "create",
				"object": myopa.M{
					"type": "page",
				},
				"user": "user-1234",
			},
			defined: true,
		},
		{
			name:  "all user can read a page",
			query: defaultQuery,
			input: myopa.M{
				"action": "read",
				"object": myopa.M{
					"type": "page",
					"id":   "page-1234",
				},
				"user": "user-1234",
			},
			defined: true,
		},
		{
			name:  "a page manager is allowed to update the page",
			query: defaultQuery,
			input: myopa.M{
				"action": "update",
				"object": myopa.M{
					"type": "page",
					"id":   "page-1234",
				},
				"user": "user-1234",
			},
			defined:   true,
			exprCount: 2,
		},
		{
			name:  "page admin is allowed to delete the page",
			query: defaultQuery,
			input: myopa.M{
				"action": "delete",
				"object": myopa.M{
					"type": "page",
					"id":   "page-1234",
				},
				"user": "user-1234",
			},
			defined:   true,
			exprCount: 3,
		},
	}

	e, err := myopa.New("example.rego")
	require.NoError(t, err)

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
	e, err := myopa.New("example.rego")
	require.NoError(b, err)

	defaultQuery := "data.example.authz.allow = true"
	unknowns := []string{
		"data.pages",
		"data.page_managers",
	}
	input := myopa.M{
		"action": "create",
		"object": myopa.M{
			"type": "page",
		},
		"user": "anon",
	}
	for i := 0; i < b.N; i++ {
		_, err = e.Compile(context.TODO(), defaultQuery, unknowns, input)
		require.NoError(b, err)
	}
}
