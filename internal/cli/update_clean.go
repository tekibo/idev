package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	helptree "github.com/tekibo/idev/internal/cli/help"
	"github.com/tekibo/idev/internal/manager"
)

func updateClean(args []string) int {
	dryRun, yes := false, false
	asJSON := JSON
	for _, arg := range args {
		switch arg {
		case "--dry-run":
			dryRun = true
		case "--yes":
			yes = true
		case "--json":
			asJSON = true
		case "-h", "--help":
			return helptree.Show([]string{"update", "clean"}, JSON)
		default:
			fmt.Fprintf(os.Stderr, "idev update clean: unknown option %q\n", arg)
			return 2
		}
	}
	manager := manager.New()
	plan := manager.CleanPlan()
	if asJSON {
		if dryRun || plan.Bytes == 0 {
			return cleanJSON(plan)
		}
		if !yes {
			fmt.Fprintln(os.Stderr, "idev update clean: --json needs --dry-run or --yes")
			return 2
		}
		result, err := manager.CleanApply(plan)
		if err != nil {
			fmt.Fprintf(os.Stderr, "idev update clean: %v\n", err)
			return 1
		}
		return cleanJSON(result)
	}
	if plan.Bytes == 0 {
		fmt.Println("nothing to clean")
		return 0
	}
	fmt.Printf("clean: %s in %d paths\n", cleanSize(plan.Bytes), len(plan.Targets))
	for _, target := range plan.Targets {
		fmt.Printf("  %s (%s)\n", target.Path, cleanSize(target.Bytes))
	}
	if dryRun {
		return 0
	}
	if !yes {
		if !manager.IsTerminal(os.Stdin.Fd()) {
			fmt.Fprintln(os.Stderr, "idev update clean: no terminal for confirmation; rerun with --yes")
			return 2
		}
		fmt.Print("Remove these caches? [y/N] ")
		line, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil && strings.TrimSpace(line) == "" {
			fmt.Println("cancelled")
			return 0
		}
		if !cleanConfirmed(line) {
			fmt.Println("cancelled")
			return 0
		}
	}
	result, err := manager.CleanApply(plan)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev update clean: %v\n", err)
		return 1
	}
	fmt.Printf("removed %s from %d paths\n", cleanSize(result.Bytes), len(result.Removed))
	return 0
}

func cleanJSON(value any) int {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev update clean: %v\n", err)
		return 2
	}
	fmt.Println(string(data))
	return 0
}
