package cli

import (
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

// printUninstallReport prints the typed uninstall report.
func printUninstallReport(report manager.UninstallServiceReport) {
	fmt.Println("Components:")
	for _, component := range report.Components {
		action := "remove"
		if component.Keep {
			action = "keep"
		}
		kind := "fixed"
		if component.Keepable {
			kind = "keepable"
		}
		fmt.Printf("  %-6s %-16s %-8s %s\n", action, component.ID, kind, component.Label)
	}
	fmt.Println("Will remove:")
	for _, path := range report.Remove {
		fmt.Printf("  %s\n", path)
	}
	fmt.Println("Will keep:")
	for _, path := range report.Keep {
		fmt.Printf("  %s\n", path)
	}
	if len(report.Keep) > 0 {
		fmt.Println("This is a partial uninstall. Kept tools remain on disk and will no longer be managed by idev.")
	}
	if len(report.System) > 0 {
		fmt.Println("System packages installed by idev (not removed):")
		for _, item := range report.System {
			fmt.Printf("  %s\n", item.Package)
		}
	}
	if report.SystemErr != "" {
		fmt.Fprintf(os.Stderr, "changeset: %s\n", report.SystemErr)
	}
	fmt.Println("User configs (preserved):")
	for _, item := range report.Configs {
		fmt.Printf("  %s: %s\n", item.Path, item.Contents)
	}
	fmt.Println("Project data (preserved):")
	for _, item := range report.ProjectData {
		fmt.Printf("  %s\n", item)
	}
}

// printRemaining prints unexpected paths left inside the prefix.
func printRemaining(paths []string) {
	fmt.Println("Remaining inside the prefix (not removed):")
	for _, path := range paths {
		fmt.Printf("  %s\n", path)
	}
}
