// Package next reports the setup state and the next command.
package next

import (
	"os"

	"github.com/tekibo/idev/internal/manager/doctor"
	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/manifest"
	"github.com/tekibo/idev/internal/platform/project"
)

// Step is one setup state line.
type Step struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Done    bool   `json:"done"`
	Detail  string `json:"detail,omitempty"`
	Command string `json:"command,omitempty"`
}

// Plan is the agent-facing setup state.
type Plan struct {
	Version string `json:"version"`
	Home    string `json:"home"`
	Project string `json:"project,omitempty"`
	Kind    string `json:"kind,omitempty"`
	Steps   []Step `json:"steps"`
	Next    string `json:"next,omitempty"`
	Done    bool   `json:"done"`
}

// Build reads the setup state without changing it.
func Build(cfg config.Config, projectDir string) Plan {
	if projectDir == "" {
		projectDir = cfg.Project
	}
	if projectDir == "" {
		projectDir, _ = os.Getwd()
	}
	report := doctor.Build(cfg, projectDir)
	checks := map[string]doctor.Check{}
	for _, check := range report.Checks {
		checks[check.ID] = check
	}
	info := project.Detect(projectDir)
	steps := []Step{
		sourceStep(cfg),
		systemStep(report),
		baseStep(cfg, checks),
		sdkStep(checks),
		appleToolsStep(checks),
		flowsStep(cfg, info),
		agentsStep(cfg),
	}
	projectStep, bundleStep := projectSteps(info)
	steps = append(steps, projectStep)
	if bundleStep != nil {
		steps = append(steps, *bundleStep)
	}
	plan := Plan{
		Version: manifest.Version,
		Home:    cfg.Home,
		Project: projectDir,
		Kind:    project.LabelFor(info.Kind),
		Steps:   steps,
	}
	for _, step := range steps {
		if !step.Done && step.Command != "" {
			plan.Next = step.Command
			return plan
		}
	}
	plan.Done = true
	plan.Next = action(info)
	return plan
}
