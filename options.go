package xormcache

import (
	"log"
	"os"
)

var defaultOptions = options{
	logger: log.New(os.Stdout, `[cache]`, log.LstdFlags),
	prefix: `cache`,
	sep:    `-`,
	coder:  GobCoder{},
}

type options struct {
	logger *log.Logger
	prefix string
	sep    string
	coder  Coder
}

type Option interface {
	apply(*options)
}

type funcOption struct {
	f func(*options)
}

func (fdo *funcOption) apply(do *options) {
	fdo.f(do)
}
func newFuncOption(f func(*options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func WithCoder(coder Coder) Option {
	return newFuncOption(func(o *options) {
		if coder == nil {
			o.coder = GobCoder{}
		} else {
			o.coder = coder
		}
	})
}
func WithKeyPrefix(prefix string) Option {
	return newFuncOption(func(o *options) {
		o.prefix = prefix
	})
}
func WithKeySeparators(sep string) Option {
	return newFuncOption(func(o *options) {
		if sep == `` {
			sep = `-`
		}
		o.sep = sep
	})
}
func WithLogger(logger *log.Logger) Option {
	return newFuncOption(func(o *options) {
		o.logger = logger
	})
}
