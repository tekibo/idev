// Package official installs pinned official binary releases.
package official

import (
	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/release"
)

// Spec is one official binary install plan.
type Spec struct {
	Name        string
	Version     string
	Source      string
	Bin         string
	VersionArgs []string
	Component   release.Component
}

// NewSpec returns the official binary install plan.
func NewSpec(cfg config.Config, name, source string, component release.Component, versionArgs []string) Spec {
	return Spec{
		Name:        name,
		Version:     component.Version,
		Source:      source,
		Bin:         cfg.Path("bin", name),
		VersionArgs: versionArgs,
		Component:   component,
	}
}
