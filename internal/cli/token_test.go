package cli

import (
	"slices"
	"testing"
)

func TestToTokens_ReturnsExpectedValues(t *testing.T) {
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

		tokens := toTokens(testCase.args)
		for tokens.next() {
			token := tokens.current()
			actual = append(actual, token.getValue())
		}

		if !slices.Equal(testCase.expected, actual) {
			t.Errorf("unexpected token value(s): { expected = %q, actual = %q }", testCase.expected, actual)
		}
	}
}

func TestToTokens_ReturnsExpectedTypes(t *testing.T) {
	testCases := []struct {
		args     []string
		expected []tokenType
	}{
		{
			args: []string{"new", "test", "--name", "some-test-case", "--", "-v"},
			expected: []tokenType{
				tokenTypeMaybePositionalArguement,
				tokenTypeMaybePositionalArguement,
				tokenTypeMaybeOption,
				tokenTypeMaybeOptionArgument,
				tokenTypeTerminator,
				tokenTypePositionalArguement,
			},
		},
	}

	for _, testCase := range testCases {
		var actual []tokenType

		tokens := toTokens(testCase.args)
		for tokens.next() {
			token := tokens.current()
			actual = append(actual, token.getType())
		}

		if !slices.Equal(testCase.expected, actual) {
			t.Errorf("unexpected token value(s): { expected = %v, actual = %v }", testCase.expected, actual)
		}
	}
}
