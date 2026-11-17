package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// XcodeSelection records the current Xcode.xip basename.
type XcodeSelection struct {
	XIP string `json:"xip"`
}

// LoadXcodeSelection reads the current xip selection.
func LoadXcodeSelection(home string) (XcodeSelection, error) {
	path := filepath.Join(home, "config", "xcode.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return XcodeSelection{}, nil
		}
		return XcodeSelection{}, err
	}
	var selection XcodeSelection
	if err := json.Unmarshal(data, &selection); err != nil {
		return XcodeSelection{}, err
	}
	return selection, nil
}

// SaveXcodeSelection writes the current xip selection.
func SaveXcodeSelection(home string, selection XcodeSelection) error {
	path := filepath.Join(home, "config", "xcode.json")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(selection, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
