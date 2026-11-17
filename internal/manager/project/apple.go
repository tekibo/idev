package project

import (
	"fmt"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
)

// AppleCommand builds one apple-dev command.
func (s *Service) AppleCommand(name string, args []string) (execx.Command, error) {
	script := s.Config.Path("apple-tools", name)
	if !config.IsExec(script) {
		return execx.Command{}, fmt.Errorf("%s not found; run idev setup --apple-tools", script)
	}
	return execx.Command{Name: script, Args: args, Env: s.Config.SwiftBuildEnv()}, nil
}
