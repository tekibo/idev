// Package expo owns the project Expo CLI adapter.
package expo

import (
	"fmt"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
)

// ResolveCLI finds the project local Expo CLI.
func ResolveCLI(project string) (string, error) {
	binary := filepath.Join(project, "node_modules", ".bin", "expo")
	if config.IsExec(binary) {
		return binary, nil
	}
	return "", fmt.Errorf("no Expo CLI at %s; install project dependencies first", binary)
}

// HookPath returns the manager-owned Node hook.
func HookPath(home string) (string, error) {
	hook := filepath.Join(home, "react-native", "tools", "node-darwin.js")
	if !config.IsFile(hook) {
		return "", fmt.Errorf("expo arm not installed at %s; run idev setup", hook)
	}
	return hook, nil
}
