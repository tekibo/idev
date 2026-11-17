package manager

import "github.com/tekibo/idev/internal/manager/doctor"

// Report is the diagnostic report exposed by the app layer.
type Report = doctor.Report

// DoctorFor returns the full diagnostic report for one project directory.
func (a *App) DoctorFor(project string) Report {
	return doctor.Build(a.Config, project)
}

// PrivateDoctorFor returns the private and project checks without system probes.
func (a *App) PrivateDoctorFor(project string) Report {
	return doctor.BuildPrivate(a.Config, project)
}
