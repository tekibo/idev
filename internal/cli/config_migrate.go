package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

func sourceMigrate(args []string) int {
	yes, dryRun := false, false
	for _, arg := range args {
		switch arg {
		case "--yes":
			yes = true
		case "--dry-run":
			dryRun = true
		case "-h", "--help":
			fmt.Println("Usage: idev config migrate [--yes] [--dry-run]")
			return 0
		default:
			fmt.Fprintf(os.Stderr, "idev config migrate: unknown option %q\n", arg)
			return 2
		}
	}
	manager := manager.New()
	if manager.Config.HomeErr != nil {
		fmt.Fprintf(os.Stderr, "idev config: %v\n", manager.Config.HomeErr)
		return 1
	}
	plan, err := manager.BuildSourcePlan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev config: %v\n", err)
		return 1
	}
	printSourcePlan(plan)
	if dryRun {
		return 0
	}
	keep := map[string]bool{}
	if !yes && (len(plan.Conflicts) > 0 || len(plan.Warnings) > 0) {
		if !manager.IsTerminal(os.Stdin.Fd()) {
			fmt.Fprintln(os.Stderr, "idev config migrate: no terminal for confirmation; rerun with --yes")
			return 2
		}
		reader := bufio.NewReader(os.Stdin)
		for _, conflict := range plan.Conflicts {
			value := conflict.Default
			if conflict.Class == "required" {
				value = conflict.Required
			}
			keep[conflict.Field] = !askKeep(reader, conflict.Field, value, conflict.Current)
		}
		seen := conflictFields(plan)
		for _, warning := range plan.Warnings {
			if seen[warning.Field] {
				continue
			}
			keep[warning.Field] = !askKeep(reader, warning.Field, warning.Required, warning.Current)
		}
	}
	backup, err := manager.ApplySourcePlan(plan, keep)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev config: %v\n", err)
		return 1
	}
	if backup != "" {
		fmt.Printf("backup: %s\n", backup)
	}
	fmt.Println("sources migrated")
	return 0
}

func setupMigrate(dryRun bool) bool {
	if dryRun {
		return true
	}
	return autoMigrateSources() == 0
}

func autoMigrateSources() int {
	manager := manager.New()
	if manager.Config.HomeErr != nil {
		fmt.Fprintf(os.Stderr, "idev: %v\n", manager.Config.HomeErr)
		return 1
	}
	plan, err := manager.AutoMigrateSources()
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev: %v; run idev config migrate\n", err)
		return 1
	}
	if plan == nil {
		return 0
	}
	for _, warning := range plan.Warnings {
		fmt.Fprintf(os.Stderr, "warning: %s is %s, required %s\n",
			warning.Field, warning.Current, warning.Required)
	}
	return 0
}
