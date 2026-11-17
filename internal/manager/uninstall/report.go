// Package uninstall builds and applies the idev kill switch.
package uninstall

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/tekibo/idev/internal/platform/config"
)

// Config is one preserved user config.
type Config struct {
	Path     string `json:"path"`
	Contents string `json:"contents"`
}

// SystemPackage is one package installed by an idev transaction.
type SystemPackage struct {
	Package string `json:"package"`
}

// Report is the pre-uninstall report.
type Report struct {
	Home        string          `json:"home"`
	Remove      []string        `json:"remove"`
	Keep        []string        `json:"keep"`
	Containers  []string        `json:"containers"`
	Components  []Component     `json:"components"`
	System      []SystemPackage `json:"system"`
	Configs     []Config        `json:"configs"`
	ProjectData []string        `json:"project_data"`
	SystemErr   string          `json:"system_error,omitempty"`
}

// Build returns the component-aware allowlist report.
func Build(cfg config.Config, requested []string, installed []string) (Report, error) {
	if cfg.HomeErr != nil {
		return Report{}, cfg.HomeErr
	}
	if err := config.ValidateHome(cfg.Home); err != nil {
		return Report{}, fmt.Errorf("uninstall: %w", err)
	}
	userHome := config.UserHome()
	report := Report{
		Home: cfg.Home,
		Containers: []string{
			cfg.Path("lib", "xtool"),
			cfg.Path("lib"),
			cfg.Path("bin"),
			cfg.Path("toolchains"),
			cfg.Path(".config", "swiftpm", "swift-sdks"),
			cfg.Path(".config", "swiftpm"),
			cfg.Path(".config"),
		},
		Components: buildComponents(cfg),
	}
	keep, err := resolveKeep(report.Components, requested)
	if err != nil {
		return Report{}, err
	}
	for index := range report.Components {
		component := &report.Components[index]
		component.Keep = keep[component.ID]
		if component.Keep {
			report.Keep = appendUnique(report.Keep, component.Paths...)
		} else {
			report.Remove = appendUnique(report.Remove, component.Paths...)
		}
	}
	shared := existing([]string{cfg.Path("lib", "compat")})
	if keep["swift"] || keep["xtool"] {
		report.Keep = appendUnique(report.Keep, shared...)
	} else {
		report.Remove = appendUnique(report.Remove, shared...)
	}
	link := filepath.Join(userHome, ".local", "bin", "idev")
	if ownedLink(link, cfg.Home) || ownedExecutable(link) {
		report.Remove = append(report.Remove, link)
	}
	sort.Strings(report.Remove)
	report.Configs = userConfigs(userHome)
	report.ProjectData = projectData()
	report.System = systemPackages(installed)
	return report, nil
}
