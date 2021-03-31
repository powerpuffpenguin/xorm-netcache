package cache

import (
	"time"

	"xorm.io/xorm/caches"
)

var defaultCacher *caches.LRUCacher

func SetDefaultCacher(cacher *caches.LRUCacher) {
	defaultCacher = cacher
}
func DefaultCacher() *caches.LRUCacher {
	return defaultCacher
}

type Element struct {
	Data    []byte
	Modtime time.Time
}
