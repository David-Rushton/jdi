package main

import (
	"fmt"

	"github.com/David-Rushton/jdi/internal/cli"
)

type TestCommand struct {
	Name string `cli:"-n|--name|give me a name"`
	Path string `cli:"0|<path>|path to some thing"`
}

func (t *TestCommand) Invoke() error {
	fmt.Println("Test command invoked")
	fmt.Println("--------------------")
	fmt.Printf("<path>       %s\n", t.Path)
	fmt.Printf("-n|--name    %s\n", t.Name)
	return nil
}

func main() {
	fmt.Println("hola mundo")

	app := cli.App{Name: "test app"}
	if e := app.AddCommnad("test", "some command", &TestCommand{}); e != nil {
		panic(e)
	}

	args := []string{"test", "-n", "some name"}
	if e := app.Run(args[1:]); e != nil {
		panic(e)
	}
}
