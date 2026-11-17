package project

import (
	proj "github.com/tekibo/idev/internal/platform/project"
)

// ProjectType returns the detected project kind.
func (s *Service) ProjectType() proj.Kind {
	return proj.Detect(s.Config.Project).Kind
}

// Inspect returns the detected project fields.
func (s *Service) Inspect(dir string) proj.Info {
	return proj.Inspect(dir)
}

// Workspace returns the apps detected under one directory.
func (s *Service) Workspace(dir string) proj.Workspace {
	return proj.DetectWorkspace(dir)
}

// SetBundleID updates the top-level bundleID in xtool.yml.
func (s *Service) SetBundleID(dir, id string) error {
	return proj.WriteBundleID(dir, id)
}
