package xtool

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
)

func build(cfg config.Config, runner execx.Runner, source string, stdout, stderr io.Writer) error {
	swift := filepath.Join(cfg.SwiftBin(), "swift")
	if !config.IsExec(swift) {
		return ErrSwiftMissing
	}
	return runner.Run(execx.Command{
		Name:   swift,
		Args:   []string{"build", "-c", "release", "--product", "xtool", "-Xswiftc", "-no-toolchain-stdlib-rpath"},
		Dir:    source,
		Env:    cfg.XtoolBuildEnv(),
		Stdout: stdout,
		Stderr: stderr,
	})
}

func binPath(cfg config.Config, runner execx.Runner, source string) (string, error) {
	output, err := runner.Output(execx.Command{
		Name: filepath.Join(cfg.SwiftBin(), "swift"),
		Args: []string{"build", "-c", "release", "--show-bin-path"},
		Dir:  source,
		Env:  cfg.XtoolBuildEnv(),
	})
	if err != nil {
		return "", err
	}
	path := strings.TrimSpace(string(output))
	if path == "" {
		return "", fmt.Errorf("empty Swift build path")
	}
	return path, nil
}
