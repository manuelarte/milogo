package pkg

type Config struct {
	QueryParamField string
	WrapperField    string
	Parser          Parser
}

func DefaultConfig(configOptions ...ConfigOption) Config {
	c := Config{
		QueryParamField: "fields",
		Parser:          NewParser(),
	}
	for _, co := range configOptions {
		co(&c)
	}

	return c
}

type ConfigOption func(c *Config)

func WithWrapField(wrapperField string) (ConfigOption, error) {
	// validate wrapper field is not several words, but only one and alphanumeric
	return func(c *Config) {
		c.WrapperField = wrapperField
	}, nil
}
