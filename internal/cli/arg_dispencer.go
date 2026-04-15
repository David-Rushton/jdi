package cli

// take args
// match options
// then match subcommands
//
// inspect remaining

import "slices"

type argDispencer struct {
	args map[int]*arg
}

type arg struct {
	value     string
	dispensed bool
}

func newArgDispencer(args []string) argDispencer {
	result := argDispencer{
		args: map[int]*arg{},
	}

	for i, value := range args {
		result.args[i] = &arg{
			value:     value,
			dispensed: false,
		}
	}

	return result
}

func (a *argDispencer) getSubcommand(position int) (string, bool) {
	if position < 0 || position >= len(a.args) {
		return "", false
	}

	available := -1
	for i := 0; i < len(a.args); i++ {
		if !a.args[i].dispensed {
			available++

			if available == position {
				a.args[i].dispensed = true
				return a.args[i].value, true
			}

			if available > position {
				break
			}
		}
	}

	return "", false
}

func (a *argDispencer) getOptionArgument(shortName, longName string) (string, bool) {
	var value string
	var found bool

	names := []string{}
	for _, name := range []string{shortName, longName} {
		if name != "" {
			names = append(names, name)
		}
	}

	for i, arg := range a.args {
		if arg.dispensed {
			continue
		}

		if slices.Contains(names, arg.value) {
			if nextArg, exists := a.args[i+1]; exists {
				if !nextArg.dispensed {
					arg.dispensed = true
					nextArg.dispensed = true
					found = true
					value = nextArg.value
				}
			}

			break
		}
	}

	return value, found
}

func (a *argDispencer) optionExists(shortName, longName string) bool {
	names := []string{}
	for _, name := range []string{shortName, longName} {
		if name != "" {
			names = append(names, name)
		}
	}

	for _, arg := range a.args {
		if arg.dispensed {
			continue
		}

		if slices.Contains(names, arg.value) {
			arg.dispensed = true
			return true
		}
	}

	return false
}
