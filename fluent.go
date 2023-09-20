package fluentio

import (
	"encoding/json"
	"io"

	"github.com/fluent/fluent-logger-golang/fluent"
)

// Writer is an io.Writer that writes to fluentd.
type Writer struct {
	client *fluent.Fluent
	tag    string
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
		client: client,
		tag:    tag,
	}, nil
}

// Write is the implementation of io.Writer.
func (f *Writer) Write(p []byte) (n int, err error) {
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
	return f.client.Close()
}
