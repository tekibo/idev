package config

import (
	"errors"
	"os"
	"sort"
	"strings"
)

// SourcePlan is one three-way sources migration.
type SourcePlan struct {
	Path      string           `json:"path"`
	Backup    string           `json:"backup,omitempty"`
	Conflicts []SourceConflict `json:"conflicts,omitempty"`
	Warnings  []SourceWarning  `json:"warnings,omitempty"`
	Changed   bool             `json:"changed"`
	local     map[string]any
	merged    map[string]any
	defaults  map[string]any
	applied   string
}

// SourceConflictError reports fields that need confirmation.
type SourceConflictError struct {
	Conflicts []SourceConflict
}

func (e *SourceConflictError) Error() string {
	fields := make([]string, 0, len(e.Conflicts))
	for _, conflict := range e.Conflicts {
		fields = append(fields, conflict.Field)
	}
	return "sources: confirmation needed for " + strings.Join(fields, ", ")
}

// BuildSourcePlan compares local, applied, and embedded sources.
func BuildSourcePlan(home string) (*SourcePlan, error) {
	if os.Getenv("IDEV_SOURCES") != "" {
		return nil, errors.New("sources: migration does not apply to IDEV_SOURCES")
	}
	defaults, err := defaultSourceMap()
	if err != nil {
		return nil, err
	}
	local, localOK, err := loadSourceMap(SourcePath(home))
	if err != nil {
		return nil, err
	}
	if !localOK {
		local = defaults
	}
	applied, appliedOK, err := loadSourceMap(AppliedPath(home))
	if err != nil {
		return nil, err
	}
	if !appliedOK {
		applied = local
	}
	merged, conflicts := mergeSource(flatSource(applied), flatSource(local), flatSource(defaults))
	warnings := sourceWarnings(merged, flatSource(defaults))
	conflicts = appendBlockConflicts(conflicts, warnings)
	return &SourcePlan{
		Path:      SourcePath(home),
		Conflicts: conflicts,
		Warnings:  warnings,
		Changed:   !localOK || !sourceMapEqual(local, unflatSource(merged)),
		local:     flatSource(local),
		merged:    merged,
		defaults:  flatSource(defaults),
		applied:   AppliedPath(home),
	}, nil
}

// AutoMigrate applies a conflict-free migration.
func AutoMigrate(home string) (*SourcePlan, error) {
	if os.Getenv("IDEV_SOURCES") != "" {
		return nil, nil
	}
	plan, err := BuildSourcePlan(home)
	if err != nil {
		return nil, err
	}
	if len(plan.Conflicts) > 0 {
		return plan, &SourceConflictError{Conflicts: plan.Conflicts}
	}
	keep := map[string]bool{}
	for _, warning := range plan.Warnings {
		keep[warning.Field] = true
	}
	_, err = ApplySourcePlan(plan, keep)
	return plan, err
}

func appendBlockConflicts(conflicts []SourceConflict, warnings []SourceWarning) []SourceConflict {
	seen := map[string]bool{}
	for _, conflict := range conflicts {
		seen[conflict.Field] = true
	}
	for _, warning := range warnings {
		if !warning.Block || seen[warning.Field] {
			continue
		}
		conflicts = append(conflicts, SourceConflict{
			Field: warning.Field, Class: pinRequired, Severity: pinBlock,
			Current: warning.Current, Required: warning.Required, Default: warning.Required,
			Reason: warning.Reason,
		})
	}
	sort.Slice(conflicts, func(i, j int) bool { return conflicts[i].Field < conflicts[j].Field })
	return conflicts
}
