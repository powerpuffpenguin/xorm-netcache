package redis

import (
	"context"
	"time"
)

var defaultOptions = options{
	ctx:          context.Background(),
	timeout:      time.Hour,
	expireWorker: 5,
}

type options struct {
	ctx          context.Context
	write        Redis
	read         Redis
	timeout      time.Duration
	expireWorker int
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
func WithRead(read Redis) Option {
	return newFuncOption(func(o *options) {
		o.read = read
	})
}

func WithContext(ctx context.Context) Option {
	return newFuncOption(func(o *options) {
		if ctx == nil {
			ctx = context.Background()
		}
		o.ctx = ctx
	})
}

// Set the cache expiration time if it is less than 0, it will never expire
func WithTimeout(duration time.Duration) Option {
	return newFuncOption(func(o *options) {
		if duration < 0 {
			duration = 0
		}
		o.timeout = duration
	})
}

// How much to give to the worker to update the expiration time
func WithExpireWorker(count int) Option {
	return newFuncOption(func(o *options) {
		if count < 1 {
			count = 1
		}
		o.expireWorker = count
	})
}
