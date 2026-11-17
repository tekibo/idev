package manager

import (
	"github.com/tekibo/idev/internal/platform/config"
)

// SourcePlan is one sources migration plan.
type SourcePlan = config.SourcePlan

// SourceConflict is one source field conflict.
type SourceConflict = config.SourceConflict

// SourceWarning is one unsupported required pin.
type SourceWarning = config.SourceWarning

// BuildSourcePlan reads the local, applied, and embedded sources.
func (a *App) BuildSourcePlan() (*SourcePlan, error) {
	return config.BuildSourcePlan(a.Config.Home)
}

// AutoMigrateSources applies a conflict-free sources migration.
func (a *App) AutoMigrateSources() (*SourcePlan, error) {
	return config.AutoMigrate(a.Config.Home)
}

// ApplySourcePlan writes the chosen sources values.
func (a *App) ApplySourcePlan(plan *SourcePlan, keep map[string]bool) (string, error) {
	return config.ApplySourcePlan(plan, keep)
}

// RefreshSources resets the source config to the embedded defaults.
func (a *App) RefreshSources() (string, error) {
	return config.RefreshSources(a.Config.Home)
}
