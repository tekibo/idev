package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

// Update handles the manager and component update commands.
func Update(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: idev update <check|apply|clean> [--all] [component]")
		return 2
	}
	switch args[0] {
	case "check":
		return updateCheck(args[1:])
	case "apply":
		return updateApply(args[1:])
	case "clean":
		return updateClean(args[1:])
	case "-h", "--help":
		fmt.Println("Usage: idev update <check|apply|clean> [--all] [component]")
		return 0
	default:
		fmt.Fprintf(os.Stderr, "idev update: unknown command %q\n", args[0])
		return 2
	}
}

func updateCheck(args []string) int {
	asJSON := false
	for _, arg := range args {
		switch arg {
		case "--json":
			asJSON = true
		case "-h", "--help":
			fmt.Println("Usage: idev update check [--json]")
			return 0
		default:
			fmt.Fprintf(os.Stderr, "idev update check: unknown option %q\n", arg)
			return 2
		}
	}
	if code := autoMigrateSources(); code != 0 {
		return code
	}
	updates, err := manager.New().UpdateCheck()
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev update: %v\n", err)
		return 1
	}
	if asJSON {
		data, err := json.MarshalIndent(updates, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "idev update: %v\n", err)
			return 2
		}
		fmt.Println(string(data))
		return 0
	}
	fmt.Printf("%-18s %-20s %-20s %s\n", "component", "installed", "available", "status")
	for _, update := range updates {
		status := "update"
		if update.Current {
			status = "current"
		}
		if update.Warning != "" {
			status = "unsupported"
		}
		fmt.Printf("%-18s %-20s %-20s %s\n", update.Component, update.Installed, update.Available, status)
	}
	return 0
}
