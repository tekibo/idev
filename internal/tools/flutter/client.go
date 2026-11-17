// Package flutter wraps the pinned Apple tooling Flutter scripts.
package flutter

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
)

// ErrMissing reports that the Flutter scripts are absent.
var ErrMissing = errors.New("Flutter scripts not found; run idev setup")

// ErrInvalidCommand rejects path separators and traversal in a script name.
var ErrInvalidCommand = errors.New("invalid Flutter command")

// ValidName reports whether one script name stays inside the Flutter tree.
func ValidName(name string) bool {
	return name != "" && name != "." && name != ".." && filepath.Base(name) == name && !strings.ContainsAny(name, `/\\`)
}

// Client runs the pinned Flutter scripts.
type Client struct {
	Config config.Config
	Runner execx.Runner
}

// New builds a Flutter script client.
func New(cfg config.Config, runner execx.Runner) *Client {
	return &Client{Config: cfg, Runner: runner}
}

// Command returns the command for one Flutter script.
func (c *Client) Command(name string, args []string) (execx.Command, error) {
	if !ValidName(name) {
		return execx.Command{}, ErrInvalidCommand
	}
	path := c.Config.Path("flutter", name)
	if !config.IsExec(path) {
		return execx.Command{}, ErrMissing
	}
	return execx.Command{Name: path, Args: args, Env: c.Config.SwiftEnv()}, nil
}

// Preview returns the command preview for one Flutter script.
func (c *Client) Preview(name string, args []string) string {
	if !ValidName(name) {
		return ""
	}
	return execx.Command{Name: c.Config.Path("flutter", name), Args: args}.Preview()
}
