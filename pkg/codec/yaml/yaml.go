// Package yaml provides ...
package yaml

import (
	"github.com/trustasia-com/go-van/pkg/codec"

	"gopkg.in/yaml.v3"
)

// NewCodec new encoder & decoder
func NewCodec() codec.Codec {
	return yamlCodec{}
}

// yamlCodec yaml codec
type yamlCodec struct{}

// Marshal returns the wire format of v.
func (yamlCodec) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

// Unmarshal parses the wire format into v.
func (yamlCodec) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}
