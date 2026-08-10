package taskrunner

import (
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	"github.com/blocknextai/go-packages/cast"
)

type DataProcessor struct {
	outputStore *OutputStore
}

func NewDataProcessor(outputStore *OutputStore) *DataProcessor {
	return &DataProcessor{
		outputStore: outputStore,
	}
}

var (
	methodCallRegex  = regexp.MustCompile(`\$([\w.]+_[\w-]+)\.(first|last|get)\((\d*)\)\.([a-zA-Z0-9_.]+)`)
	arrayAccessRegex = regexp.MustCompile(`\$([\w.]+_[\w-]+)(?:\[(\d+|\*)\])?\.([a-zA-Z0-9_.]+)`)
	triggerVarRegex  = regexp.MustCompile(`\$trigger\.([\w.]+)`)
	inputRegex       = regexp.MustCompile(`\$input\.([a-zA-Z0-9_.]+)`)
)

func (p *DataProcessor) ProcessNodeData(data map[string]any, itemIndex int, triggerData map[string]any, view *BranchView) map[string]any {
	processed := make(map[string]any, len(data))

	for k, v := range data {
		switch value := v.(type) {
		case string:
			processed[k] = p.processStringValue(value, itemIndex, triggerData, view)
		case map[string]any:
			processed[k] = p.ProcessNodeData(value, itemIndex, triggerData, view)
		default:
			processed[k] = v
		}
	}
	return processed
}

func (p *DataProcessor) processStringValue(value string, currentItemIndex int, triggerData map[string]any, view *BranchView) any {
	resolvedStr := value

	if triggerData != nil {
		resolvedStr = p.resolveTriggerVariables(resolvedStr, triggerData)
	}

	resolvedStr = p.resolveInputReferences(resolvedStr, currentItemIndex, view)

	resolvedStr = p.resolveMethodCalls(resolvedStr, view)

	resolvedStr = p.resolveArrayAccess(resolvedStr, currentItemIndex, view)

	return resolvedStr
}

func (p *DataProcessor) resolveTriggerVariables(value string, triggerData map[string]any) string {
	return triggerVarRegex.ReplaceAllStringFunc(value, func(match string) string {
		path := match[len("$trigger."):]
		resolved := p.getNestedValueWithArrayAccess(triggerData, path)
		if resolved == nil {
			return match
		}
		return cast.ToString(resolved)
	})
}

func (p *DataProcessor) resolveInputReferences(value string, currentItemIndex int, view *BranchView) string {
	if p.outputStore == nil || view.Input() == "" {
		return value
	}

	matches := view.Matches("input", value, func(template string) [][]string {
		return inputRegex.FindAllStringSubmatch(template, -1)
	})
	if len(matches) == 0 {
		return value
	}

	inputOutputs, ok := p.outputStore.Get(view.Task(), view.Input())
	if !ok || len(inputOutputs) == 0 {
		return value
	}

	itemIndex := currentItemIndex
	if itemIndex >= len(inputOutputs) {
		itemIndex = 0
	}

	resolvedStr := value
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		fieldPath := match[1]
		resolvedValue := p.getNestedValueWithArrayAccess(inputOutputs[itemIndex], fieldPath)
		if resolvedValue == nil || p.isComplexType(resolvedValue) {
			continue
		}

		placeholder := "$input." + fieldPath
		replacementValue := cast.ToString(resolvedValue)
		if strings.TrimSpace(value) != placeholder {
			var b strings.Builder
			b.WriteString(`"`)
			b.WriteString(replacementValue)
			b.WriteString(`"`)
			replacementValue = b.String()
		}

		resolvedStr = strings.ReplaceAll(resolvedStr, placeholder, replacementValue)
	}

	return resolvedStr
}

func (p *DataProcessor) resolveMethodCalls(value string, view *BranchView) string {
	if p.outputStore == nil {
		return value
	}

	matches := view.Matches("method", value, func(template string) [][]string {
		return methodCallRegex.FindAllStringSubmatch(template, -1)
	})
	resolvedStr := value

	for _, match := range matches {
		if len(match) < 5 {
			continue
		}

		nodeKey := match[1]
		method := match[2]
		methodArg := match[3]
		fieldPath := match[4]

		nodeOutputs, ok := p.outputStore.Get(view.Task(), view.StoreKey(nodeKey))
		if !ok || len(nodeOutputs) == 0 {
			continue
		}

		targetResult := p.getTargetResult(nodeOutputs, method, methodArg)
		if targetResult == nil {
			continue
		}

		resolvedValue := p.getNestedValueWithArrayAccess(targetResult, fieldPath)
		if resolvedValue == nil || p.isComplexType(resolvedValue) {
			continue
		}

		placeholder := p.buildMethodCallPlaceholder(nodeKey, method, methodArg, fieldPath)

		replacementValue := cast.ToString(resolvedValue)
		if strings.TrimSpace(value) != placeholder {
			var b strings.Builder
			b.WriteString(`"`)
			b.WriteString(replacementValue)
			b.WriteString(`"`)
			replacementValue = b.String()
		}

		resolvedStr = strings.ReplaceAll(resolvedStr, placeholder, replacementValue)
	}

	return resolvedStr
}

