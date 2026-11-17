package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// FlowSelection is the saved setup group selection.
type FlowSelection struct {
	Flows []string `json:"flows"`
	Expo  bool     `json:"expo"`
}

// LoadFlows reads the saved setup selection.
func LoadFlows(home string) (FlowSelection, error) {
	selection := FlowSelection{Expo: true}
	path := filepath.Join(home, "config", "flows.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return selection, nil
		}
		return FlowSelection{}, err
	}
	if err := json.Unmarshal(data, &selection); err != nil {
		return FlowSelection{}, err
	}
	return selection, nil
}

// SaveFlows writes the setup selection.
func SaveFlows(home string, selection FlowSelection) error {
	path := filepath.Join(home, "config", "flows.json")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(selection, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
