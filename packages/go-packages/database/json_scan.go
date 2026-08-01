package database

import (
	"github.com/blocknextai/go-packages/json"
)

// UnmarshalJSONColumn decodes a JSON column into target, treating a nil byte
// slice as a no-op so NULL columns leave target unchanged.
func UnmarshalJSONColumn[T any](data []byte, target *T) error {
	if data == nil {
		return nil
	}
	return json.Unmarshal(data, target)
}