func (p *DataProcessor) resolveArrayAccess(value string, currentItemIndex int, view *BranchView) string {
	if p.outputStore == nil {
		return value
	}

	matches := view.Matches("array", value, func(template string) [][]string {
		return arrayAccessRegex.FindAllStringSubmatch(template, -1)
	})
	resolvedStr := value

	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		nodeKey := match[1]
		indexStr := match[2]
		fieldPath := match[3]

		nodeOutputs, ok := p.outputStore.Get(view.Task(), view.StoreKey(nodeKey))
		if !ok {
			slog.Warn("reference has no stored outputs to resolve against",
				"component", "data_processor",
				"task_id", view.Task(),
				"node_key", nodeKey)
			continue
		}
		if len(nodeOutputs) == 0 {
			continue
		}

		resolvedValue := p.resolveByIndex(nodeOutputs, nodeKey, indexStr, fieldPath, view.ItemIndex(nodeKey, currentItemIndex), view)
		if resolvedValue == nil || p.isComplexType(resolvedValue) {
			continue
		}

		placeholder := p.buildArrayAccessPlaceholder(nodeKey, indexStr, fieldPath)

		replacementValue := cast.ToString(resolvedValue)
		if strings.TrimSpace(value) != placeholder {
			var b strings.Builder
			b.WriteString(`"`)
			b.WriteString(replacementValue)
			b.WriteString(`"`)
			replacementValue = b.String()
		}

		resolvedStr = strings.ReplaceAll(resolvedStr, placeholder, replacementValue)
	}

	return resolvedStr
}

func (p *DataProcessor) getTargetResult(nodeOutputs []map[string]any, method, methodArg string) map[string]any {
	switch method {
	case "first":
		if len(nodeOutputs) > 0 {
			return nodeOutputs[0]
		}
	case "last":
		if len(nodeOutputs) > 0 {
			return nodeOutputs[len(nodeOutputs)-1]
		}
	case "get":
		if methodArg != "" {
			if index, err := strconv.Atoi(methodArg); err == nil && index >= 0 && index < len(nodeOutputs) {
				return nodeOutputs[index]
			}
		}
	}
	return nil
}

func (p *DataProcessor) resolveByIndex(
	nodeOutputs []map[string]any,
	nodeKey, indexStr, fieldPath string,
	currentItemIndex int,
	view *BranchView,
) any {
	switch indexStr {
	case "":
		if currentItemIndex < len(nodeOutputs) {
			return p.getNestedValueWithArrayAccess(nodeOutputs[currentItemIndex], fieldPath)
		}
		if len(nodeOutputs) > 0 {
			return p.getNestedValueWithArrayAccess(nodeOutputs[0], fieldPath)
		}
	case "*":
		cacheKey := nodeKey + "[*]." + fieldPath
		if cached, ok := view.Collected(cacheKey); ok {
			if len(cached) > 0 {
				return cached
			}
			return nil
		}

		allValues := make([]any, 0, len(nodeOutputs))
		for _, result := range nodeOutputs {
			if val := p.getNestedValueWithArrayAccess(result, fieldPath); val != nil {
				allValues = append(allValues, val)
			}
		}
		view.Collect(cacheKey, allValues)

		if len(allValues) > 0 {
			return allValues
		}
	default:
		if index, err := strconv.Atoi(indexStr); err == nil && index >= 0 && index < len(nodeOutputs) {
			return p.getNestedValueWithArrayAccess(nodeOutputs[index], fieldPath)
		}
	}
	return nil
}

func (p *DataProcessor) isComplexType(value any) bool {
	switch value.(type) {
	case map[string]any, []any, []map[string]any:
		return true
	default:
		return false
	}
}

func (p *DataProcessor) buildMethodCallPlaceholder(nodeKey, method, methodArg, fieldPath string) string {
	var builder strings.Builder
	builder.WriteString("$")
	builder.WriteString(nodeKey)
	builder.WriteString(".")
	builder.WriteString(method)
	builder.WriteString("(")
	if methodArg != "" {
		builder.WriteString(methodArg)
	}
	builder.WriteString(").")
	builder.WriteString(fieldPath)
	return builder.String()
}

func (p *DataProcessor) buildArrayAccessPlaceholder(nodeKey, indexStr, fieldPath string) string {
	var builder strings.Builder
	builder.WriteString("$")
	builder.WriteString(nodeKey)
	if indexStr != "" {
		builder.WriteString("[")
		builder.WriteString(indexStr)
		builder.WriteString("]")
	}
	builder.WriteString(".")
	builder.WriteString(fieldPath)
	return builder.String()
}

func (p *DataProcessor) getNestedValueWithArrayAccess(data any, path string) any {
	parts := p.splitPathPreservingParentheses(path)
	if len(parts) == 0 {
		return data
	}

	current := data
	part := parts[0]

	if strings.Contains(part, "(") && strings.Contains(part, ")") {
		current = p.getArrayAccessValue(current, part)
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
		return p.getFirstElementIfSlice(current)
	}

	return p.getNestedValueWithArrayAccess(current, strings.Join(parts[1:], "."))
}

func (p *DataProcessor) getFirstElementIfSlice(current any) any {
	if arr, ok := current.([]any); ok && len(arr) > 0 {
		return arr[0]
	}
	if arr, ok := current.([]map[string]any); ok && len(arr) > 0 {
		return arr[0]
	}
	return current
}

func (p *DataProcessor) splitPathPreservingParentheses(path string) []string {
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

func (p *DataProcessor) getArrayAccessValue(data any, accessor string) any {
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
		if len(slice) == 0 {
			return nil
		}
		return slice[0]
	case "last":
		if len(slice) == 0 {
			return nil
		}
		return slice[len(slice)-1]
	default:
		return nil
	}
}
