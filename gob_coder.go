package xormcache

import (
	"bytes"
	"encoding/gob"
)

type GobCoder struct{}

func (c GobCoder) Encode(key string, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	e := enc.Encode(&data)
	if e != nil {
		return nil, e
	}
	return buf.Bytes(), nil
}
func (c GobCoder) Decode(key string, data []byte) (interface{}, error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	var to interface{}
	e := dec.Decode(to)
	if e != nil {
		return nil, e
	}
	return to, nil
}
