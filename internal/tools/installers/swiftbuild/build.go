package swiftbuild

import (
	"fmt"
	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
	"github.com/tekibo/idev/internal/tools/installers/installutil"
	"io"
	"path/filepath"
)

func buildProduct(cfg config.Config, runner execx.Runner, dir, product string, extra map[string]string, stdout, stderr io.Writer) (string, error) {
	swift := filepath.Join(cfg.SwiftBin(), "swift")
	if !config.IsExec(swift) {
		return "", ErrSwiftMissing
	}
	if err := runner.Run(execx.Command{
		Name:   swift,
		Args:   []string{"build", "-c", "release", "--product", product},
		Dir:    dir,
		Env:    installutil.SwiftEnv(cfg, extra),
		Stdout: stdout,
		Stderr: stderr,
	}); err != nil {
		return "", err
	}
	binary := filepath.Join(dir, ".build", "release", product)
	if !config.IsFile(binary) {
		return "", fmt.Errorf("%w: %s", ErrBuildMissing, binary)
	}
	return binary, nil
}
