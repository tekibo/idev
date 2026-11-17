package device

import (
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/platform/execx"
)

// InstallPlan is the ordered terminate, install, and launch commands.
type InstallPlan struct {
	Terminate execx.Command
	Install   execx.Command
	Launch    func() (execx.Command, error)
}

// InstallPlan builds the install sequence.
func (c *Client) InstallPlan(bundle, bundleID, udid, devURL string) (InstallPlan, error) {
	if _, err := os.Stat(bundle); err != nil {
		return InstallPlan{}, fmt.Errorf("%s: no such file", bundle)
	}
	if !c.Xtool.Available() {
		return InstallPlan{}, fmt.Errorf("xtool not found")
	}
	plan := InstallPlan{}
	if pmd, ok := c.pymobile(); ok {
		plan.Terminate = execx.Command{
			Name: pmd,
			Args: []string{"developer", "dvt", "pkill", "--bundle", c.ResolveBundle(bundleID, udid)},
		}
	}
	install, err := c.Xtool.InstallCommand(bundle, udid)
	if err != nil {
		return InstallPlan{}, err
	}
	plan.Install = install
	plan.Launch = func() (execx.Command, error) {
		pmd, ok := c.pymobile()
		if !ok {
			return execx.Command{}, ErrNoLaunchTool
		}
		args := []string{"developer", "core-device", "launch-application", c.ResolveBundle(bundleID, udid), "--kill-existing"}
		if devURL != "" {
			args = append(args, "--", "--initialUrl", devURL)
		}
		return execx.Command{Name: pmd, Args: args}, nil
	}
	return plan, nil
}
