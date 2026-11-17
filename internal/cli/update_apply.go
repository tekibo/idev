package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

func updateApply(args []string) int {
	all := false
	components := []string{}
	for _, arg := range args {
		switch {
		case arg == "--all":
			all = true
		case arg == "-h" || arg == "--help":
			fmt.Println("Usage: idev update apply [--all] [component]")
			return 0
		case len(arg) > 0 && arg[0] == '-':
			fmt.Fprintf(os.Stderr, "idev update apply: unknown option %q\n", arg)
			return 2
		default:
			components = append(components, arg)
		}
	}
	if all && len(components) > 0 {
		fmt.Fprintln(os.Stderr, "idev update apply: --all runs alone")
		return 2
	}
	if !all && len(components) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: idev update apply [--all] [component]")
		return 2
	}
	if code := autoMigrateSources(); code != 0 {
		return code
	}
	manager := manager.New()
	if all {
		if err := manager.UpdateApplyAll(context.Background(), os.Stdout, os.Stderr); err != nil {
			fmt.Fprintf(os.Stderr, "idev update: %v\n", err)
			return 1
		}
		printUpdateHints(manager)
		fmt.Println("updates applied; restart idev")
		return 0
	}
	for _, component := range components {
		if err := manager.UpdateApply(context.Background(), component, os.Stdout, os.Stderr); err != nil {
			fmt.Fprintf(os.Stderr, "idev update: %v\n", err)
			return 1
		}
		fmt.Printf("updated %s\n", component)
	}
	printUpdateHints(manager)
	return 0
}
