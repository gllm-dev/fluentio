package fluentio

import "github.com/fluent/fluent-logger-golang/fluent"

// Config is used to configure the Writer.
type Config struct {
	fluentConfig            *fluent.Config
	basicConfig             *BasicConfig
	tag                     string
	discardWritesAfterClose bool
}

// BasicConfig is used to configure the Writer with a basic configuration.
type BasicConfig struct {
	FluentHost   string
	FluentPort   int
	Milliseconds bool
}

// WithBasicConfig returns a function that can be used to configure the Writer with a basic configuration.
func WithBasicConfig(host string, port int, milliseconds bool) func(*Config) {
	return func(c *Config) {
		c.basicConfig = &BasicConfig{
			FluentHost:   host,
			FluentPort:   port,
			Milliseconds: milliseconds,
		}
	}
}

// WithFluentConfig returns a function that can be used to configure the Writer with the standard fluent-logger-golang configuration.
func WithFluentConfig(config *fluent.Config) func(*Config) {
	return func(c *Config) {
		c.fluentConfig = config
	}
}

// WithTag returns a function that can be used to configure the Writer with a tag.
func WithTag(tag string) func(*Config) {
	return func(c *Config) {
		c.tag = tag
	}
}

// WithDiscardWritesAfterClose returns a function that can be used to configure the Writer to discard writes after Close()
func WithDiscardWritesAfterClose() func(*Config) {
	return func(c *Config) {
		c.discardWritesAfterClose = true
	}
}
