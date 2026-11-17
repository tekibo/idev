package xtool

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
	"github.com/tekibo/idev/internal/tools/installers/installutil"
)

func copyRuntime(cfg config.Config, runner execx.Runner, binDir, dest string) error {
	built := filepath.Join(binDir, "xtool")
	output, err := runner.Output(execx.Command{Name: "ldd", Args: []string{built}, Env: cfg.XtoolBuildEnv()})
	if err != nil {
		return err
	}
	prefix := cfg.SwiftLib() + string(os.PathSeparator)
	for _, line := range strings.Split(string(output), "\n") {
		fields := strings.Fields(line)
		for index, field := range fields {
			if field != "=>" || index+1 >= len(fields) {
				continue
			}
			target := fields[index+1]
			if target == "not found" {
				return fmt.Errorf("unresolved runtime library: %s", line)
			}
			if !strings.HasPrefix(target, prefix) {
				continue
			}
			if err := installutil.CopyFile(target, filepath.Join(dest, filepath.Base(target)), 0644); err != nil {
				return err
			}
		}
	}
	libXADI := filepath.Join(binDir, "libXADI.so")
	if config.IsFile(libXADI) {
		if err := installutil.CopyFile(libXADI, filepath.Join(dest, "libXADI.so"), 0644); err != nil {
			return err
		}
	}
	return nil
}
