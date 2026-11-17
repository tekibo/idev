package config

import (
	"encoding/json"
	"sort"
)

const (
	pinRequired = "required"
	pinEditable = "editable"
	pinWarn     = "warn"
	pinBlock    = "block"
)

type sourcePinSpec struct {
	Field    string
	Class    string
	Severity string
	Reason   string
}

var sourceSpecs = []sourcePinSpec{
	{"release.manifest", pinEditable, pinWarn, "release source"},
	{"swift.version", pinRequired, pinWarn, "tested Swift toolchain"},
	{"swift.amd64.url", pinEditable, pinWarn, "Swift archive mirror"},
	{"swift.amd64.sha256", pinRequired, pinWarn, "tested Swift archive"},
	{"swift.arm64.url", pinEditable, pinWarn, "Swift archive mirror"},
	{"swift.arm64.sha256", pinRequired, pinWarn, "tested Swift archive"},
}

// SourceWarning is one unsupported required pin.
type SourceWarning struct {
	Field    string `json:"field"`
	Current  string `json:"current"`
	Required string `json:"required"`
	Reason   string `json:"reason"`
	Block    bool   `json:"block,omitempty"`
}

// SourceWarnings compares the effective sources with the required pins.
func SourceWarnings(s Sources) []SourceWarning {
	defaults, err := defaultSourceMap()
	if err != nil {
		return nil
	}
	current, err := sourceMapFromStruct(s)
	if err != nil {
		return nil
	}
	return sourceWarnings(flatSource(current), flatSource(defaults))
}

func sourceWarnings(current, defaults map[string]any) []SourceWarning {
	warnings := []SourceWarning{}
	for _, spec := range sourceSpecs {
		if spec.Class != pinRequired {
			continue
		}
		required, ok := defaults[spec.Field]
		if !ok || sourceEqual(current[spec.Field], required) {
			continue
		}
		warnings = append(warnings, SourceWarning{
			Field:    spec.Field,
			Current:  sourceValue(current[spec.Field]),
			Required: sourceValue(required),
			Reason:   spec.Reason,
			Block:    spec.Severity == pinBlock,
		})
	}
	sort.Slice(warnings, func(i, j int) bool { return warnings[i].Field < warnings[j].Field })
	return warnings
}

func sourceMapFromStruct(s Sources) (map[string]any, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, err
	}
	return value, nil
}

func sourceSpec(field string) (sourcePinSpec, bool) {
	for _, spec := range sourceSpecs {
		if spec.Field == field {
			return spec, true
		}
	}
	return sourcePinSpec{}, false
}
