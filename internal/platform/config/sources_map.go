package config

import (
	"encoding/json"
	"fmt"
)

func flattenSource(prefix string, value any, out map[string]any) {
	if object, ok := value.(map[string]any); ok {
		for key, child := range object {
			name := key
			if prefix != "" {
				name = prefix + "." + key
			}
			flattenSource(name, child, out)
		}
		return
	}
	out[prefix] = value
}

func flatSource(value map[string]any) map[string]any {
	out := map[string]any{}
	flattenSource("", value, out)
	return out
}

func unflatSource(flat map[string]any) map[string]any {
	out := map[string]any{}
	for name, value := range flat {
		parts := splitSourcePath(name)
		node := out
		for _, part := range parts[:len(parts)-1] {
			next, ok := node[part].(map[string]any)
			if !ok {
				next = map[string]any{}
				node[part] = next
			}
			node = next
		}
		node[parts[len(parts)-1]] = value
	}
	return out
}

func splitSourcePath(name string) []string {
	parts := []string{}
	start := 0
	for index, char := range name {
		if char == '.' {
			parts = append(parts, name[start:index])
			start = index + 1
		}
	}
	return append(parts, name[start:])
}

func sourceValue(value any) string {
	if text, ok := value.(string); ok {
		return text
	}
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprint(value)
	}
	return string(data)
}

func sourceEqual(left, right any) bool {
	return sourceValue(left) == sourceValue(right)
}

func sourceMapEqual(left, right map[string]any) bool {
	leftFlat := flatSource(left)
	rightFlat := flatSource(right)
	if len(leftFlat) != len(rightFlat) {
		return false
	}
	for key, value := range leftFlat {
		if !sourceEqual(value, rightFlat[key]) {
			return false
		}
	}
	return true
}
