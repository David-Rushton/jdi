package cli

import (
	"testing"
)

type testCommand struct {
	Path string `cli:"0|<path>|path to some thing"`
	Name string `cli:"-n|--name|give me a name"`
}

func (t *testCommand) Invoke() error {
	// no-op.
	return nil
}

func TestMatch(t *testing.T) {
	testCommand := &testCommand{
		Path: "",
		Name: "",
	}
	expectedName := "test-path"
	expectedPath := "test-name"
	testArgs := []string{expectedPath, "--name", expectedName}

	command, err := parseCommand("test", "test description", testCommand)
	if err != nil {
		t.Errorf("command parsing failed: %v", err)
	}

	err = match(command, toTokens(testArgs))
	if err != nil {
		t.Errorf("param/arg matching failed: %v", err)
	}

	if testCommand.Path != expectedPath {
		t.Errorf("param/arg matching failed: { expected path: %s, actual path: %s}", expectedPath, testCommand.Path)
	}

	if testCommand.Name != expectedName {
		t.Errorf("param/arg matching failed: { expected name: %s, actual name: %s}", expectedName, testCommand.Name)
	}
}
