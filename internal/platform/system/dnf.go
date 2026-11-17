package system

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/tekibo/idev/internal/platform/execx"
)

type dnf struct {
	runner execx.Runner
	distro string
}

func (m *dnf) Name() string { return "dnf" }

func (m *dnf) Distro() string { return m.distro }

func (m *dnf) Install(packages []string) (string, []string) {
	return "dnf", append([]string{"install", "-y"}, packages...)
}

func (m *dnf) Preview(packages []string) (string, error) {
	args := append([]string{"install", "--assumeno"}, packages...)
	var stderr bytes.Buffer
	output, err := m.runner.Output(execx.Command{Name: "dnf", Args: args, Stderr: &stderr})
	text := string(output) + stderr.String()
	if strings.Contains(text, "Upgrading:") {
		return "transaction would upgrade an installed package", nil
	}
	if err != nil && !strings.Contains(text, "Installing:") {
		return "", fmt.Errorf("dnf preview: %w", err)
	}
	return "", nil
}

func (m *dnf) Installed(pkg string) (string, error) {
	output, err := m.runner.Output(execx.Command{Name: "rpm", Args: []string{"-q", "--qf", "%{VERSION}-%{RELEASE}", pkg}})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
