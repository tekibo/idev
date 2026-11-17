package privileged

import (
	"context"
	"io"
	"time"

	"github.com/tekibo/idev/internal/platform/execx"
)

// PasswordProvider supplies one sudo password.
type PasswordProvider interface {
	Password(context.Context) ([]byte, error)
}

// Runner runs exact privileged plans.
type Runner struct {
	Exec          execx.Runner
	Sudo          string
	Provider      PasswordProvider
	Recorder      Recorder
	PromptTimeout time.Duration
}

// New builds a privileged runner.
func New(execRunner execx.Runner, provider PasswordProvider, recorder Recorder) *Runner {
	return &Runner{Exec: execRunner, Sudo: "sudo", Provider: provider, Recorder: recorder}
}

// Run authenticates when needed and runs the exact command.
func (r *Runner) Run(ctx context.Context, plan Plan) error {
	if plan.Command.Name == "" {
		return ErrEmptyCommand
	}
	if r.canRun() {
		if err := r.runDirect(plan); err != nil {
			return err
		}
	} else {
		if r.Provider == nil {
			return ErrPasswordRequired
		}
		if err := r.authenticate(ctx, plan); err != nil {
			return err
		}
	}
	if r.Recorder != nil {
		return r.Recorder.Record(ctx, plan)
	}
	return nil
}

func (r *Runner) canRun() bool {
	return r.Exec.Run(execx.Command{
		Name:   r.sudo(),
		Args:   []string{"-n", "true"},
		Stdout: io.Discard,
		Stderr: io.Discard,
	}) == nil
}

func (r *Runner) runDirect(plan Plan) error {
	command := plan.Command
	command.Name = r.sudo()
	command.Args = append([]string{"-n", plan.Command.Name}, plan.Command.Args...)
	return r.Exec.Run(command)
}

func (r *Runner) sudo() string {
	if r.Sudo == "" {
		return "sudo"
	}
	return r.Sudo
}
