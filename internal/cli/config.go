package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tekibo/idev/internal/manager"
)

// ConfigCommand handles source migration and refresh.
func ConfigCommand(args []string) int {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Fprintln(os.Stderr, "Usage: idev config <migrate|refresh> [options]")
		return 0
	}
	switch args[0] {
	case "migrate":
		return sourceMigrate(args[1:])
	case "refresh":
		return sourceRefresh(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "idev config: unknown command %q\n", args[0])
		return 2
	}
}

func sourceRefresh(args []string) int {
	for _, arg := range args {
		switch arg {
		case "-h", "--help":
			fmt.Println("Usage: idev config refresh")
			return 0
		default:
			fmt.Fprintf(os.Stderr, "idev config refresh: unknown option %q\n", arg)
			return 2
		}
	}
	manager := manager.New()
	if manager.Config.HomeErr != nil {
		fmt.Fprintf(os.Stderr, "idev config: %v\n", manager.Config.HomeErr)
		return 1
	}
	backup, err := manager.RefreshSources()
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev config: %v\n", err)
		return 1
	}
	if backup != "" {
		fmt.Printf("backup: %s\n", backup)
	}
	fmt.Println("sources refreshed")
	return 0
}

func printSourcePlan(plan *manager.SourcePlan) {
	fmt.Printf("sources: %s\n", plan.Path)
	for _, conflict := range plan.Conflicts {
		fmt.Printf("conflict %s: current %s, default %s", conflict.Field, conflict.Current, conflict.Default)
		if conflict.Required != "" {
			fmt.Printf(", required %s", conflict.Required)
		}
		if conflict.Reason != "" {
			fmt.Printf(" (%s)", conflict.Reason)
		}
		fmt.Println()
	}
	seen := conflictFields(plan)
	for _, warning := range plan.Warnings {
		if seen[warning.Field] {
			continue
		}
		fmt.Printf("warning %s: current %s, required %s (%s)\n",
			warning.Field, warning.Current, warning.Required, warning.Reason)
	}
	if !plan.Changed && len(plan.Conflicts) == 0 && len(plan.Warnings) == 0 {
		fmt.Println("no source changes")
	}
}

func conflictFields(plan *manager.SourcePlan) map[string]bool {
	fields := map[string]bool{}
	for _, conflict := range plan.Conflicts {
		fields[conflict.Field] = true
	}
	return fields
}

func askKeep(reader *bufio.Reader, field, required, current string) bool {
	fmt.Printf("%s: use %s instead of %s? [Y/n] ", field, required, current)
	line, err := reader.ReadString('\n')
	if err != nil && strings.TrimSpace(line) == "" {
		return false
	}
	answer := strings.ToLower(strings.TrimSpace(line))
	return answer == "" || answer == "y" || answer == "yes"
}
