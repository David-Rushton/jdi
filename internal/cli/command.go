package cli

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Invoker interface {
	Invoke() error
}

type command struct {
	name                 string
	description          string
	invoke               Invoker
	posistionalArguments map[int]*positionalArgument
	options              map[string]*option
}

type positionalArgument struct {
	position    int
	name        string
	description string
	value       reflect.Value
}

func (pa *positionalArgument) required() bool {
	return strings.HasPrefix(pa.name, PositionalArgumentRequiredPrefix)
}

func (pa *positionalArgument) kind() reflect.Kind {
	return pa.value.Kind()
}

type option struct {
	shortName   string
	longName    string
	description string
	value       reflect.Value
}

func (o *option) names() []string {
	var result []string

	if o.shortName != "" {
		result = append(result, o.shortName)
	}

	if o.longName != "" {
		result = append(result, o.longName)
	}

	return result
}

func (o *option) requiresArgument() bool {
	return o.kind() != reflect.Bool
}

func (o *option) kind() reflect.Kind {
	return o.value.Kind()
}

// TODO: We should return aggregated errors.
func parseCommand(name, description string, invoke Invoker) (*command, error) {
	cmd := &command{
		name,
		description,
		invoke,
		map[int]*positionalArgument{},
		map[string]*option{},
	}

	paramsT := reflect.TypeOf(cmd.invoke)
	paramsV := reflect.ValueOf(cmd.invoke)

	if paramsT.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("command %s does not point to a command object", cmd.name)
	}

	for i := 0; i < paramsT.Elem().NumField(); i++ {
		fieldT := paramsT.Elem().Field(i)
		fieldV := paramsV.Elem().Field(i)
		tag := fieldT.Tag.Get(CliTag)

		if tag != "" {
			switch strings.HasPrefix(tag, OptionalTagPrefix) {
			case true:
				option, err := parseOption(tag, fieldV)
				if err != nil {
					return nil, err
				}

				for _, name := range option.names() {
					if name != "" {
						if _, exists := cmd.options[name]; exists {
							return nil, fmt.Errorf("duplicate option name: %s", name)
						}
					}

					cmd.options[name] = option
				}

			default:
				positionalArg, err := parsePosistionalArgument(tag, fieldV)
				if err != nil {
					return nil, err
				}

				if _, exists := cmd.posistionalArguments[positionalArg.position]; exists {
					return nil, fmt.Errorf("duplicate argument position: %v", positionalArg.position)
				}

				cmd.posistionalArguments[positionalArg.position] = positionalArg
			}
		}
	}

	// Validate positional arguments.
	mustBeOptional := false
	for i := 0; i < len(cmd.posistionalArguments); i++ {
		if _, exists := cmd.posistionalArguments[i]; !exists {
			return nil, fmt.Errorf("positional arguments should be consecutive, missing position %d", cmd.posistionalArguments[i].position)
		}

		if mustBeOptional && cmd.posistionalArguments[i].required() {
			return nil, fmt.Errorf("required arguments cannot appear after optional arguments")
		}

		if !cmd.posistionalArguments[i].required() {
			mustBeOptional = true
		}
	}

	return cmd, nil
}

func parseOption(tag string, field reflect.Value) (*option, error) {
	var tagElements = split(tag, '|')
	if len(tagElements) == 0 || len(tagElements) > 3 {
		return nil, fmt.Errorf("invalid option tag format: %s", tag)
	}

	shortName := ""
	longName := ""
	for i := 0; i < len(tagElements)-1; i++ {
		if shortOptionRegex.MatchString(tagElements[i]) {
			shortName = tagElements[i]
			continue
		}

		if longOptionRegex.MatchString(tagElements[i]) {
			longName = tagElements[i]
			continue
		}
	}

	if shortName == "" || longName == "" {
		return nil, fmt.Errorf("option tags must contain at least one of short and long name: %s", tag)
	}

	description := tagElements[len(tagElements)-1]
	if description == "" {
		return nil, fmt.Errorf("option tags must contain a description: %s", tag)
	}

	return &option{
		shortName,
		longName,
		description,
		field,
	}, nil
}

func parsePosistionalArgument(tag string, field reflect.Value) (*positionalArgument, error) {
	const positionIndex = 0
	const nameIndex = 1
	const descriptionIndex = 2

	var tagElements = split(tag, '|')
	if len(tagElements) != 3 {
		return nil, fmt.Errorf("invalid positional argument tag format: %s", tag)
	}

	if !positionalArgumentPositionRegex.MatchString(tagElements[positionIndex]) {
		return nil, fmt.Errorf("invalid positional argument in tag: %s", tag)
	}

	if !positionalArgumentNameRegex.MatchString(tagElements[nameIndex]) {
		return nil, fmt.Errorf("invalid positional argument name in tag: %s", tag)
	}

	position, _ := strconv.ParseInt(tagElements[positionIndex], 10, 64)

	return &positionalArgument{
		int(position),
		tagElements[nameIndex],
		tagElements[descriptionIndex],
		field,
	}, nil
}
