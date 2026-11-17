package xtool

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
	"github.com/tekibo/idev/internal/platform/release"
	"github.com/tekibo/idev/internal/tools/installers/installutil"
)

// Result reports one xtool install.
type Result struct {
	Pin        string
	Root       string
	Link       string
	Source     string
	Already    bool
	Shim       string
	ShimOnPath bool
}

// Install fetches, builds, and links the pinned xtool.
func Install(ctx context.Context, cfg config.Config, runner execx.Runner, force bool, stdout, stderr io.Writer) (Result, error) {
	spec := NewSpec(cfg)
	if cfg.SourcesErr != nil {
		return Result{}, cfg.SourcesErr
	}
	if !config.IsExec(filepath.Join(cfg.SwiftBin(), "swift")) {
		return Result{}, ErrSwiftMissing
	}
	manifest, err := release.Load(cfg.Sources.Release.Manifest)
	if err != nil {
		return Result{}, err
	}
	component, ok := manifest.Components["xtool"]
	if !ok {
		return Result{}, fmt.Errorf("release: xtool component missing")
	}
	if component.Version != spec.Pin {
		return Result{}, fmt.Errorf("release: xtool %s does not match lock %s", component.Version, spec.Pin)
	}
	if !force && config.IsExec(spec.Bin) {
		if err := verify(cfg, runner); err == nil {
			shim, onPath, err := registerShim(cfg)
			if err != nil {
				return Result{}, err
			}
			_ = os.RemoveAll(spec.Source)
			return Result{Pin: spec.Pin, Root: spec.Root, Link: spec.Link, Source: spec.Source, Already: true, Shim: shim, ShimOnPath: onPath}, nil
		}
		return Result{}, fmt.Errorf("existing xtool is incomplete; run with --force")
	}
	if force {
		_ = os.RemoveAll(spec.Root)
		_ = os.Remove(spec.Link)
		_ = os.RemoveAll(spec.Source)
	}
	if err := os.MkdirAll(spec.Root, 0755); err != nil {
		return Result{}, err
	}
	keep := false
	defer func() {
		if !keep {
			rollback(cfg)
		}
	}()
	if _, err := release.Fetch(ctx, cfg.Sources.Release.Manifest, component, spec.Source, cfg.Path("cache", "downloads")); err != nil {
		return Result{}, err
	}
	if err := build(cfg, runner, spec.Source, stdout, stderr); err != nil {
		return Result{}, err
	}
	binDir, err := binPath(cfg, runner, spec.Source)
	if err != nil {
		return Result{}, err
	}
	if err := installutil.CopyFile(filepath.Join(binDir, "xtool"), spec.Bin, 0644); err != nil {
		return Result{}, err
	}
	if err := os.Chmod(spec.Bin, 0755); err != nil {
		return Result{}, err
	}
	if err := copyRuntime(cfg, runner, binDir, spec.Root); err != nil {
		return Result{}, err
	}
	if err := os.MkdirAll(cfg.Path("bin"), 0755); err != nil {
		return Result{}, err
	}
	_ = os.Remove(spec.Link)
	if err := os.Symlink(spec.Bin, spec.Link); err != nil {
		return Result{}, err
	}
	if err := verify(cfg, runner); err != nil {
		return Result{}, err
	}
	shim, onPath, err := registerShim(cfg)
	if err != nil {
		return Result{}, err
	}
	keep = true
	_ = os.RemoveAll(spec.Source)
	return Result{Pin: spec.Pin, Root: spec.Root, Link: spec.Link, Source: spec.Source, Shim: shim, ShimOnPath: onPath}, nil
}
