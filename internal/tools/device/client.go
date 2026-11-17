// Package device owns device discovery, install, and launch.
package device

import (
	"fmt"
	"io"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
	"github.com/tekibo/idev/internal/tools/xtool"
)

// Client owns device discovery, install, and launch.
type Client struct {
	Config config.Config
	Runner execx.Runner
	Xtool  *xtool.Client
}

// NewWithRunner builds a device client with a caller-owned runner.
func NewWithRunner(cfg config.Config, runner execx.Runner) *Client {
	return &Client{Config: cfg, Runner: runner, Xtool: xtool.NewWithRunner(cfg, runner)}
}

// pymobile returns the private pymobiledevice3 binary.
func (c *Client) pymobile() (string, bool) {
	path := c.Config.PymobileBin()
	return path, config.IsExec(path)
}

// List prints connected devices.
func (c *Client) List(stdin io.Reader, stdout, stderr io.Writer) error {
	if c.Xtool.Available() {
		if err := c.Xtool.Devices(stdin, stdout, stderr); err == nil {
			return nil
		}
		if _, err := fmt.Fprintln(stderr, "idev devices: xtool devices failed; trying pymobiledevice3"); err != nil {
			return err
		}
	}
	pmd, ok := c.pymobile()
	if !ok {
		return ErrNoDeviceTool
	}
	return c.Runner.Run(execx.Command{
		Name:   pmd,
		Args:   []string{"usbmux", "list"},
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	})
}
