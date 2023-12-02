package jsonldhelpers

import (
	"fediverse/jsonhelpers"
	"fediverse/pointers"

	"github.com/mitchellh/mapstructure"
	"github.com/piprate/json-gold/ld"
)

// ValueNode represents a subset of an object, that contains that object's
// @value field
type ValueNode[T any] struct {
	Value T `mapstructure:"@value"`
}

type IDNode struct {
	ID string `mapstructure:"@id"`
}

// Decoder is what is used for decoding a JSON-LD document into a Go struct.
type Decoder struct {
	processor *ld.JsonLdProcessor
	options   *ld.JsonLdOptions
}

var decoder Decoder

var processor = ld.NewJsonLdProcessor()
var options = ld.NewJsonLdOptions("")

func init() {
	decoder = Decoder{
		processor: processor,
		options:   options,
	}
}

// Decode decodes a JSON-LD document into a Go struct.
func (d Decoder) Decode(source []byte, destination any) error {
	objOutput, err := jsonhelpers.UnmarshalAny(source)
	if err != nil {
		return err
	}

	return d.DecodeValue(objOutput, &destination)
}

// DecodeValue decodes the given deserialized object into to the given
// destination object.
func (d Decoder) DecodeValue(source any, destination any) error {
	expanded, err := pointers.
		ValueOrDefault(d.processor, processor).
		Expand(source, pointers.ValueOrDefault(d.options, options))
	if err != nil {
		return err
	}

	return mapstructure.Decode(expanded, destination)
}

// DecodeBytes decodes a JSON-LD document from a byte slice into a go value
// (ideally a struct that uses mitchelh)
func DecodeBytes(source []byte, destination any) error {
	return decoder.Decode(source, &destination)
}

// DecodeValue decodes the deserialized value (ideally a map[string]any or
// slices of map[stringany]) of a JSON-LD document into the given destination.
func DecodeValue(source any, destination any) error {
	return decoder.DecodeValue(source, &destination)
}
