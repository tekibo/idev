package xtool

import (
	"github.com/tekibo/idev/internal/platform/execx"
)

// InstallCommand builds the xtool install command.
func (c *Client) InstallCommand(bundle, udid string) (execx.Command, error) {
	if !c.Available() {
		return execx.Command{}, ErrMissing
	}
	args := []string{"install"}
	if udid != "" {
		args = append(args, "--udid", udid)
	} else {
		args = append(args, "--usb")
	}
	args = append(args, bundle)
	return execx.Command{Name: c.Path, Args: args, Env: c.Env}, nil
}
