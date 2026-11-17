package official

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
	"github.com/tekibo/idev/internal/platform/release"
	"github.com/tekibo/idev/internal/tools/installers/installutil"
)

// Result reports one official binary install.
type Result struct {
	Name    string
	Version string
	Bin     string
	Already bool
}

// Install downloads and verifies one official binary.
func Install(ctx context.Context, cfg config.Config, runner execx.Runner, spec Spec, force bool) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	if spec.Component.URL == "" || spec.Component.SHA256 == "" {
		return Result{}, ErrMissing
	}
	if !force && config.IsExec(spec.Bin) && installutil.VersionMatches(runner, spec.Bin, spec.VersionArgs, spec.Version) {
		return Result{Name: spec.Name, Version: spec.Version, Bin: spec.Bin, Already: true}, nil
	}
	if force {
		_ = os.Remove(spec.Bin)
	}
	if err := os.MkdirAll(filepath.Dir(spec.Bin), 0755); err != nil {
		return Result{}, err
	}
	if err := release.FetchBinary(ctx, spec.Source, spec.Component, spec.Name, spec.Bin, cfg.Path("cache", "downloads")); err != nil {
		return Result{}, err
	}
	if !installutil.VersionMatches(runner, spec.Bin, spec.VersionArgs, spec.Version) {
		return Result{}, fmt.Errorf("%w: %s", ErrVersionMismatch, spec.Version)
	}
	return Result{Name: spec.Name, Version: spec.Version, Bin: spec.Bin}, nil
}
