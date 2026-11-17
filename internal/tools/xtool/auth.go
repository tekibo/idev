package xtool

import (
	"io"

	"github.com/tekibo/idev/internal/platform/redact"
)

// Auth runs xtool auth. Secret values must be typed at the prompt.
func (c *Client) Auth(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if redact.HasSecretFlag(args) {
		return ErrSecretArg
	}
	if len(args) > 0 && args[0] == "status" {
		return c.AuthStatus(stdout)
	}
	return c.run(append([]string{"auth"}, args...), stdin, stdout, stderr)
}

// AuthStatus prints the auth status with the token redacted.
func (c *Client) AuthStatus(stdout io.Writer) error {
	output, err := c.Output([]string{"auth", "status"})
	if err != nil {
		return err
	}
	_, err = io.WriteString(stdout, redact.RedactText(string(output)))
	return err
}

// Devices runs xtool devices.
func (c *Client) Devices(stdin io.Reader, stdout, stderr io.Writer) error {
	return c.run([]string{"devices", "--no-wait"}, stdin, stdout, stderr)
}

func (c *Client) run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if !c.Available() {
		return ErrMissing
	}
	return c.Run(args, stdin, stdout, stderr)
}
