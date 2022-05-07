package redis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd
}
type Store struct {
	opts *options
}

func New(redis Redis, opt ...Option) (s *Store, e error) {
	opts := defaultOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	if redis == nil {
		e = errors.New(`redis not supported nil`)
		return
	}
	opts.write = redis
	if opts.read == nil {
		opts.read = redis
	}
	s = &Store{
		opts: &opts,
	}
	return
}
func (s *Store) Put(key string, value []byte) error {
	return s.opts.write.Set(s.opts.ctx, key, value, s.opts.timeout).Err()
}

func (s *Store) Get(key string) ([]byte, error) {
	return s.opts.read.Get(s.opts.ctx, key).Bytes()
}

func (s *Store) Del(key string) error {
	return s.opts.write.Del(s.opts.ctx, key).Err()
}

func (s *Store) DelPrefix(prefix string) error {
	var (
		cursor uint64
		count  int64 = 1000
		match        = prefix + "*"
	)
	for {
		scan := s.opts.write.Scan(s.opts.ctx, cursor,
			match,
			count,
		)
		var (
			keys []string
			e    error
		)
		keys, cursor, e = scan.Result()
		if e != nil {
			return e
		}
		if len(keys) != 0 {
			e = s.opts.write.Del(s.opts.ctx, keys...).Err()
			if e != nil {
				return e
			}
		}
		if cursor == 0 {
			break
		}
	}
	return nil
}
