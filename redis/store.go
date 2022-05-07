package redis

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
}
type Store struct {
	opts  *options
	mutex sync.Mutex
	keys  map[string]bool
	done  chan struct{}
	ch    chan string
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
		keys: make(map[string]bool),
	}
	if s.opts.timeout > time.Second {
		s.done = make(chan struct{})
		s.ch = make(chan string)
		for i := 0; i < 5; i++ {
			go s.work()
		}
		runtime.SetFinalizer(s, (*Store).stop)
	}
	return
}
func (s *Store) stop() {
	close(s.done)
}
func (s *Store) Put(key string, value []byte) error {
	return s.opts.write.Set(s.opts.ctx, key, value, s.opts.timeout).Err()
}

func (s *Store) Get(key string) ([]byte, error) {
	b, e := s.opts.read.Get(s.opts.ctx, key).Bytes()
	if e == redis.Nil {
		return nil, nil
	} else if e != nil {
		return nil, e
	}
	if s.ch != nil {
		s.send(key)
	}
	return b, nil
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
func (s *Store) send(key string) {
	s.mutex.Lock()
	if _, ok := s.keys[key]; !ok {
		s.keys[key] = true
		select {
		case s.ch <- key:
		case <-s.done:
		default:
			go s.ttl(key)
		}
	}
	s.mutex.Unlock()
}

func (s *Store) work() {
	for {
		select {
		case <-s.done:
			return
		case key := <-s.ch:
			s.ttl(key)
		}
	}
}
func (s *Store) ttl(key string) {
	s.opts.write.Expire(s.opts.ctx, key, s.opts.timeout)
	s.mutex.Lock()
	delete(s.keys, key)
	s.mutex.Unlock()
}
