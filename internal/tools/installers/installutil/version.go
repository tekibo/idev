package installutil

import (
	"strings"

	"github.com/tekibo/idev/internal/platform/execx"
)

// VersionMatches runs one version command and checks the pin.
func VersionMatches(runner execx.Runner, bin string, args []string, version string) bool {
	output, err := runner.Output(execx.Command{Name: bin, Args: args})
	return err == nil && strings.Contains(string(output), version)
}
