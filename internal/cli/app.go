package cli

import (
	"fmt"
	"reflect"
)

type App struct {
	Name string

	commands map[string]*command
}

func (a *App) Run(args []string) error {
	fmt.Println("running!")

	if len(args) == 0 {
		// TODO: Show help.
		return fmt.Errorf("No command requested")
	}

	command, found := a.commands[args[0]]
	if !found {
		return fmt.Errorf("Cannot find command %s", args[0])
	}

	return command.invoke.Invoke()
}

func (a *App) AddCommnad(name, description string, invoke Invoker) error {
	if a.commands == nil {
		a.commands = map[string]*command{}
	}

	if _, exists := a.commands[name]; exists {
		return fmt.Errorf("A command with the name %s already exists", name)
	}

	a.commands[name] = &command{
		name,
		description,
		invoke,
		[]*subcommand{},
		[]*option{},
		map[string]reflect.Value{},
	}

	if e := a.commands[name].parse(); e != nil {
		return e
	}

	return nil
}
