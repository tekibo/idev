package manager

import (
	"fmt"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
)

// signedPythonTool runs a helper that imports cryptography from the private venv.
func (a *App) signedPythonTool(name string, args []string) (UpstreamCommand, error) {
	path := a.Config.Path("apple-tools", "tools", name)
	if !config.IsFile(path) {
		return UpstreamCommand{}, fmt.Errorf("%s not found; run idev setup", path)
	}
	python := filepath.Join(a.Config.PymobileRoot(), "bin", "python")
	if !config.IsExec(python) {
		return UpstreamCommand{}, fmt.Errorf("private pymobiledevice3 venv missing; run idev setup --pymobile")
	}
	return UpstreamCommand{Name: python, Args: append([]string{path}, args...), Env: a.Config.SwiftEnv()}, nil
}

func (a *App) pythonTool(name string, args []string) (UpstreamCommand, error) {
	path := a.Config.Path("apple-tools", "tools", name)
	if !config.IsFile(path) {
		return UpstreamCommand{}, fmt.Errorf("%s not found; run idev setup", path)
	}
	python, ok := a.Config.FindBinary("python3")
	if !ok {
		return UpstreamCommand{}, fmt.Errorf("python3 not found")
	}
	return UpstreamCommand{Name: python, Args: append([]string{path}, args...), Env: a.Config.SwiftEnv()}, nil
}
