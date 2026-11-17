package cli

import (
	"fmt"
	"os"

	helptree "github.com/tekibo/idev/internal/cli/help"
	"github.com/tekibo/idev/internal/manager"
)

// Agents installs or removes the idev agent skill.
func Agents(args []string) int {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		return helptree.Show([]string{"agents"}, JSON)
	}
	if len(args) == 2 && (args[1] == "-h" || args[1] == "--help") {
		return helptree.Show([]string{"agents", args[0]}, JSON)
	}
	switch args[0] {
	case "install":
		return agentsInstall(args[1:])
	case "uninstall":
		return agentsUninstall(args[1:])
	case "status":
		return agentsStatus(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "idev agents: unknown command %q\n", args[0])
		return 2
	}
}

func agentsInstall(args []string) int {
	asJSON, ok := agentsJSON(args)
	if !ok {
		return 2
	}
	result, err := manager.New().InstallAgents()
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev agents: %v\n", err)
		return 1
	}
	if asJSON {
		return agentsPrint(result)
	}
	fmt.Printf("skill:    %s\n", result.SkillDir)
	for _, link := range result.Links {
		fmt.Printf("linked:   %s\n", link)
	}
	for _, path := range result.Instructions {
		fmt.Printf("instructions: %s\n", path)
	}
	for _, path := range result.Skipped {
		fmt.Printf("skipped:  %s\n", path)
	}
	return 0
}

func agentsUninstall(args []string) int {
	asJSON, ok := agentsJSON(args)
	if !ok {
		return 2
	}
	result, err := manager.New().UninstallAgents()
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev agents: %v\n", err)
		return 1
	}
	if asJSON {
		return agentsPrint(result)
	}
	for _, link := range result.RemovedLinks {
		fmt.Printf("removed:  %s\n", link)
	}
	if result.RemovedSkill {
		fmt.Println("removed:  skill directory")
	}
	for _, path := range result.RemovedInstructions {
		fmt.Printf("removed:  instructions in %s\n", path)
	}
	return 0
}

func agentsStatus(args []string) int {
	asJSON, ok := agentsJSON(args)
	if !ok {
		return 2
	}
	status := manager.New().AgentStatus()
	if asJSON {
		return agentsPrint(status)
	}
	if status.Installed && status.Stale {
		fmt.Printf("skill:    %s (stale; run idev agents install)\n", status.SkillDir)
	} else if status.Installed {
		fmt.Printf("skill:    %s (installed)\n", status.SkillDir)
	} else {
		fmt.Printf("skill:    %s (missing)\n", status.SkillDir)
	}
	for _, link := range status.Links {
		state := "missing"
		if link.OK {
			state = "ok"
		} else if link.Detail != "" {
			state = link.Detail
		}
		fmt.Printf("link:     %s (%s)\n", link.Path, state)
	}
	for _, item := range status.Instructions {
		state := "missing"
		if item.Present {
			state = "present"
		}
		fmt.Printf("rules:    %s (%s)\n", item.Path, state)
	}
	return 0
}
