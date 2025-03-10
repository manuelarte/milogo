package config

import "github.com/manuelarte/milogo/internal/parser"

type Config struct {
	QueryParamField string
	WrapperField    string
	Parser          parser.Parser
}

func DefaultConfig(configOptions ...Option) Config {
	c := Config{
		QueryParamField: "fields",
		Parser:          parser.NewParser(),
	}
	for _, co := range configOptions {
		co(&c)
	}

	return c
}

type Option func(c *Config)

func WithWrapField(wrapperField string) (Option, error) {
	// validate wrapper field is not several words, but only one and alphanumeric
	return func(c *Config) {
		c.WrapperField = wrapperField
	}, nil
}
