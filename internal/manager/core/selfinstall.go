package core

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
)

func (s *Services) InstallSelf() (string, bool, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", false, err
	}
	destination := filepath.Join(config.UserHome(), ".local", "bin", "idev")
	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return "", false, err
	}
	data, err := os.ReadFile(executable)
	if err != nil {
		return "", false, err
	}
	if info, statErr := os.Lstat(destination); statErr == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			target, targetErr := filepath.EvalSymlinks(destination)
			source, sourceErr := filepath.EvalSymlinks(executable)
			if targetErr == nil && sourceErr == nil && target == source {
				return destination, false, nil
			}
			return "", false, fmt.Errorf("refusing to replace symlink %s", destination)
		}
		current, readErr := os.ReadFile(destination)
		if readErr == nil && bytes.Equal(current, data) {
			return destination, false, nil
		}
		return "", false, fmt.Errorf("refusing to replace %s; remove it or use the installer", destination)
	}
	if current, err := os.ReadFile(destination); err == nil && string(current) == string(data) {
		return destination, false, nil
	}
	temporary := destination + ".tmp"
	if err := os.WriteFile(temporary, data, 0755); err != nil {
		return "", false, err
	}
	if err := os.Rename(temporary, destination); err != nil {
		return "", false, err
	}
	return destination, true, nil
}
