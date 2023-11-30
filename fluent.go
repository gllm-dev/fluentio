package fluentio

import (
	"encoding/json"
	"io"

	"github.com/fluent/fluent-logger-golang/fluent"
)

// Writer is an io.Writer that writes to fluentd.
type Writer struct {
	config        *Config
	client        *fluent.Fluent
	tag           string
	discardWrites bool
}

var _ io.WriteCloser = (*Writer)(nil)

// New creates a new Writer.
// It accepts a variadic number of options that can be used to configure the Writer.
// If no options are provided, it will return an error.
func New(opts ...func(config *Config)) (*Writer, error) {
	config := new(Config)
	for _, opt := range opts {
		opt(config)
	}

	var cfg *fluent.Config
	if config.basicConfig != nil {
		cfg = &fluent.Config{
			FluentHost:         config.basicConfig.FluentHost,
			FluentPort:         config.basicConfig.FluentPort,
			SubSecondPrecision: config.basicConfig.Milliseconds,
		}
	} else if config.fluentConfig != nil {
		cfg = config.fluentConfig
	} else {
		return nil, ErrNoConfigProvided
	}

	client, err := fluent.New(*cfg)
	if err != nil {
		return nil, err
	}

	var tag string
	if config.tag != "" {
		tag = config.tag
	}

	return &Writer{
		config:        config,
		client:        client,
		tag:           tag,
		discardWrites: false,
	}, nil
}

// Write is the implementation of io.Writer.
func (f *Writer) Write(p []byte) (n int, err error) {
	// If the Writer has been closed,
	if f.discardWrites {
		return len(p), nil
	}

	var m map[string]interface{}
	err = json.Unmarshal(p, &m)
	if err != nil {
		return 0, err
	}

	err = f.client.Post(f.tag, m)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

// Close is the implementation of io.Closer.
func (f *Writer) Close() error {
	if f.config.discardWritesAfterClose {
		f.discardWrites = true
	}
	return f.client.Close()
}
