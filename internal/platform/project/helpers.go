package project

import (
	"os"
	"path/filepath"
)

func hasFile(dir, name string) bool {
	info, err := os.Stat(filepath.Join(dir, name))
	return err == nil && !info.IsDir()
}

func hasXcodeProject(dir string) string {
	matches, _ := filepath.Glob(filepath.Join(dir, "*.xcodeproj"))
	if len(matches) == 0 {
		return ""
	}
	return matches[0]
}

func pathExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

func fileExists(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func addWarning(info *Info, text string) {
	info.Warnings = append(info.Warnings, text)
}

func absolute(dir, path string) string {
	if path == "" || filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(dir, path)
}
