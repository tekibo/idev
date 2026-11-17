package system

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/tekibo/idev/internal/platform/execx"
)

// Probe performs source-agnostic system detection.
type Probe struct {
	Runner    execx.Runner
	once      sync.Once
	libraries []byte
	libraryOK bool
}

// New builds a probe.
func New(runner execx.Runner) *Probe {
	return &Probe{Runner: runner}
}

// Command finds a command on PATH.
func (p *Probe) Command(name string) (string, bool) {
	path, err := exec.LookPath(name)
	return path, err == nil
}

// PkgConfig reports a pkg-config module.
func (p *Probe) PkgConfig(name string) bool {
	err := p.Runner.Run(execx.Command{
		Name:   "pkg-config",
		Args:   []string{"--exists", name},
		Stdout: io.Discard,
		Stderr: io.Discard,
	})
	return err == nil
}

// Library reports a shared library in the loader cache.
func (p *Probe) Library(name string) bool {
	p.once.Do(func() {
		output, err := p.Runner.Output(execx.Command{Name: "ldconfig", Args: []string{"-p"}})
		p.libraries = output
		p.libraryOK = err == nil
	})
	return p.libraryOK && strings.Contains(string(p.libraries), name)
}

// Socket reports a socket path.
func (p *Probe) Socket(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode()&os.ModeSocket != 0
}
