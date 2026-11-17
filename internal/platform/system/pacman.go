package system

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/tekibo/idev/internal/platform/execx"
)

type pacman struct {
	runner execx.Runner
	distro string
}

func (m *pacman) Name() string { return "pacman" }

func (m *pacman) Distro() string { return m.distro }

func (m *pacman) Install(packages []string) (string, []string) {
	return "pacman", append([]string{"-S", "--needed", "--noconfirm"}, packages...)
}

func (m *pacman) Preview(packages []string) (string, error) {
	args := []string{"-S", "--needed", "--print", "--print-format", "%n %v"}
	args = append(args, packages...)
	var stderr bytes.Buffer
	output, err := m.runner.Output(execx.Command{Name: "pacman", Args: args, Stderr: &stderr})
	if err != nil {
		detail := strings.TrimSpace(stderr.String())
		if detail == "" {
			detail = err.Error()
		}
		return "", fmt.Errorf("pacman preview: %s", detail)
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return "", fmt.Errorf("pacman preview returned no packages")
	}
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		installed, err := m.Installed(fields[0])
		if err == nil && installed != fields[1] {
			return fmt.Sprintf("transaction would upgrade %s from %s to %s", fields[0], installed, fields[1]), nil
		}
	}
	return "", nil
}

func (m *pacman) Installed(pkg string) (string, error) {
	output, err := m.runner.Output(execx.Command{Name: "pacman", Args: []string{"-Q", pkg}})
	if err != nil {
		return "", err
	}
	fields := strings.Fields(string(output))
	if len(fields) < 2 {
		return "", fmt.Errorf("pacman -Q %s: unexpected output", pkg)
	}
	return fields[1], nil
}
