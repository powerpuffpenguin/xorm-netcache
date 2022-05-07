package xormcache

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

// How to encode and decode cached data to interact with backend storage devices
type Coder interface {
	Encode(data interface{}) ([]byte, error)
	Decode(data []byte, to interface{}) error
}

type JsonCoder struct{}

func (c JsonCoder) Encode(data interface{}) ([]byte, error) {
	val, e := json.Marshal(data)
	if e != nil {
		return nil, e
	}
	return val, nil
}
func (c JsonCoder) Decode(data []byte, to interface{}) error {
	return json.Unmarshal(data, to)
}

type GobCoder struct{}

func (c GobCoder) Encode(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	e := enc.Encode(&data)
	if e != nil {
		return nil, e
	}
	return buf.Bytes(), nil
}
func (c GobCoder) Decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}
