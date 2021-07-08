package apollo

import (
	"github.com/zouyx/agollo/v4/env/config"
)

// Option apollo option
type Option func(opts *Options)

// Options apollo options
type Options = config.AppConfig

// WithConfig Set apollo server ip adder
func WithConfig(conf config.AppConfig) Option {
	return func(opts *Options) { *opts = conf }
}
