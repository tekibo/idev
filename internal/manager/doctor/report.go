package doctor

import (
	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
	"github.com/tekibo/idev/internal/platform/manifest"
	"github.com/tekibo/idev/internal/platform/system"
)

// Report is the full diagnostic snapshot.
type Report struct {
	Version    string         `json:"version"`
	Home       string         `json:"home"`
	Manager    string         `json:"manager"`
	Distro     string         `json:"distro"`
	Checks     []Check        `json:"checks"`
	Summary    map[string]int `json:"summary"`
	SystemPlan system.Plan    `json:"system_plan"`
}

// Build creates the report.
func Build(cfg config.Config, project string) Report {
	runner := execx.Exec{}
	probe := system.New(runner)
	manager, managerErr := system.Detect(runner)
	checks := system.Checks(probe, manager, false)
	if managerErr != nil {
		checks = append(checks, system.ManagerCheck(managerErr))
	}
	checks = append(checks, privateChecks(cfg, project)...)
	managerName := "unsupported"
	distro := system.DistroName()
	if manager != nil {
		managerName = manager.Name()
		distro = manager.Distro()
	}
	report := Report{
		Version:    manifest.Version,
		Home:       cfg.Home,
		Manager:    managerName,
		Distro:     distro,
		Checks:     checks,
		Summary:    map[string]int{},
		SystemPlan: system.BuildPlan(manager, checks),
	}
	for _, check := range checks {
		report.Summary[check.Status]++
	}
	return report
}

// BuildPrivate creates a report without the system package probes.
func BuildPrivate(cfg config.Config, project string) Report {
	checks := privateChecks(cfg, project)
	report := Report{
		Version: manifest.Version,
		Home:    cfg.Home,
		Checks:  checks,
		Summary: map[string]int{},
	}
	for _, check := range checks {
		report.Summary[check.Status]++
	}
	return report
}

// CriticalMissing reports whether a critical check failed.
func (r Report) CriticalMissing() bool {
	for _, check := range r.Checks {
		if check.Critical && !check.OK() {
			return true
		}
	}
	return false
}
