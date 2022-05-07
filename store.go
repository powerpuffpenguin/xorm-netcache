package xormcache

type Store interface {
	Put(key string, value []byte) error
	Get(key string) ([]byte, error)
	Del(key string) error
	DelPrefix(prefix string) error
}
