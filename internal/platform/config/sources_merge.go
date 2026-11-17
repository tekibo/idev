package config

import "sort"

// SourceConflict is one field whose user value differs from a changed default.
type SourceConflict struct {
	Field    string `json:"field"`
	Class    string `json:"class"`
	Severity string `json:"severity,omitempty"`
	Current  string `json:"current"`
	Base     string `json:"base"`
	Default  string `json:"default"`
	Required string `json:"required,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

func mergeSource(base, local, changed map[string]any) (map[string]any, []SourceConflict) {
	keys := map[string]bool{}
	for key := range base {
		keys[key] = true
	}
	for key := range local {
		keys[key] = true
	}
	for key := range changed {
		keys[key] = true
	}
	merged := map[string]any{}
	conflicts := []SourceConflict{}
	for key := range keys {
		baseValue, baseOK := base[key]
		localValue, localOK := local[key]
		newValue, newOK := changed[key]
		switch {
		case !localOK:
			if newOK {
				merged[key] = newValue
			}
		case !baseOK:
			merged[key] = localValue
		case sourceEqual(localValue, baseValue):
			if newOK {
				merged[key] = newValue
			} else {
				merged[key] = localValue
			}
		case !newOK || sourceEqual(newValue, baseValue):
			merged[key] = localValue
		case sourceEqual(localValue, newValue):
			merged[key] = localValue
		default:
			merged[key] = localValue
			conflicts = append(conflicts, sourceConflict(key, localValue, baseValue, newValue))
		}
	}
	sort.Slice(conflicts, func(i, j int) bool { return conflicts[i].Field < conflicts[j].Field })
	return merged, conflicts
}

func sourceConflict(field string, current, base, next any) SourceConflict {
	conflict := SourceConflict{
		Field:   field,
		Class:   pinEditable,
		Current: sourceValue(current),
		Base:    sourceValue(base),
		Default: sourceValue(next),
	}
	if spec, ok := sourceSpec(field); ok {
		conflict.Class = spec.Class
		conflict.Severity = spec.Severity
		conflict.Reason = spec.Reason
		if spec.Class == pinRequired {
			conflict.Required = sourceValue(next)
		}
	}
	return conflict
}

func cloneSource(value map[string]any) map[string]any {
	out := make(map[string]any, len(value))
	for key, item := range value {
		out[key] = item
	}
	return out
}
