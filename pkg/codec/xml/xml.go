// Package xml provides ...
package xml

import (
	"encoding/xml"

	"github.com/trustasia-com/go-van/pkg/codec"
)

// NewCodec new encoder & decoder
func NewCodec() codec.Codec {
	return xmlCodec{}
}

// xmlCodec yaml codec
type xmlCodec struct{}

// Marshal returns the wire format of v.
func (xmlCodec) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

// Unmarshal parses the wire format into v.
func (xmlCodec) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}
