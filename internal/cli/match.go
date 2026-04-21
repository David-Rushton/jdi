package cli

import (
	"errors"
	"fmt"
)

const (
	Base10 = 10
)

func match(cmd *command, tokens argTokens) error {
	var errs []error

	// Optional arguments.
	for tokens.next() {
		current := tokens.current()
		next, nextExists := tokens.peek()

		if !nextExists {
			continue
		}

		requiredTypesFound := current.getType() == tokenTypeMaybeOption && next.getType() == tokenTypeMaybeOptionArgument
		if !requiredTypesFound {
			continue
		}

		// TODO: command -n name -n new-name == {{ name = new-name }}
		if option, optionExists := cmd.optionalArguments[current.getValue()]; optionExists {
			if option.field.hasValue {
				errs = append(errs, fmt.Errorf("cannot set value of optional argument because it has already been set"))
			}

			if err := option.field.set(next.getValue()); err != nil {
				errs = append(errs, fmt.Errorf("cannot set value of optional argument because %e", err))
			}

			current.consume()
			next.consume()
		}
	}

	// Optional switches.
	tokens.resetPosition()
	for tokens.next() {
		current := tokens.current()

		if current.getType() != tokenTypeMaybeOption {
			continue
		}

		// TODO: command -v -v == {{ -v == true }}
		if option, optionExists := cmd.optionalSwitches[current.getValue()]; optionExists {
			if option.field.hasValue {
				errs = append(errs, fmt.Errorf("cannot set value of optional argument because it has already been set"))
			}

			if err := option.field.set("true"); err != nil {
				errs = append(errs, fmt.Errorf("cannot set value of optional switch because %e", err))
			}

			current.consume()
		}
	}

	// Positional arguments.
	positionalTokens := tokens.getRemaining()
	for i, positionalToken := range positionalTokens {
		if i >= len(cmd.posistionalArguments) {
			errs = append(errs, fmt.Errorf("unexpected positional argument '%s' at position %d", positionalToken.getValue(), i))
		}

		positionalArg := cmd.posistionalArguments[i]
		if err := positionalArg.field.set(positionalToken.getValue()); err != nil {
			errs = append(errs, fmt.Errorf("cannot set value of positional argument %d because %e", i, err))
		}
	}

	return errors.Join(errs...)
}
