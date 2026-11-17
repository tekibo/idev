package update

import (
	"context"
	"fmt"
	"io"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
)

func (s *Service) refreshAgents(ctx context.Context, target string, stdout, stderr io.Writer) {
	if !config.IsFile(s.Config.Path("skills", "idev", ".idev-managed")) {
		return
	}
	err := s.Runner.RunContext(ctx, execx.Command{
		Name:   target,
		Args:   []string{"agents", "install"},
		Stdout: stdout,
		Stderr: stderr,
	})
	if err != nil {
		fmt.Fprintf(stderr, "idev update: refresh agent skill: %v\n", err)
		return
	}
	fmt.Fprintln(stdout, "agents: skill refreshed")
}
