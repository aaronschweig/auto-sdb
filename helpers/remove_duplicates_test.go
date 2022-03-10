package helpers_test

import (
	"testing"

	"github.com/aaronschweig/auto-sdb/helpers"
	"github.com/stretchr/testify/require"
)

func TestRemoveDuplicates(t *testing.T) {
	testCases := []struct {
		name     string
		in       []string
		expected []string
	}{
		{
			name:     "should dedupe array but keep distinct items",
			in:       []string{"H222", "H222", "H221"},
			expected: []string{"H222", "H221"},
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			result := helpers.RemoveDuplicates(test.in)
			require.ElementsMatch(t, test.expected, result)
		})
	}
}
