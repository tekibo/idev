package xtool

import (
	"fmt"
	"strings"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
)

func verify(cfg config.Config, runner execx.Runner) error {
	output, err := runner.Output(execx.Command{
		Name: cfg.XtoolLink(),
		Args: []string{"--version"},
		Env:  cfg.SwiftEnv(),
	})
	if err != nil {
		return fmt.Errorf("xtool --version: %w", err)
	}
	if !strings.Contains(string(output), cfg.Lock.AppleDev.XtoolSHA) {
		return fmt.Errorf("%w: %s", ErrPinMismatch, strings.TrimSpace(string(output)))
	}
	return nil
}
