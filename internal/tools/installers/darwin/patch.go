package darwin

import (
	"io"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
	"github.com/tekibo/idev/internal/tools/installers/installutil"
)

func patchSDK(cfg config.Config, runner execx.Runner, spec Spec, stdout, stderr io.Writer) error {
	python, ok := cfg.FindBinary("python3")
	if !ok {
		return ErrPythonMissing
	}
	return runner.Run(execx.Command{
		Name:   python,
		Args:   []string{"-c", patchScript, spec.SDKRoot, spec.SwiftRoot},
		Env:    installutil.SwiftEnv(cfg, nil),
		Stdout: stdout,
		Stderr: stderr,
	})
}
