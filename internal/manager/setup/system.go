package setup

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/tekibo/idev/internal/platform/execx"
	"github.com/tekibo/idev/internal/platform/security/changeset"
	"github.com/tekibo/idev/internal/platform/security/privileged"
	"github.com/tekibo/idev/internal/platform/system"
	"github.com/tekibo/idev/internal/platform/terminal"
)

// SystemPlan is the system transaction view.
type SystemPlan struct {
	Manager  string
	Packages []string
	Command  string
	Name     string
	Args     []string
	Blocked  string
}

// SystemPlan builds the one missing-package command.
func (s *Service) SystemPlan(all bool) SystemPlan {
	probe := system.New(s.Runner)
	manager, err := system.Detect(s.Runner)
	if err != nil {
		return SystemPlan{Blocked: err.Error()}
	}
	checks := system.Checks(probe, manager, all)
	plan := system.BuildPlan(manager, checks)
	return SystemPlan{
		Manager:  plan.Manager,
		Packages: plan.Packages,
		Command:  plan.Command,
		Name:     plan.Name,
		Args:     plan.Args,
		Blocked:  plan.Blocked,
	}
}

// PasswordProvider supplies one sudo password.
type PasswordProvider interface {
	Password(context.Context) ([]byte, error)
}

// TerminalPasswordProvider reads the password from a terminal.
func (s *Service) TerminalPasswordProvider(in *os.File, out io.Writer) PasswordProvider {
	return terminalPassword{in: in, out: out}
}

// ApplySystemPlan runs the exact package transaction.
func (s *Service) ApplySystemPlan(ctx context.Context, plan SystemPlan, provider PasswordProvider, stdout, stderr io.Writer) error {
	if plan.Blocked != "" {
		return fmt.Errorf("blocked: %s", plan.Blocked)
	}
	if len(plan.Packages) == 0 {
		return nil
	}
	if plan.Name == "" {
		return fmt.Errorf("no package manager command")
	}
	command := execx.Command{
		Name:   plan.Name,
		Args:   plan.Args,
		Stdout: stdout,
		Stderr: stderr,
	}
	runner := privileged.New(s.Runner, provider, privileged.ChangeRecorder{Store: changeset.New(s.Config.Home)})
	return runner.Run(ctx, privileged.Plan{Command: command, Reason: "system packages", Packages: plan.Packages})
}

type terminalPassword struct {
	in  *os.File
	out io.Writer
}

func (p terminalPassword) Password(context.Context) ([]byte, error) {
	return terminal.ReadPassword("sudo password: ", p.in, p.out)
}
