package cocoapods

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
)

// Result reports one CocoaPods install.
type Result struct {
	Version string
	Root    string
	Bin     string
	Link    string
	Already bool
}

// Install creates the private gem tree and the pod wrapper. The gem tree also
// owns erb because some distributions ship it outside the system Ruby.
func Install(ctx context.Context, cfg config.Config, runner execx.Runner, force bool, stdout, stderr io.Writer) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	spec := NewSpec(cfg)
	if !force && config.IsExec(spec.Link) {
		if version, err := podVersion(runner, spec); err == nil && strings.Contains(version, spec.Version) {
			if err := writeWrapper(spec.Link); err != nil {
				return Result{}, err
			}
			return Result{Version: spec.Version, Root: spec.Root, Bin: spec.Bin, Link: spec.Link, Already: true}, nil
		}
	}
	if force {
		if err := os.RemoveAll(spec.Root); err != nil {
			return Result{}, err
		}
		_ = os.Remove(spec.Link)
	}
	gem, ok := cfg.FindBinary("gem")
	if !ok {
		return Result{}, ErrGemMissing
	}
	if err := os.MkdirAll(spec.Cache, 0755); err != nil {
		return Result{}, err
	}
	if err := os.MkdirAll(spec.Root, 0755); err != nil {
		return Result{}, err
	}
	for _, item := range []struct {
		name    string
		version string
	}{
		{name: "cocoapods", version: spec.Version},
		{name: "erb", version: spec.Erb},
	} {
		if err := runner.Run(execx.Command{
			Name: gem,
			Args: []string{
				"install", item.name, "-v", item.version,
				"--install-dir", spec.Root,
				"--bindir", filepath.Join(spec.Root, "bin"),
				"--no-document",
			},
			Env:    cfg.GemEnv(),
			Stdout: stdout,
			Stderr: stderr,
		}); err != nil {
			return Result{}, fmt.Errorf("install %s %s: %w", item.name, item.version, err)
		}
	}
	if err := writeWrapper(spec.Link); err != nil {
		return Result{}, err
	}
	version, err := podVersion(runner, spec)
	if err != nil || !strings.Contains(version, spec.Version) {
		return Result{}, fmt.Errorf("%w: expected %s, got %q: %v", ErrVersionMismatch, spec.Version, version, err)
	}
	return Result{Version: spec.Version, Root: spec.Root, Bin: spec.Bin, Link: spec.Link}, nil
}

func podVersion(runner execx.Runner, spec Spec) (string, error) {
	var output bytes.Buffer
	err := runner.Run(execx.Command{
		Name:   spec.Link,
		Args:   []string{"--version"},
		Stdout: &output,
		Stderr: &output,
	})
	return strings.TrimSpace(output.String()), err
}
