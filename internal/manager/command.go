package manager

import (
	"context"
	"io"

	"github.com/tekibo/idev/internal/platform/execx"
)

// RunCommand runs one external command with the shared terminal contract.
func (a *App) RunCommand(ctx context.Context, name string, args, env []string, dir string,
	stdin io.Reader, stdout, stderr io.Writer, jsonMode bool) error {
	return execx.RunSession(ctx, execx.Command{
		Name: name, Args: args, Env: env, Dir: dir,
	}, stdin, stdout, stderr, jsonMode)
}

// RunApple runs one apple-dev command through the terminal contract.
func (a *App) RunApple(ctx context.Context, name string, args []string,
	stdin io.Reader, stdout, stderr io.Writer, jsonMode bool) error {
	command, err := a.ProjectService.AppleCommand(name, args)
	if err != nil {
		return err
	}
	return a.RunCommand(ctx, command.Name, command.Args, command.Env, command.Dir,
		stdin, stdout, stderr, jsonMode)
}

// LogSink is one line-oriented event writer.
type LogSink interface {
	io.Writer
	Flush()
	Done(code int)
}

// LogWriter returns a line-oriented event writer.
func (a *App) LogWriter(out io.Writer, command, stream string) LogSink {
	return execx.LogWriter(out, command, stream)
}
