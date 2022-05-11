package xormcache

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type BinaryValue interface {
	TypeID() uint32
	Marshal() (data []byte, err error)
}
type BinaryDecoder interface {
	Unmarshal(data []byte) (interface{}, error)
}
type BinaryCoder struct {
	types map[uint32]BinaryDecoder
}

func NewBinaryCoder() *BinaryCoder {
	return &BinaryCoder{
		types: make(map[uint32]BinaryDecoder),
	}
}
func (c *BinaryCoder) Register(key uint32, decoder BinaryDecoder) {
	if _, ok := c.types[key]; ok {
		panic(fmt.Errorf(`decoder of %d already exists`, key))
	}
	c.types[key] = decoder
}
func (c *BinaryCoder) Encode(key string, data interface{}) ([]byte, error) {
	val, ok := data.(BinaryValue)
	if !ok {
		return nil, nil
	}
	tid := val.TypeID()
	_, exists := c.types[tid]
	if !exists {
		return nil, fmt.Errorf(`encoding not found decoder of %d`, tid)
	}

	b, e := val.Marshal()
	if e != nil {
		return nil, e
	}
	result := make([]byte, 4+len(b))
	binary.BigEndian.PutUint32(result, tid)
	copy(result[4:], b)
	return result, nil
}
func (c *BinaryCoder) Decode(key string, data []byte) (interface{}, error) {
	if len(data) < 4 {
		return nil, errors.New("decoding encoded mismatch")
	}
	tid := binary.BigEndian.Uint32(data)
	coder, ok := c.types[tid]
	if !ok {
		return nil, fmt.Errorf("decoding not found decoder of %d", tid)
	}
	return coder.Unmarshal(data[4:])
}
