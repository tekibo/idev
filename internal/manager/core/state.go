package core

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func ComponentsPath(home string) string {
	return filepath.Join(home, "state", "components.json")
}

func LoadComponents(home string) (map[string]string, error) {
	data, err := os.ReadFile(ComponentsPath(home))
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	state := map[string]string{}
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return state, nil
}

func RecordComponents(home string, values map[string]string) error {
	state, err := LoadComponents(home)
	if err != nil {
		return err
	}
	for name, version := range values {
		state[name] = version
	}
	if err := os.MkdirAll(filepath.Dir(ComponentsPath(home)), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ComponentsPath(home), data, 0644)
}
