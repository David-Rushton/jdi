package cli

import "regexp"

const (
	CliTag          = "cli"
	CliTagSeparator = '|'
	// OptionGroupPattern                = `^-[a-zA-Z0-9]{2,}$`
	ShortOptionPattern                = `^-[a-zA-Z0-9]$`
	LongOptionPattern                 = `^--[a-zA-Z0-9][a-zA-Z0-9-]+$`
	PositionalArgumentNamePattern     = `^[\[<][a-zA-Z0-9][a-zA-Z0-9-]*[\]>]$`
	PositionalArgumentPositionPattern = `^\d+$`
	PositionalArgumentRequiredPrefix  = "<"
	OptionalTagPrefix                 = "-"
)

var (
	// optionGroupRegex                = regexp.MustCompile(OptionGroupPattern)
	shortOptionRegex                = regexp.MustCompile(ShortOptionPattern)
	longOptionRegex                 = regexp.MustCompile(LongOptionPattern)
	positionalArgumentNameRegex     = regexp.MustCompile(PositionalArgumentNamePattern)
	positionalArgumentPositionRegex = regexp.MustCompile(PositionalArgumentPositionPattern)
)
