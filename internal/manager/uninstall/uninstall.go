package uninstall

import (
	"context"
	"sort"

	"github.com/tekibo/idev/internal/platform/security/changeset"
)

// UninstallReport returns the component-aware kill-switch report.
func (s *Service) UninstallReport(keep []string) (Report, error) {
	packages, err := installedPackages(s.Config.Home)
	if err != nil {
		return Report{}, err
	}
	return Build(s.Config, keep, packages)
}

func installedPackages(home string) ([]string, error) {
	changes, err := changeset.New(home).List()
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	packages := make([]string, 0)
	for _, change := range changes {
		if change.Kind != "system" || change.Status != "applied" {
			continue
		}
		for _, name := range change.Packages {
			if name == "" || seen[name] {
				continue
			}
			seen[name] = true
			packages = append(packages, name)
		}
	}
	sort.Strings(packages)
	return packages, nil
}

// UninstallApply removes the allowlist paths and writes the receipt.
func (s *Service) UninstallApply(ctx context.Context, report Report) (Receipt, error) {
	return Apply(ctx, report)
}
