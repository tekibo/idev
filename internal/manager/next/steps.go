package next

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/tekibo/idev/internal/manager/doctor"
	"github.com/tekibo/idev/internal/platform/config"
)

func sourceStep(cfg config.Config) Step {
	warnings := config.SourceWarnings(cfg.Sources)
	if len(warnings) == 0 {
		return Step{ID: "sources", Label: "Source pins", Done: true}
	}
	return Step{
		ID:      "sources",
		Label:   "Source pins",
		Detail:  warnings[0].Field + " is " + warnings[0].Current + ", required " + warnings[0].Required,
		Command: "idev config migrate",
	}
}

func systemStep(report doctor.Report) Step {
	if len(report.SystemPlan.Packages) == 0 {
		return Step{ID: "system", Label: "System packages", Done: true}
	}
	return Step{
		ID:      "system",
		Label:   "System packages",
		Detail:  report.SystemPlan.Command,
		Command: "idev system apply",
	}
}

func baseStep(cfg config.Config, checks map[string]doctor.Check) Step {
	ids := []string{
		"private.swift", "private.xtool", "private.pymobiledevice3",
		"private.rcodesign", "private.ipsw",
	}
	missing := []string{}
	for _, id := range ids {
		if check, ok := checks[id]; !ok || !check.OK() {
			missing = append(missing, strings.TrimPrefix(id, "private."))
		}
	}
	if !config.IsFile(cfg.Path("apple-tools", "ship.sh")) {
		missing = append(missing, "apple-tools")
	}
	if len(missing) == 0 {
		return Step{ID: "base", Label: "Private toolchain", Done: true}
	}
	return Step{
		ID:      "base",
		Label:   "Private toolchain",
		Detail:  "missing: " + strings.Join(missing, ", "),
		Command: "idev setup --base",
	}
}

func agentsStep(cfg config.Config) Step {
	link := filepath.Join(config.UserHome(), ".agents", "skills", "idev")
	if _, err := os.Lstat(link); err == nil {
		return Step{ID: "agents", Label: "Agent skill", Done: true}
	}
	return Step{
		ID:      "agents",
		Label:   "Agent skill",
		Detail:  "not linked into the agent skill directory",
		Command: "idev agents install",
	}
}

func sdkStep(checks map[string]doctor.Check) Step {
	check, ok := checks["private.sdk"]
	if ok && check.OK() {
		return Step{ID: "sdk", Label: "Darwin SDK", Done: true}
	}
	return Step{
		ID:      "sdk",
		Label:   "Darwin SDK",
		Detail:  "the Xcode.xip is user-supplied",
		Command: "idev sdk install <Xcode.xip>",
	}
}

func appleToolsStep(checks map[string]doctor.Check) Step {
	check, ok := checks["private.apple_tools"]
	if ok && check.OK() {
		return Step{ID: "apple-tools", Label: "Darwin tools and OAM", Done: true}
	}
	return Step{
		ID:      "apple-tools",
		Label:   "Darwin tools and OAM",
		Detail:  "actool is missing",
		Command: "idev setup --apple-tools",
	}
}
