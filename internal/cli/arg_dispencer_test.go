package cli

import "testing"

func TestArgDispencer_ReturnsSubcommands(t *testing.T) {
	var testCases = []struct {
		args          []string
		posistion     int
		expectedValue string
		expectedFound bool
	}{
		{
			args:          []string{"x", "y", "z"},
			posistion:     0,
			expectedValue: "x",
			expectedFound: true,
		},
		{
			args:          []string{"x", "y", "z"},
			posistion:     1,
			expectedValue: "y",
			expectedFound: true,
		},
		{
			args:          []string{"x", "y", "z"},
			posistion:     2,
			expectedValue: "z",
			expectedFound: true,
		},
		{
			args:          []string{"x", "y", "z"},
			posistion:     3,
			expectedValue: "",
			expectedFound: false,
		},
	}

	for _, testCase := range testCases {
		dispencer := newArgDispencer(testCase.args)
		actualValue, actualFound := dispencer.getSubcommand(testCase.posistion)

		if actualFound != testCase.expectedFound {
			t.Errorf("unexpected found status { expected: %t, actual %t }", testCase.expectedFound, actualFound)
		}

		if testCase.expectedFound {
			if actualValue != testCase.expectedValue {
				t.Errorf("unexpected value { expected: %s, actual %s }", testCase.expectedValue, actualValue)
			}
		}
	}
}

func TestArgDispencer_ReturnsOptionArguments(t *testing.T) {
	var testCases = []struct {
		args          []string
		shortName     string
		longName      string
		expectedValue string
		expectedFound bool
	}{
		{
			args:          []string{"x", "--long-name", "value", "--not-tested"},
			shortName:     "-l",
			longName:      "--long-name",
			expectedValue: "value",
			expectedFound: true,
		},
		{
			args:          []string{"x", "-s", "value", "--not-tested"},
			shortName:     "-s",
			longName:      "--short-name-not-in-args",
			expectedValue: "value",
			expectedFound: true,
		},
		{
			args:          []string{"x", "--long-name", "value-1", "-s", "123"},
			shortName:     "-s",
			longName:      "",
			expectedValue: "123",
			expectedFound: true,
		},
		{
			args:          []string{"x", "--long-name", "value", "--not-tested"},
			shortName:     "",
			longName:      "--arg-not-provided",
			expectedValue: "",
			expectedFound: false,
		},
	}

	for _, testCase := range testCases {
		dispencer := newArgDispencer(testCase.args)
		actualValue, actualFound := dispencer.getOptionArgument(testCase.shortName, testCase.longName)

		if actualFound != testCase.expectedFound {
			t.Errorf("unexpected found status { expected: %t, actual %t }", testCase.expectedFound, actualFound)
		}

		if testCase.expectedFound {
			if actualValue != testCase.expectedValue {
				t.Errorf("unexpected value { expected: %s, actual %s }", testCase.expectedValue, actualValue)
			}
		}
	}
}
