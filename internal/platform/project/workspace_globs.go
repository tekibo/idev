package project

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func workspaceCandidates(root string) []string {
	dirs := []string{root}
	patterns := workspacePatterns(root)
	if len(patterns) == 0 {
		if _, ok := inspectNode(root, root); ok {
			return dirs
		}
		patterns = []string{"apps/*", "packages/*"}
	}
	for _, pattern := range patterns {
		if strings.Contains(pattern, "**") {
			continue
		}
		matches, _ := filepath.Glob(filepath.Join(root, pattern))
		for _, match := range matches {
			if info, err := os.Stat(match); err == nil && info.IsDir() {
				dirs = append(dirs, match)
			}
		}
	}
	return uniqueDirs(dirs)
}

func workspacePatterns(root string) []string {
	patterns := []string{}
	patterns = append(patterns, packageWorkspacePatterns(root)...)
	patterns = append(patterns, pnpmWorkspacePatterns(root)...)
	patterns = append(patterns, lernaWorkspacePatterns(root)...)
	return uniqueStrings(patterns)
}

func packageWorkspacePatterns(root string) []string {
	data, err := os.ReadFile(filepath.Join(root, "package.json"))
	if err != nil {
		return nil
	}
	var raw struct {
		Workspaces json.RawMessage `json:"workspaces"`
	}
	if json.Unmarshal(data, &raw) != nil || len(raw.Workspaces) == 0 {
		return nil
	}
	var list []string
	if json.Unmarshal(raw.Workspaces, &list) == nil {
		return list
	}
	var object struct {
		Packages []string `json:"packages"`
	}
	if json.Unmarshal(raw.Workspaces, &object) == nil {
		return object.Packages
	}
	return nil
}

func uniqueDirs(dirs []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		key := filepath.Clean(dir)
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, dir)
	}
	return out
}

func uniqueStrings(values []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		if seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	return out
}
