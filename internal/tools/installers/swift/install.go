package swift

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/download"
	"github.com/tekibo/idev/internal/platform/execx"
)

// Install downloads, verifies, extracts, and links the official Swift toolchain.
func Install(ctx context.Context, cfg config.Config, runner execx.Runner, force bool, stdout io.Writer) (Result, error) {
	spec, err := ForArch(cfg, runtime.GOARCH)
	if err != nil {
		return Result{}, err
	}
	if !force && config.IsExec(filepath.Join(cfg.SwiftBin(), "swift")) {
		if err := InstallLinkers(cfg); err != nil {
			return Result{}, err
		}
		if err := Verify(ctx, cfg, runner); err == nil {
			return Result{Version: spec.Version, Root: spec.Root, Cache: spec.Cache, Already: true}, nil
		}
		return Result{}, fmt.Errorf("existing Swift toolchain is incomplete; run with --force")
	}
	if err := download.File(ctx, spec.URL, spec.SHA256, spec.Cache, download.Log(stdout, "swift")); err != nil {
		return Result{}, err
	}
	parent := cfg.Path("toolchains")
	if err := os.MkdirAll(parent, 0755); err != nil {
		return Result{}, err
	}
	temporary, err := os.MkdirTemp(parent, ".swift-")
	if err != nil {
		return Result{}, err
	}
	defer func() { _ = os.RemoveAll(temporary) }()
	if err := Extract(spec.Cache, temporary); err != nil {
		return Result{}, err
	}
	if !config.IsExec(filepath.Join(temporary, "usr", "bin", "swift")) {
		return Result{}, fmt.Errorf("archive did not contain usr/bin/swift")
	}
	if force {
		if err := os.RemoveAll(spec.Root); err != nil {
			return Result{}, err
		}
	} else if _, err := os.Stat(spec.Root); err == nil {
		return Result{}, fmt.Errorf("Swift root already exists: %s", spec.Root)
	}
	if err := os.Rename(temporary, spec.Root); err != nil {
		return Result{}, err
	}
	keep := false
	defer func() {
		if !keep {
			rollback(cfg, spec.Root)
		}
	}()
	if err := InstallCompat(cfg); err != nil {
		return Result{}, err
	}
	if err := InstallLinks(cfg); err != nil {
		return Result{}, err
	}
	if err := InstallLinkers(cfg); err != nil {
		return Result{}, err
	}
	if err := Verify(ctx, cfg, runner); err != nil {
		return Result{}, err
	}
	keep = true
	return Result{Version: spec.Version, Root: spec.Root, Cache: spec.Cache}, nil
}
