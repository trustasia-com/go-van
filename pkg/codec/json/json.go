// Package json provides ...
package json

import (
	"encoding/json"

	"github.com/trustasia-com/go-van/pkg/codec"
)

// NewCodec new encoder & decoder
func NewCodec() codec.Codec {
	return jsonCodec{}
}

// jsonCodec yaml codec
type jsonCodec struct{}

// Marshal returns the wire format of v.
func (jsonCodec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal parses the wire format into v.
func (jsonCodec) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
