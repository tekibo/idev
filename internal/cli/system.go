package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

// System handles the system plan and apply commands.
func System(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: idev system <plan|apply> [--all]")
		return 2
	}
	switch args[0] {
	case "plan":
		return systemPlan(args[1:])
	case "apply":
		return systemApply(args[1:])
	case "-h", "--help":
		fmt.Println("Usage: idev system <plan|apply> [--all]")
		return 0
	default:
		fmt.Fprintf(os.Stderr, "idev system: unknown command %q\n", args[0])
		return 2
	}
}

func systemPlan(args []string) int {
	asJSON := false
	filtered := make([]string, 0, len(args))
	for _, arg := range args {
		if arg == "--json" {
			asJSON = true
			continue
		}
		filtered = append(filtered, arg)
	}
	all, status := systemOptions(filtered)
	if status >= 0 {
		return status
	}
	plan := manager.New().SystemPlan(all)
	if asJSON {
		data, err := json.MarshalIndent(plan, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "idev system: %v\n", err)
			return 2
		}
		fmt.Println(string(data))
		if plan.Blocked != "" {
			return 1
		}
		return 0
	}
	if len(plan.Packages) == 0 && plan.Blocked == "" {
		fmt.Println("No missing system packages.")
		return 0
	}
	printSystemPlan(plan)
	if plan.Blocked != "" {
		return 1
	}
	return 0
}

func systemApply(args []string) int {
	all, status := systemOptions(args)
	if status >= 0 {
		return status
	}
	manager := manager.New()
	plan := manager.SystemPlan(all)
	if len(plan.Packages) == 0 && plan.Blocked == "" {
		fmt.Println("No missing system packages.")
		return 0
	}
	printSystemPlan(plan)
	if plan.Blocked != "" {
		return 1
	}
	provider := manager.TerminalPasswordProvider(os.Stdin, os.Stderr)
	if JSON {
		provider = stdinPassword{}
	}
	if err := manager.ApplySystemPlan(context.Background(), plan, provider, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "idev system: %v\n", err)
		return 1
	}
	fmt.Println("System packages installed.")
	return 0
}
