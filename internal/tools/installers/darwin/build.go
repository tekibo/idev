package darwin

import (
	"fmt"
	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
	"github.com/tekibo/idev/internal/tools/installers/installutil"
	"io"
	"os"
	"path/filepath"
)

func buildActool(cfg config.Config, runner execx.Runner, payload string, force bool, stdout, stderr io.Writer) (string, error) {
	return buildProduct(cfg, runner, filepath.Join(payload, "darwin-tools"), "actool", nil, force, stdout, stderr)
}

func buildOAM(cfg config.Config, runner execx.Runner, payload string, force bool, stdout, stderr io.Writer) (string, error) {
	return buildProduct(cfg, runner, filepath.Join(payload, "oam"), "OpenAppleMacrosServer", []string{"--build-system", "native", "--static-swift-stdlib"}, force, stdout, stderr)
}

func buildProduct(cfg config.Config, runner execx.Runner, dir, product string, extra []string, force bool, stdout, stderr io.Writer) (string, error) {
	swift := filepath.Join(cfg.SwiftBin(), "swift")
	if !config.IsExec(swift) {
		return "", ErrSwiftMissing
	}
	if force {
		if err := os.RemoveAll(filepath.Join(dir, ".build")); err != nil {
			return "", err
		}
	}
	args := []string{"build", "-c", "release"}
	args = append(args, extra...)
	args = append(args, "--product", product)
	if err := runner.Run(execx.Command{
		Name:   swift,
		Args:   args,
		Dir:    dir,
		Env:    installutil.SwiftEnv(cfg, nil),
		Stdout: stdout,
		Stderr: stderr,
	}); err != nil {
		return "", err
	}
	binary := filepath.Join(dir, ".build", "release", product)
	if !config.IsFile(binary) {
		return "", fmt.Errorf("%w: %s", ErrToolMissing, binary)
	}
	return binary, nil
}
