package execx

import (
	"os/exec"
)

// Runner runs external commands.
type Runner interface {
	Run(Command) error
	Output(Command) ([]byte, error)
}

// ExitCode returns the process exit code, 1 for other errors.
func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	if exitError, ok := err.(*exec.ExitError); ok {
		return exitError.ExitCode()
	}
	return 1
}
