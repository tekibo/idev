// Package xtool owns xtool command wrappers.
package xtool

import (
	"io"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
)

// Client runs xtool through one runner.
type Client struct {
	Runner execx.Runner
	Path   string
	Env    []string
}

// NewWithRunner resolves xtool with a caller-owned runner.
func NewWithRunner(cfg config.Config, runner execx.Runner) *Client {
	return &Client{Runner: runner, Path: cfg.XtoolBin(), Env: cfg.BaseEnv()}
}

// Available reports whether xtool was found.
func (c *Client) Available() bool {
	return c.Path != "" && config.IsExec(c.Path)
}

// Run runs xtool with inherited stdio.
func (c *Client) Run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	return c.RunEnv(c.Env, args, stdin, stdout, stderr)
}

// RunEnv runs xtool with an explicit environment.
func (c *Client) RunEnv(env []string, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	return c.Runner.Run(execx.Command{
		Name:   c.Path,
		Args:   args,
		Env:    env,
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	})
}

// Output runs xtool and captures stdout.
func (c *Client) Output(args []string) ([]byte, error) {
	if !c.Available() {
		return nil, ErrMissing
	}
	return c.Runner.Output(execx.Command{Name: c.Path, Args: args, Env: c.Env})
}
