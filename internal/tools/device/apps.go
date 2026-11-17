package device

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/tekibo/idev/internal/platform/execx"
)

// ResolveBundle returns the installed identifier for one project bundle ID.
func (c *Client) ResolveBundle(bundleID, udid string) string {
	pmd, ok := c.pymobile()
	if !ok {
		return bundleID
	}
	args := []string{"apps", "list"}
	if udid != "" {
		args = append(args, "--udid", udid)
	}
	output, err := c.Runner.Output(execx.Command{Name: pmd, Args: args})
	if err != nil {
		return bundleID
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(output, &raw); err != nil {
		return bundleID
	}
	if apps, ok := raw["Apps"]; ok {
		var nested map[string]json.RawMessage
		if err := json.Unmarshal(apps, &nested); err == nil {
			raw = nested
		}
	}
	if _, ok := raw[bundleID]; ok {
		return bundleID
	}
	suffix := "." + bundleID
	for candidate := range raw {
		if strings.HasSuffix(candidate, suffix) {
			return candidate
		}
	}
	return bundleID
}

// LaunchApp launches an installed app, resolving the signed installed ID.
func (c *Client) LaunchApp(bundleID, udid string, stdout, stderr io.Writer) error {
	return c.pymobileCommand(udid, stdout, stderr,
		"developer", "dvt", "launch", c.ResolveBundle(bundleID, udid))
}

// UninstallApp uninstalls an installed app, resolving the signed installed ID.
func (c *Client) UninstallApp(bundleID, udid string, stdout, stderr io.Writer) error {
	return c.pymobileCommand(udid, stdout, stderr,
		"apps", "uninstall", c.ResolveBundle(bundleID, udid))
}

func (c *Client) pymobileCommand(udid string, stdout, stderr io.Writer, args ...string) error {
	pmd, ok := c.pymobile()
	if !ok {
		return ErrNoDeviceTool
	}
	env := os.Environ()
	if udid != "" {
		env = append(env, "PYMOBILEDEVICE3_UDID="+udid)
	}
	return c.Runner.Run(execx.Command{Name: pmd, Args: args, Env: env, Stdout: stdout, Stderr: stderr})
}
