package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

// Doctor prints the diagnostic report.
func Doctor(args []string) int {
	manager := manager.New()
	asJSON, private, projectSet := false, false, false
	project := manager.Config.Project
	for index := 0; index < len(args); index++ {
		switch args[index] {
		case "--json":
			asJSON = true
		case "--private":
			private = true
		case "--project":
			index++
			if index >= len(args) {
				fmt.Fprintln(os.Stderr, "idev doctor: --project needs a directory")
				return 2
			}
			project, projectSet = args[index], true
		case "-h", "--help":
			fmt.Println("Usage: idev doctor [--json] [--private] [--project DIR]")
			return 0
		default:
			fmt.Fprintf(os.Stderr, "idev doctor: unknown option %q\n", args[index])
			return 2
		}
	}
	if private && !projectSet {
		project = ""
	}
	report := manager.DoctorFor(project)
	if private {
		report = manager.PrivateDoctorFor(project)
	}
	if asJSON {
		output, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "idev doctor: %v\n", err)
			return 2
		}
		fmt.Println(string(output))
		return exitFor(report.CriticalMissing())
	}
	printDoctor(report)
	return exitFor(report.CriticalMissing())
}

func exitFor(critical bool) int {
	if critical {
		return 1
	}
	return 0
}
