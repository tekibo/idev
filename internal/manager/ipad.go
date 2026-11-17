package manager

import (
	"fmt"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
)

// IpadScreenCommand returns the iPad viewer command.
func (a *App) IpadScreenCommand(args []string) (UpstreamCommand, error) {
	path := a.Config.Path("apple-tools", "tools", "ipad-screen.sh")
	if !config.IsExec(path) {
		return UpstreamCommand{}, fmt.Errorf("%s not found; run idev setup --apple-tools", path)
	}
	overrides := map[string]string{
		"IPAD_SCREEN_HOME":   a.Config.Home,
		"IPAD_SCREEN_CONFIG": a.Config.Path("config", "ipad-screen"),
		"IPAD_SCREEN_STATE":  a.Config.Path("state", "ipad-screen"),
		"IPAD_SCREEN_BIN":    a.Config.Path("bin", "ipad-screen"),
	}
	python := filepath.Join(a.Config.PymobileRoot(), "bin", "python")
	if config.IsExec(python) {
		overrides["IPAD_SCREEN_PYTHON"] = python
	}
	return UpstreamCommand{Name: path, Args: args, Env: a.Config.SwiftEnvWith(overrides)}, nil
}
