package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// ApplySourcePlan writes the chosen values and records the applied defaults.
func ApplySourcePlan(plan *SourcePlan, keep map[string]bool) (string, error) {
	result := cloneSource(plan.merged)
	for _, conflict := range plan.Conflicts {
		if keep[conflict.Field] {
			if value, ok := plan.local[conflict.Field]; ok {
				result[conflict.Field] = value
			} else {
				delete(result, conflict.Field)
			}
			continue
		}
		if value, ok := plan.defaults[conflict.Field]; ok {
			result[conflict.Field] = value
		} else {
			delete(result, conflict.Field)
		}
	}
	for _, warning := range plan.Warnings {
		if keep[warning.Field] {
			continue
		}
		if value, ok := plan.defaults[warning.Field]; ok {
			result[warning.Field] = value
		}
	}
	value := unflatSource(result)
	if err := validateSourceMap(value); err != nil {
		return "", err
	}
	backup := ""
	if !sourceMapEqual(unflatSource(plan.local), value) {
		var err error
		backup, err = backupSource(plan.Path)
		if err != nil {
			return "", err
		}
		plan.Backup = backup
	}
	if err := writeSourceMap(plan.Path, value); err != nil {
		return "", err
	}
	if err := writeSourceMap(plan.applied, unflatSource(plan.defaults)); err != nil {
		return "", err
	}
	return backup, nil
}

// RefreshSources resets the user copy to the embedded defaults.
func RefreshSources(home string) (string, error) {
	if os.Getenv("IDEV_SOURCES") != "" {
		return "", errors.New("sources: refresh does not apply to IDEV_SOURCES")
	}
	defaults, err := defaultSourceMap()
	if err != nil {
		return "", err
	}
	path := SourcePath(home)
	backup, err := backupSource(path)
	if err != nil {
		return "", err
	}
	if err := writeSourceMap(path, defaults); err != nil {
		return "", err
	}
	if err := writeSourceMap(AppliedPath(home), defaults); err != nil {
		return "", err
	}
	return backup, nil
}

func validateSourceMap(value map[string]any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	var sources Sources
	if err := json.Unmarshal(data, &sources); err != nil {
		return fmt.Errorf("sources: %w", err)
	}
	if sources.Release.Manifest == "" || sources.Swift.Version == "" {
		return errors.New("sources: release manifest or Swift version is empty")
	}
	return nil
}
