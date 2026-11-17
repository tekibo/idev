package expo

import (
	"fmt"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
)

// Command builds the installed Expo wrapper command with the manager environment.
func Command(cfg config.Config, args []string) (execx.Command, error) {
	wrapper := filepath.Join(cfg.Home, "react-native", "expo", "bin", "expo")
	if !config.IsExec(wrapper) {
		return execx.Command{}, fmt.Errorf("expo arm not installed at %s; run idev setup", wrapper)
	}
	if _, err := ResolveCLI(cfg.Project); err != nil {
		return execx.Command{}, err
	}
	if _, err := HookPath(cfg.Home); err != nil {
		return execx.Command{}, err
	}
	return execx.Command{Name: wrapper, Args: args, Dir: cfg.Project, Env: cfg.ReactNativeEnv("")}, nil
}
