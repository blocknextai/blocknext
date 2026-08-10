package taskrunner

import (
	"strconv"
	"strings"

	"github.com/blocknextai/go-packages/cast"
)

func getNestedValueWithArrayAccess(data any, path string) any {
	parts := splitPathPreservingParentheses(path)
	if len(parts) == 0 {
		return data
	}

	current := data
	part := parts[0]

	if strings.Contains(part, "(") && strings.Contains(part, ")") {
		current = getArrayAccessValue(current, part)
	} else {
		currentMap, ok := current.(map[string]any)
		if !ok {
			return nil
		}
		val, ok := currentMap[part]
		if !ok {
			return nil
		}
		current = val
	}

	if len(parts) == 1 {
		return getFirstElementIfSlice(current)
	}

	return getNestedValueWithArrayAccess(current, strings.Join(parts[1:], "."))
}

func getFirstElementIfSlice(current any) any {
	if arr, ok := current.([]any); ok && len(arr) > 0 {
		return arr[0]
	}
	if arr, ok := current.([]map[string]any); ok && len(arr) > 0 {
		return arr[0]
	}
	return current
}

func splitPathPreservingParentheses(path string) []string {
	var parts []string
	var current strings.Builder
	parenCount := 0

	for _, char := range path {
		switch char {
		case '(':
			parenCount++
		case ')':
			parenCount--
		}

		if char == '.' && parenCount == 0 {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

func getArrayAccessValue(data any, accessor string) any {
	methodStart := strings.Index(accessor, "(")
	methodEnd := strings.Index(accessor, ")")

	if methodStart == -1 || methodEnd == -1 {
		return nil
	}

	methodName := accessor[:methodStart]
	argsStr := accessor[methodStart+1 : methodEnd]

	slice := cast.ToSlice(data)
	if len(slice) == 0 {
		return nil
	}

	switch methodName {
	case "get":
		index, err := strconv.Atoi(argsStr)
		if err != nil || index < 0 || index >= len(slice) {
			return nil
		}
		return slice[index]
	case "first":
		return slice[0]
	case "last":
		return slice[len(slice)-1]
	default:
		return nil
	}
}

func isComplexType(value any) bool {
	switch value.(type) {
	case map[string]any, []any, []map[string]any:
		return true
	default:
		return false
	}
}
