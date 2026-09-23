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

func substitute(template string, matches [][]string, resolve func(match []string) (string, any)) string {
	resolved := template

	for _, match := range matches {
		placeholder, value := resolve(match)
		if placeholder == "" || value == nil || isComplexType(value) {
			continue
		}

		replacement := cast.ToString(value)
		if strings.TrimSpace(template) != placeholder {
			var builder strings.Builder
			builder.WriteString(`"`)
			builder.WriteString(replacement)
			builder.WriteString(`"`)
			replacement = builder.String()
		}

		resolved = strings.ReplaceAll(resolved, placeholder, replacement)
	}

	return resolved
}

func (p *DataProcessor) resolveTriggerVariables(value string, triggerData map[string]any) string {
	return triggerVarRegex.ReplaceAllStringFunc(value, func(match string) string {
		path := match[len("$trigger."):]
		resolved := getNestedValueWithArrayAccess(triggerData, path)
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

	return substitute(value, matches, func(match []string) (string, any) {
		if len(match) < 2 {
			return "", nil
		}

		fieldPath := match[1]
		return "$input." + fieldPath, getNestedValueWithArrayAccess(inputOutputs[itemIndex], fieldPath)
	})
}

func (p *DataProcessor) resolveMethodCalls(value string, view *BranchView) string {
	if p.outputStore == nil {
		return value
	}

	matches := view.Matches("method", value, func(template string) [][]string {
		return methodCallRegex.FindAllStringSubmatch(template, -1)
	})

	return substitute(value, matches, func(match []string) (string, any) {
		if len(match) < 5 {
			return "", nil
		}

		nodeKey, method, methodArg, fieldPath := match[1], match[2], match[3], match[4]

		nodeOutputs, ok := p.outputStore.Get(view.Task(), view.StoreKey(nodeKey))
		if !ok || len(nodeOutputs) == 0 {
			return "", nil
		}

		targetResult := getTargetResult(nodeOutputs, method, methodArg)
		if targetResult == nil {
			return "", nil
		}

		return buildMethodCallPlaceholder(nodeKey, method, methodArg, fieldPath),
			getNestedValueWithArrayAccess(targetResult, fieldPath)
	})
}

func (p *DataProcessor) resolveArrayAccess(value string, currentItemIndex int, view *BranchView) string {
	if p.outputStore == nil {
		return value
	}

	matches := view.Matches("array", value, func(template string) [][]string {
		return arrayAccessRegex.FindAllStringSubmatch(template, -1)
	})

	return substitute(value, matches, func(match []string) (string, any) {
		if len(match) < 4 {
			return "", nil
		}

		nodeKey, indexStr, fieldPath := match[1], match[2], match[3]

		nodeOutputs, ok := p.outputStore.Get(view.Task(), view.StoreKey(nodeKey))
		if !ok {
			slog.Warn("reference has no stored outputs to resolve against",
				"component", "data_processor",
				"task_id", view.Task(),
				"node_key", nodeKey)
			return "", nil
		}
		if len(nodeOutputs) == 0 {
			return "", nil
		}

		return buildArrayAccessPlaceholder(nodeKey, indexStr, fieldPath),
			resolveByIndex(nodeOutputs, nodeKey, indexStr, fieldPath, view.ItemIndex(nodeKey, currentItemIndex), view)
	})
}

func getTargetResult(nodeOutputs []map[string]any, method, methodArg string) map[string]any {
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

func resolveByIndex(
	nodeOutputs []map[string]any,
	nodeKey, indexStr, fieldPath string,
	currentItemIndex int,
	view *BranchView,
) any {
	switch indexStr {
	case "":
		if currentItemIndex < len(nodeOutputs) {
			return getNestedValueWithArrayAccess(nodeOutputs[currentItemIndex], fieldPath)
		}
		if len(nodeOutputs) > 0 {
			return getNestedValueWithArrayAccess(nodeOutputs[0], fieldPath)
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
			if val := getNestedValueWithArrayAccess(result, fieldPath); val != nil {
				allValues = append(allValues, val)
			}
		}
		view.Collect(cacheKey, allValues)

		if len(allValues) > 0 {
			return allValues
		}
	default:
		if index, err := strconv.Atoi(indexStr); err == nil && index >= 0 && index < len(nodeOutputs) {
			return getNestedValueWithArrayAccess(nodeOutputs[index], fieldPath)
		}
	}
	return nil
}

func buildMethodCallPlaceholder(nodeKey, method, methodArg, fieldPath string) string {
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

func buildArrayAccessPlaceholder(nodeKey, indexStr, fieldPath string) string {
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
