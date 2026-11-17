package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/tekibo/idev/internal/manager"
)

// Uninstall prints the report and runs the kill switch after typed confirmation.
func Uninstall(args []string) int {
	dryRun, asJSON, latest := false, false, false
	confirm := ""
	var keep []string
	for index := 0; index < len(args); index++ {
		arg := args[index]
		switch {
		case arg == "--dry-run":
			dryRun = true
		case arg == "--json":
			asJSON = true
		case arg == "--latest":
			latest = true
		case arg == "--confirm":
			index++
			if index >= len(args) {
				return uninstallUsage("--confirm needs a value")
			}
			confirm = args[index]
		case arg == "--keep":
			index++
			if index >= len(args) {
				return uninstallUsage("--keep needs a component id")
			}
			keep = append(keep, splitKeep(args[index])...)
		case strings.HasPrefix(arg, "--keep="):
			keep = append(keep, splitKeep(strings.TrimPrefix(arg, "--keep="))...)
		case arg == "-h" || arg == "--help":
			fmt.Println("Usage: idev uninstall [--dry-run|--latest] [--json] [--keep <id>]... [--confirm \"remove idev\"]")
			return 0
		default:
			return uninstallUsage("unknown option " + arg)
		}
	}
	manager := manager.New()
	if latest {
		receipt, err := manager.LatestReceipt()
		if err != nil {
			fmt.Fprintf(os.Stderr, "idev uninstall: %v\n", err)
			return 1
		}
		if asJSON {
			return printJSON(receipt)
		}
		fmt.Println(receipt.Path)
		return 0
	}
	report, err := manager.UninstallReport(keep)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev uninstall: %v\n", err)
		return 2
	}
	if dryRun {
		if asJSON {
			return printJSON(report)
		}
		printUninstallReport(report)
		return 0
	}
	if asJSON {
		return uninstallUsage("--json applies only to --dry-run and --latest")
	}
	printUninstallReport(report)
	if confirm == "" {
		fmt.Print("Type remove idev to continue: ")
		line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		confirm = strings.TrimSpace(line)
	}
	if confirm != "remove idev" {
		fmt.Fprintln(os.Stderr, "idev uninstall: confirmation phrase did not match")
		return 2
	}
	receipt, err := manager.UninstallApply(context.Background(), report)
	if len(receipt.Remaining) > 0 {
		printRemaining(receipt.Remaining)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev uninstall: %v\n", err)
		if _, statErr := os.Stat(receipt.Path); statErr == nil {
			fmt.Printf("receipt: %s\n", receipt.Path)
		}
		return 1
	}
	if len(receipt.KeptComponents) > 0 {
		fmt.Printf("kept: %s\n", strings.Join(receipt.KeptComponents, ", "))
	}
	fmt.Printf("receipt: %s\n", receipt.Path)
	return 0
}

func printJSON(value any) int {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev uninstall: %v\n", err)
		return 2
	}
	fmt.Println(string(data))
	return 0
}
