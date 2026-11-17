package system

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/tekibo/idev/internal/platform/execx"
)

// Manager is one distro package manager adapter.
type Manager interface {
	Name() string
	Distro() string
	Install(packages []string) (string, []string)
	Preview(packages []string) (string, error)
	Installed(pkg string) (string, error)
}

// ErrNoManager reports that no supported package manager was found.
var ErrNoManager = errors.New("no supported package manager; install pacman, apt, dnf")

// Detect returns the first supported package manager.
func Detect(runner execx.Runner) (Manager, error) {
	distro := DistroName()
	if _, err := exec.LookPath("pacman"); err == nil {
		return &pacman{runner: runner, distro: distro}, nil
	}
	if _, err := exec.LookPath("apt"); err == nil {
		return &apt{runner: runner, distro: distro}, nil
	}
	if _, err := exec.LookPath("dnf"); err == nil {
		return &dnf{runner: runner, distro: distro}, nil
	}
	return nil, ErrNoManager
}

// DistroName returns PRETTY_NAME from /etc/os-release.
func DistroName() string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "linux"
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
		}
	}
	return "linux"
}

// ManagerCheck reports a missing package manager.
func ManagerCheck(err error) Check {
	return Check{
		ID:        "system.package_manager",
		Class:     ClassSystem,
		Status:    StatusMissing,
		FixableBy: FixUser,
		Detail:    err.Error(),
		Hint:      "install pacman, apt, dnf",
		Critical:  true,
	}
}
