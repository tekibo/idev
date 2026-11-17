package swiftbuild

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
	"github.com/tekibo/idev/internal/platform/release"
	"github.com/tekibo/idev/internal/tools/installers/installutil"
)

// Result reports one Swift Build and Expo engine install.
type Result struct {
	Version   string
	Source    string
	Bin       string
	EngineBin string
	Already   bool
}

// Install builds Swift Build and the Expo engine from the release source.
func Install(ctx context.Context, cfg config.Config, runner execx.Runner, spec Spec, force bool, stdout, stderr io.Writer) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	if !force && config.IsExec(spec.Bin) && config.IsExec(spec.EngineBin) {
		return Result{Version: spec.Version, Source: spec.SourceRoot, Bin: spec.Bin, EngineBin: spec.EngineBin, Already: true}, nil
	}
	if force {
		_ = os.RemoveAll(spec.SourceRoot)
		_ = os.Remove(spec.Bin)
		_ = os.Remove(spec.EngineBin)
		removeBundles(spec.Bundles)
	}
	if !config.IsExec(filepath.Join(cfg.SwiftBin(), "swift")) {
		return Result{}, ErrSwiftMissing
	}
	if !isDir(spec.EngineSource) {
		return Result{}, ErrEngineMissing
	}
	if spec.Component.URL == "" || spec.Component.SHA256 == "" {
		return Result{}, ErrSourceMissing
	}
	if !config.IsFile(filepath.Join(spec.SourceRoot, "Package.swift")) {
		if _, err := release.Fetch(ctx, spec.Source, spec.Component, spec.SourceRoot, cfg.Path("cache", "downloads")); err != nil {
			return Result{}, err
		}
	}
	service, err := buildProduct(cfg, runner, spec.SourceRoot, "SWBBuildServiceBundle", nil, stdout, stderr)
	if err != nil {
		return Result{}, err
	}
	engine, err := buildProduct(cfg, runner, spec.EngineSource, "expo-build", map[string]string{"SWIFT_BUILD_PATH": spec.SourceRoot}, stdout, stderr)
	if err != nil {
		return Result{}, err
	}
	if err := os.MkdirAll(spec.Bundles, 0755); err != nil {
		return Result{}, err
	}
	if err := installutil.CopyFile(service, spec.Bin, 0755); err != nil {
		return Result{}, err
	}
	if err := installutil.CopyFile(engine, spec.EngineBin, 0755); err != nil {
		return Result{}, err
	}
	_ = os.RemoveAll(filepath.Join(spec.EngineSource, ".build"))
	if err := copyBundles(spec.SourceRoot, spec.Bundles); err != nil {
		return Result{}, err
	}
	_ = os.RemoveAll(spec.SourceRoot)
	return Result{Version: spec.Version, Source: spec.SourceRoot, Bin: spec.Bin, EngineBin: spec.EngineBin}, nil
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
