package cli

import (
	"slices"
	"testing"
)

func TestSplit(t *testing.T) {
	var testCases = []struct {
		input     string
		delimiter rune
		expected  []string
	}{
		{
			input:     "some|pipe|delimted text",
			delimiter: '|',
			expected:  []string{"some", "pipe", "delimted text"},
		},
		{
			input:     "some|pipe|delimted text with escaped delmiter||",
			delimiter: '|',
			expected:  []string{"some", "pipe", "delimted text with escaped delmiter|"},
		},
		{
			input:     "space delimited text",
			delimiter: ' ',
			expected:  []string{"space", "delimited", "text"},
		},
		{
			input:     "||",
			delimiter: '|',
			expected:  []string{"|"},
		},
	}

	for _, testCase := range testCases {
		actual := split(testCase.input, testCase.delimiter)
		if !slices.Equal(testCase.expected, actual) {
			t.Errorf("split returned unexpected result: { expected = %v, actual = %v }", testCase.expected, actual)
		}
	}
}
