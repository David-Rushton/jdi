package cli

import (
	"strings"
)

// Slices s into multiple strings.
//
// Seperator can be escapped by providing two consecutive seperators.
//
// Slices are trimmed.
//
// Empty strings are suppressed.
func split(s string, separator rune) []string {
	result := []string{}

	buf := strings.Builder{}
	flush := func() {
		candidate := strings.TrimSpace(buf.String())
		if len(candidate) > 0 {
			result = append(result, candidate)
		}
		buf.Reset()
	}

	var consecutiveSeparators int
	for _, r := range s {

		if r == separator {
			consecutiveSeparators++

			if consecutiveSeparators%2 == 0 {
				buf.WriteRune(separator)
			}

			continue
		}

		if consecutiveSeparators > 0 {
			consecutiveSeparators = 0
			flush()
		}

		buf.WriteRune(r)
	}
	flush()

	return result
}
