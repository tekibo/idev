package system

import (
	"sort"
	"strings"
)

// Plan is the exact missing-package transaction.
type Plan struct {
	Manager  string   `json:"manager"`
	Packages []string `json:"packages"`
	Command  string   `json:"command"`
	Name     string   `json:"-"`
	Args     []string `json:"-"`
	Blocked  string   `json:"blocked,omitempty"`
}

// BuildPlan collects missing packages and previews the transaction.
func BuildPlan(manager Manager, checks []Check) Plan {
	if manager == nil {
		return Plan{Blocked: ErrNoManager.Error()}
	}
	seen := map[string]bool{}
	var packages []string
	for _, check := range checks {
		if check.Status != StatusMissing || check.Package == "" || seen[check.Package] {
			continue
		}
		seen[check.Package] = true
		packages = append(packages, check.Package)
	}
	sort.Strings(packages)
	name, args := manager.Install(packages)
	plan := Plan{
		Manager:  manager.Name(),
		Packages: packages,
		Name:     name,
		Args:     args,
		Command:  "sudo " + name + " " + strings.Join(args, " "),
	}
	if len(packages) == 0 {
		return plan
	}
	blocked, err := manager.Preview(packages)
	if err != nil {
		plan.Blocked = "preview failed: " + err.Error()
		return plan
	}
	plan.Blocked = blocked
	return plan
}
