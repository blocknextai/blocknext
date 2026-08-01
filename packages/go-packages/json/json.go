// Package json wraps the standard encoding/json package and provides helpers
// for converting between arbitrary values via JSON.
package json

import (
	stdjson "encoding/json"
	"errors"
)

// RawMessage is an alias for the standard library encoding/json.RawMessage.
type RawMessage = stdjson.RawMessage

// Number is an alias for the standard library encoding/json.Number.
type Number = stdjson.Number

var (
	// ErrMarshalingArgs indicates that the source value could not be marshaled to JSON.
	ErrMarshalingArgs = errors.New("marshaling args")
	// ErrUnmarshalingRequest indicates that the JSON could not be unmarshaled into the target.
	ErrUnmarshalingRequest = errors.New("unmarshaling request")
)

// Marshal returns the JSON encoding of v.
func Marshal(v any) ([]byte, error) {
	return stdjson.Marshal(v)
}

// Unmarshal parses the JSON-encoded data and stores the result in the value
// pointed to by v.
func Unmarshal(data []byte, v any) error {
	return stdjson.Unmarshal(data, v)
}

// ArgsToStruct converts args into response by marshaling args to JSON and
// unmarshaling it back into response.
func ArgsToStruct(args any, response any) error {
	argsJSON, err := Marshal(args)
	if err != nil {
		return ErrMarshalingArgs
	}

	if err := Unmarshal(argsJSON, &response); err != nil {
		return ErrUnmarshalingRequest
	}

	return nil
}
