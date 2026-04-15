package cli

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

const (
	CliTag                     = "cli"
	CliTagSeparator            = '|'
	ShortOptionPattern         = `^-[a-zA-Z0-9]$`
	LongOptionPattern          = `^--[a-zA-Z0-9][a-zA-Z0-9-]+$`
	SubcommandNamePattern      = `^[\[<][a-zA-Z0-9][a-zA-Z0-9-]*[\]>]$`
	SubcommandPosistionPattern = `^\d+$`
)

var (
	shortOptionRegex         = regexp.MustCompile(ShortOptionPattern)
	longOptionRegex          = regexp.MustCompile(LongOptionPattern)
	subcommandNameRegex      = regexp.MustCompile(SubcommandNamePattern)
	subcommandPosistionRegex = regexp.MustCompile(SubcommandPosistionPattern)
)

type Invoker interface {
	Invoke() error
}

type command struct {
	name        string
	description string
	invoke      Invoker
	subcommands []*subcommand
	options     []*option
	parameters  map[string]reflect.Value
}

type subcommand struct {
	name        string
	posistion   string
	description string
	value       reflect.Value
}

type option struct {
	shortName   string
	longName    string
	description string
	value       reflect.Value
}

// TODO: We should return aggregated errors.
func (c *command) parse() error {
	paramsT := reflect.TypeOf(c.invoke)
	paramsV := reflect.ValueOf(c.invoke)

	if paramsT.Kind() != reflect.Pointer {
		return fmt.Errorf("command %s does not point to a command object", c.name)
	}

	for i := 0; i < paramsT.Elem().NumField(); i++ {
		fieldT := paramsT.Elem().Field(i)
		fieldV := paramsV.Elem().Field(i)
		tag := fieldT.Tag.Get(CliTag)

		if tag != "" {
			// Extract values.
			tagElements := split(tag, CliTagSeparator)

			if len(tagElements) < 2 {
				return fmt.Errorf("invalid tag format: %s", tag)
			}

			tagDescription := ""
			shortOption := ""
			longOption := ""
			subcommandName := ""
			subcommandPosistion := ""

			for _, tagElement := range tagElements {

				if shortOptionRegex.MatchString(tagElement) {
					if len(shortOption) > 0 {
						return fmt.Errorf("short option defined twice")
					}
					shortOption = tagElement
					continue
				}

				if longOptionRegex.MatchString(tagElement) {
					if len(longOption) > 0 {
						return fmt.Errorf("long option defined twice")
					}
					longOption = tagElement
					continue
				}

				if subcommandNameRegex.MatchString(tagElement) {
					if len(subcommandName) > 0 {
						return fmt.Errorf("subcommand defined twice")
					}
					subcommandName = tagElement
					continue
				}

				if subcommandPosistionRegex.MatchString(tagElement) {
					if subcommandPosistion != "" {
						return fmt.Errorf("subcommand position defined twice")
					}
					subcommandPosistion = tagElement
					continue
				}

				if len(tagDescription) > 0 {
					return fmt.Errorf("description defined twice")
				}
				tagDescription = tagElement

				// Validate.
				if len(tagDescription) == 0 {
					return fmt.Errorf("description is required")
				}

				if shortOption == "" && longOption == "" && subcommandName == "" {
					return fmt.Errorf("one of option and subcommand must be provided")
				}

				if shortOption != "" || longOption != "" {
					// Option.
					if subcommandName != "" || subcommandPosistion != "" {
						return fmt.Errorf("cannot combine options and subcommands")
					}

					c.options = append(c.options, &option{
						shortOption,
						longOption,
						tagDescription,
						fieldV})

				} else {
					// Subcommand.
					if shortOption != "" || longOption != "" {
						return fmt.Errorf("cannot combine options and subcommands")
					}

					c.subcommands = append(c.subcommands, &subcommand{
						subcommandName,
						subcommandPosistion,
						tagDescription,
						fieldV})
				}
			}

			if !fieldV.CanSet() {
				return fmt.Errorf("cannot update field %s", fieldT.Name)
			}
		}
	}

	// Validate.
	for _, option := range c.options {
		if option.shortName != "" {
			if _, exists := c.parameters[option.shortName]; exists {
				return fmt.Errorf("option %s defined twice", option.shortName)
			}
			c.parameters[option.shortName] = option.value
		}

		if option.longName != "" {
			if _, exists := c.parameters[option.longName]; exists {
				return fmt.Errorf("option %s defined twice", option.longName)
			}
			c.parameters[option.longName] = option.value
		}
	}

	posistions := map[int]int{}
	for _, subcommand := range c.subcommands {
		if _, exists := c.parameters[subcommand.posistion]; exists {
			return fmt.Errorf("subcommand posisiton %s defined twice", subcommand.posistion)
		}
		c.parameters[subcommand.posistion] = subcommand.value

		posistion, e := strconv.ParseInt(subcommand.posistion, 10, 64)
		if e != nil {
			panic(fmt.Sprintf("unexpected failure to parse int %s", subcommand.posistion))
		}
		posistions[int(posistion)]++
	}

	// Validate.
	for i := 0; i < len(posistions); i++ {
		switch posistions[i] {
		case 1:
			// No-op.
			// Validation passed.
		case 0:
			return fmt.Errorf("missing subcommand posistion %d", i)
		default:
			return fmt.Errorf("duplicated subcommand posistion %d", i)
		}
	}

	return nil
}
