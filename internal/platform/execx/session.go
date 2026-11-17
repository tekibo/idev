package execx

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

// RunSession runs the command. Human mode inherits stdio; JSON mode emits events.
func RunSession(ctx context.Context, command Command, in io.Reader, stdout, stderr io.Writer, jsonMode bool) error {
	if jsonMode {
		return runEvents(ctx, command, in, stdout)
	}
	return runAttached(ctx, command, in, stdout, stderr)
}

func runAttached(ctx context.Context, command Command, in io.Reader, stdout, stderr io.Writer) error {
	command.Stdin = in
	command.Stdout = stdout
	command.Stderr = stderr
	cmd := NewCommandContext(ctx, command)
	term := attachTerminal(cmd, in)
	if err := cmd.Start(); err != nil {
		_ = term.restore()
		return err
	}
	err := waitGroup(ctx, cmd)
	return errors.Join(err, term.restore())
}

func waitGroup(ctx context.Context, cmd *exec.Cmd) error {
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(stop)
	select {
	case err := <-done:
		killGroup(cmd)
		return err
	case <-stop:
		killGroup(cmd)
		return waitKilled(cmd, done)
	case <-ctx.Done():
		killGroup(cmd)
		return waitKilled(cmd, done)
	}
}

func waitKilled(cmd *exec.Cmd, done <-chan error) error {
	select {
	case err := <-done:
		return err
	case <-time.After(3 * time.Second):
		killGroupHard(cmd)
		return <-done
	}
}

func doneEvent(err error) Event {
	event := Event{Type: "done", State: "succeeded"}
	if err == nil {
		return event
	}
	event.State = "failed"
	event.Error = err.Error()
	if errors.Is(err, context.Canceled) {
		event.State = "canceled"
		return event
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		event.Code = exitErr.ExitCode()
		if event.Code < 0 {
			event.State = "canceled"
			event.Code = 0
		}
	}
	return event
}

func closeFiles(files ...*os.File) {
	for _, file := range files {
		_ = file.Close()
	}
}
