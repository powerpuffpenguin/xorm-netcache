package xormcache

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var defaultCoder = NewJsonCoder()

func DefaultCoder() *JsonCoder {
	return defaultCoder
}

type JsonCoder struct {
	keys   map[string]reflect.Type
	mutext sync.RWMutex
}

func NewJsonCoder() *JsonCoder {
	return &JsonCoder{
		keys: make(map[string]reflect.Type),
	}
}
func (c *JsonCoder) Encode(key string, data interface{}) ([]byte, error) {
	val, e := json.Marshal(data)
	if e != nil {
		return nil, e
	}

	t := reflect.TypeOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	c.mutext.Lock()
	if val, ok := c.keys[key]; ok {
		if val != t {
			return nil, fmt.Errorf(`encode(%s:%s) type not match exists type is %s`, key, t.Name(), val.Name())
		}
	} else {
		c.keys[key] = t
	}
	c.mutext.Unlock()

	return val, nil
}
func (c *JsonCoder) Decode(key string, data []byte) (interface{}, error) {
	c.mutext.RLock()
	t, ok := c.keys[key]
	c.mutext.RUnlock()
	if !ok {
		return nil, errors.New(`decode unknow type of ` + key)
	}

	n := reflect.New(t)
	p := n.Interface()
	e := json.Unmarshal(data, p)
	if e != nil {
		return nil, e
	}
	return p, nil
}
