// Package yaml provides ...
package yaml

import (
	"github.com/deepzz0/go-van/pkg/codec"

	"gopkg.in/yaml.v3"
)

// NewCodec new encoder & decoder
func NewCodec() codec.Codec {
	return Codec{}
}

// Codec yaml codec
type Codec struct{}

// Marshal returns the wire format of v.
func (_ Codec) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

// Unmarshal parses the wire format into v.
func (_ Codec) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}
