package yesql

import (
	"github.com/izolate/yesql/bindvar"
	"github.com/izolate/yesql/template"
)

// Config stores runtime config for yesql.
type Config struct {
	driver string
	tpl    template.Executer
	bvar   bindvar.Parser
	quiet  bool
}

// NewConfig initializes a config with supplied options, or defaults.
func NewConfig(opts ...func(*Config)) *Config {
	c := new(Config)
	for _, o := range opts {
		o(c)
	}
	if c.bvar == nil {
		OptBindvar(bindvar.New(c.driver))(c)
	}
	if c.tpl == nil {
		OptTemplate(template.New())(c)
	}
	return c
}

// OptDriver sets the driver name.
func OptDriver(s string) func(c *Config) {
	return func(c *Config) {
		c.driver = s
	}
}

// OptTemplate sets the template executer.
func OptTemplate(e template.Executer) func(c *Config) {
	return func(c *Config) {
		c.tpl = e
	}
}

// OptBindvar sets the bindvar parser.
func OptBindvar(p bindvar.Parser) func(c *Config) {
	return func(c *Config) {
		c.bvar = p
	}
}

// OptQuiet disables logging.
func OptQuiet() func(c *Config) {
	return func(c *Config) {
		c.quiet = true
	}
}
