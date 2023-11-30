package jsonldhelpers

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
	"github.com/piprate/json-gold/ld"
)

// SingleValue represents a subset of an object, that contains that object's
// @value field
type SingleValue[T any] struct {
	Value T `mapstructure:"@value"`
}

// Decoder is what is used for decoding a JSON-LD document into a Go struct.
type Decoder struct {
	Processor *ld.JsonLdProcessor
	Options   *ld.JsonLdOptions
}

var decoder Decoder

func init() {
	decoder = Decoder{
		Processor: ld.NewJsonLdProcessor(),
		Options:   ld.NewJsonLdOptions(""),
	}
}

// Decode decodes a JSON-LD document into a Go struct.
func (d Decoder) Decode(source []byte, destination any) error {
	var objOutput any
	err := json.Unmarshal(source, &objOutput)
	if err != nil {
		return err
	}

	expanded, err := d.Processor.Expand(objOutput, d.Options)
	if err != nil {
		return err
	}

	return mapstructure.Decode(expanded, destination)
}

// Decode decodes a JSON-LD document into a Go struct.
func Decode(source []byte, destination any) error {
	return decoder.Decode(source, &destination)
}
