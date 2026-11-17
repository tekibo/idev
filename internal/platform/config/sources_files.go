package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func SourcePath(home string) string {
	if value := os.Getenv("IDEV_SOURCES"); value != "" {
		return value
	}
	return filepath.Join(home, "config", "sources.json")
}

func AppliedPath(home string) string {
	return filepath.Join(home, "config", "sources.applied.json")
}

func loadSourceMap(path string) (map[string]any, bool, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, false, fmt.Errorf("%s: %w", path, err)
	}
	return value, true, nil
}

func defaultSourceMap() (map[string]any, error) {
	var value map[string]any
	if err := json.Unmarshal(defaultSources, &value); err != nil {
		return nil, err
	}
	return value, nil
}

func writeSourceMap(path string, value map[string]any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func backupSource(path string) (string, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	backup := path + ".bak." + time.Now().UTC().Format("20060102-150405")
	if err := os.WriteFile(backup, data, 0644); err != nil {
		return "", err
	}
	return backup, nil
}
