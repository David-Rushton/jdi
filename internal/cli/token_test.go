package cli

import (
	"slices"
	"testing"
)

func TestToToken_ReturnsExpectedValues(t *testing.T) {
	testCases := []struct {
		args     []string
		expected []string
	}{
		{
			args:     []string{"new", "test", "--name", "some-test-case", "-v"},
			expected: []string{"new", "test", "--name", "some-test-case", "-v"},
		},
	}

	for _, testCase := range testCases {
		var actual []string

		token := tokeniseArgs(testCase.args)
		for !token.isEnd() {
			actual = append(actual, token.getValue())
			token = token.next()
		}

		if !slices.Equal(testCase.expected, actual) {
			t.Errorf("unexpected token value(s): { expected = %q, actual = %q }", testCase.expected, actual)
		}
	}
}
