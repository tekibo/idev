package system

import (
	"fmt"
	"strings"

	"github.com/tekibo/idev/internal/platform/execx"
)

type apt struct {
	runner execx.Runner
	distro string
}

func (m *apt) Name() string { return "apt" }

func (m *apt) Distro() string { return m.distro }

func (m *apt) Install(packages []string) (string, []string) {
	return "apt", append([]string{"install", "-y"}, packages...)
}

func (m *apt) Preview(packages []string) (string, error) {
	args := append([]string{"-s", "install"}, packages...)
	output, err := m.runner.Output(execx.Command{Name: "apt-get", Args: args})
	if err != nil {
		return "", fmt.Errorf("apt preview: %w", err)
	}
	for _, line := range strings.Split(string(output), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 || fields[0] != "Inst" {
			continue
		}
		if strings.HasPrefix(fields[2], "[") {
			return fmt.Sprintf("transaction would upgrade %s", fields[1]), nil
		}
	}
	return "", nil
}

func (m *apt) Installed(pkg string) (string, error) {
	output, err := m.runner.Output(execx.Command{Name: "dpkg-query", Args: []string{"-W", "-f=${Version}", pkg}})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
