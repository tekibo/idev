package project

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func pnpmWorkspacePatterns(root string) []string {
	data, err := os.ReadFile(filepath.Join(root, "pnpm-workspace.yaml"))
	if err != nil {
		return nil
	}
	patterns := []string{}
	inPackages := false
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "packages:") {
			inPackages = true
			continue
		}
		if !inPackages {
			continue
		}
		if strings.HasPrefix(trimmed, "-") {
			value := strings.TrimSpace(strings.TrimPrefix(trimmed, "-"))
			value = strings.Trim(value, `"'`)
			if value != "" {
				patterns = append(patterns, value)
			}
			continue
		}
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			break
		}
	}
	return patterns
}

func lernaWorkspacePatterns(root string) []string {
	data, err := os.ReadFile(filepath.Join(root, "lerna.json"))
	if err != nil {
		return nil
	}
	var file struct {
		Packages []string `json:"packages"`
	}
	if json.Unmarshal(data, &file) != nil {
		return nil
	}
	return file.Packages
}
