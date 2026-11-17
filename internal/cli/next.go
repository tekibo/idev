package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

// Next prints the setup state and the next command.
func Next(args []string) int {
	manager := manager.New()
	asJSON := false
	project := manager.Config.Project
	for index := 0; index < len(args); index++ {
		switch args[index] {
		case "--json":
			asJSON = true
		case "--project":
			index++
			if index >= len(args) {
				fmt.Fprintln(os.Stderr, "idev next: --project needs a directory")
				return 2
			}
			project = args[index]
		case "-h", "--help":
			fmt.Println("Usage: idev next [--json] [--project DIR]")
			return 0
		default:
			fmt.Fprintf(os.Stderr, "idev next: unknown option %q\n", args[index])
			return 2
		}
	}
	plan := manager.NextPlan(project)
	if asJSON {
		data, err := json.MarshalIndent(plan, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "idev next: %v\n", err)
			return 2
		}
		fmt.Println(string(data))
		return 0
	}
	fmt.Printf("idev next — %s\n", plan.Version)
	fmt.Printf("home:    %s\n", plan.Home)
	fmt.Printf("project: %s\n", plan.Project)
	if plan.Kind != "" {
		fmt.Printf("kind:    %s\n", plan.Kind)
	}
	fmt.Println()
	for _, step := range plan.Steps {
		mark := "[ ]"
		if step.Done {
			mark = "[x]"
		}
		fmt.Printf("  %s %s\n", mark, step.Label)
		if step.Detail != "" && !step.Done {
			fmt.Printf("        %s\n", step.Detail)
		}
	}
	if plan.Next != "" {
		fmt.Printf("\nnext: %s\n", plan.Next)
	}
	return 0
}
