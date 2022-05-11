package xormcache

// How to encode and decode cached bean to interact with backend storage devices
type Coder interface {
	Encode(key string, data interface{}) ([]byte, error)
	Decode(key string, data []byte) (interface{}, error)
}
