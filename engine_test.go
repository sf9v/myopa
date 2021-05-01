package myopa_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sf9v/myopa"
)

func TestEngine(t *testing.T) {
	defaultAllow := "data.example.authz.allow = true"
	// defaultNotAllow := "data.example.authz.allow"
	unknowns := []string{
		"data.pages",
		"data.page_managers",
	}
	tests := []struct {
		name        string
		query       string
		input       myopa.M
		resultCount int
	}{
		{
			name:  "read page",
			query: defaultAllow,
			input: myopa.M{
				"action": "read",
				"object": myopa.M{
					"type": "page",
					"id":   "page-1234",
				},
				"user": "user-1234",
			},
			resultCount: 0,
		},
		{
			name:  "update page",
			query: defaultAllow,
			input: myopa.M{
				"action": "update",
				"object": myopa.M{
					"type": "page",
					"id":   "page-1234",
				},
				"user": "user-1234",
			},
			resultCount: 3,
		},
		{
			name:  "delete page",
			query: defaultAllow,
			input: myopa.M{
				"action": "delete",
				"object": myopa.M{
					"type": "page",
					"id":   "page-1234",
				},
				"user": "user-1234",
			},
			resultCount: 4,
		},
	}

	e, err := myopa.New("example.rego")
	require.NoError(t, err)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			results, err := e.Compile(context.TODO(), defaultAllow, unknowns, tc.input)
			require.NoError(t, err)

			assert.Len(t, results, tc.resultCount)
		})
	}
}
