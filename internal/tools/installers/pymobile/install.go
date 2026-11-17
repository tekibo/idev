package pymobile

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
	"github.com/tekibo/idev/internal/tools/installers/installutil"
)

// Result reports one pymobiledevice3 install.
type Result struct {
	Version string
	Root    string
	Bin     string
	Link    string
	Already bool
}

// Install creates the private venv and pins pymobiledevice3.
func Install(ctx context.Context, cfg config.Config, runner execx.Runner, force bool, stdout, stderr io.Writer) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	spec := NewSpec(cfg)
	if !force && config.IsExec(spec.Link) && installutil.VersionMatches(runner, spec.Link, []string{"version"}, spec.Version) {
		if err := ensureCompat(spec); err != nil {
			return Result{}, err
		}
		return Result{Version: spec.Version, Root: spec.Root, Bin: spec.Bin, Link: spec.Link, Already: true}, nil
	}
	if force {
		if err := os.RemoveAll(spec.Root); err != nil {
			return Result{}, err
		}
		_ = os.Remove(spec.Link)
	}
	python, ok := cfg.FindBinary("python3")
	if !ok {
		return Result{}, ErrPythonMissing
	}
	if err := os.MkdirAll(cfg.Path("bin"), 0755); err != nil {
		return Result{}, err
	}
	if err := os.MkdirAll(spec.Cache, 0755); err != nil {
		return Result{}, err
	}
	env := config.EnvWith(map[string]string{"PIP_CACHE_DIR": spec.Cache})
	if err := runner.Run(execx.Command{
		Name:   python,
		Args:   []string{"-m", "venv", spec.Root},
		Env:    env,
		Stdout: stdout,
		Stderr: stderr,
	}); err != nil {
		return Result{}, err
	}
	pip := filepath.Join(spec.Root, "bin", "pip")
	if !config.IsExec(pip) {
		return Result{}, ErrPipMissing
	}
	if err := runner.Run(execx.Command{
		Name:   pip,
		Args:   []string{"install", "pymobiledevice3==" + spec.Version},
		Env:    env,
		Stdout: stdout,
		Stderr: stderr,
	}); err != nil {
		return Result{}, err
	}
	_ = os.Remove(spec.Link)
	if err := os.Symlink(spec.Bin, spec.Link); err != nil {
		return Result{}, err
	}
	if err := ensureCompat(spec); err != nil {
		return Result{}, err
	}
	if !installutil.VersionMatches(runner, spec.Link, []string{"version"}, spec.Version) {
		return Result{}, fmt.Errorf("%w: %s", ErrVersionMismatch, spec.Version)
	}
	return Result{Version: spec.Version, Root: spec.Root, Bin: spec.Bin, Link: spec.Link}, nil
}
