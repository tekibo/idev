package execx

import (
	"context"
	"os/exec"
	"time"
)

// Exec is the real context-aware runner.
type Exec struct{}

// NewCommand builds an exec command for a PTY host.
func NewCommand(command Command) *exec.Cmd {
	return NewCommandContext(context.Background(), command)
}

// NewCommandContext builds a context-aware exec command for a PTY host.
func NewCommandContext(ctx context.Context, command Command) *exec.Cmd {
	cmd := exec.CommandContext(ctx, command.Name, command.Args...)
	cmd.Dir = command.Dir
	cmd.Env = command.Env
	cmd.Stdin = command.Stdin
	cmd.Stdout = command.Stdout
	cmd.Stderr = command.Stderr
	cmd.WaitDelay = 5 * time.Second
	return cmd
}

// Run runs the command with its stdio.
func (Exec) Run(command Command) error {
	return Exec{}.RunContext(context.Background(), command)
}

// RunContext runs the command under a context.
func (Exec) RunContext(ctx context.Context, command Command) error {
	return NewCommandContext(ctx, command).Run()
}

// Output captures stdout.
func (Exec) Output(command Command) ([]byte, error) {
	return Exec{}.OutputContext(context.Background(), command)
}

// OutputContext captures stdout under a context.
func (Exec) OutputContext(ctx context.Context, command Command) ([]byte, error) {
	return NewCommandContext(ctx, command).Output()
}
