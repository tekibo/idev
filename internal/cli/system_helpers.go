package cli

import (
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

func systemOptions(args []string) (bool, int) {
	all := false
	for _, arg := range args {
		switch arg {
		case "--all":
			all = true
		case "-h", "--help":
			fmt.Println("Usage: idev system <plan|apply> [--all]")
			return false, 0
		default:
			fmt.Fprintf(os.Stderr, "idev system: unknown option %q\n", arg)
			return false, 2
		}
	}
	return all, -1
}

func printSystemPlan(plan manager.SystemPlan) {
	if plan.Manager != "" {
		fmt.Printf("manager: %s\n", plan.Manager)
	}
	if len(plan.Packages) > 0 {
		fmt.Println("missing packages:")
		for _, pkg := range plan.Packages {
			fmt.Printf("  %s\n", pkg)
		}
	}
	if plan.Command != "" {
		fmt.Printf("\n%s\n", plan.Command)
	}
	if plan.Blocked != "" {
		fmt.Fprintf(os.Stderr, "blocked: %s\n", plan.Blocked)
	}
}
