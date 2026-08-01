// Package cast provides helpers for converting arbitrary values to common
// primitive types.
package cast

import (
	"math"
	"strconv"
)

// ToString converts v to its string representation, returning an empty string
// for nil, unsupported types, or non-finite floats.
func ToString(v any) string {
	if v == nil {
		return ""
	}

	switch value := v.(type) {
	case string:
		return value
	case int:
		return strconv.Itoa(value)
	case int32:
		return strconv.FormatInt(int64(value), 10)
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float64:
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return ""
		}
		return strconv.FormatFloat(value, 'f', -1, 64)
	case float32:
		if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
			return ""
		}
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case bool:
		return strconv.FormatBool(value)
	default:
		return ""
	}
}

// ToBool converts v to a bool, treating "true", "1", and "yes" as true for
// strings and non-zero numbers as true.
func ToBool(v any) bool {
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return val == "true" || val == "1" || val == "yes"
	case int, int64:
		return val != 0
	case float64:
		return val != 0
	default:
		return false
	}
}

// ToFloat converts v to a float64, returning 0 for unsupported types or
// unparseable strings.
func ToFloat(v any) float64 {
	switch val := v.(type) {
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case float64:
		return val
	case string:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0
		}
		return f
	default:
		return 0
	}
}

// ToInt64 converts v to an int64, returning 0 for unsupported types or
// unparseable strings.
func ToInt64(v any) int64 {
	switch val := v.(type) {
	case int64:
		return val
	case int:
		return int64(val)
	case float64:
		return int64(val)
	case string:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return 0
		}
		return i
	default:
		return 0
	}
}

// ToSlice converts data to a []any, returning nil for unsupported types.
func ToSlice(data any) []any {
	switch v := data.(type) {
	case []any:
		return v
	case []string:
		slice := make([]any, len(v))
		for i, item := range v {
			slice[i] = item
		}
		return slice
	case []map[string]any:
		slice := make([]any, len(v))
		for i, item := range v {
			slice[i] = item
		}
		return slice
	default:
		return nil
	}
}
