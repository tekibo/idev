package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

// Project prints project detection data.
func Project(args []string) int {
	if len(args) == 1 && args[0] == "type" {
		fmt.Println(manager.New().ProjectType())
		return 0
	}
	if len(args) >= 1 && args[0] == "info" {
		return projectInfo(args[1:])
	}
	if len(args) >= 1 && args[0] == "apps" {
		return projectApps(args[1:])
	}
	if len(args) >= 2 && args[0] == "set-bundle-id" {
		return projectSetBundleID(args[1:])
	}
	if len(args) == 1 && args[0] == "clean" {
		manager := manager.New()
		if _, err := manager.CleanBuild(manager.Config.Project, func(line string) { fmt.Println(line) }); err != nil {
			fmt.Fprintf(os.Stderr, "idev project clean: %v\n", err)
			return 1
		}
		return 0
	}
	if len(args) == 1 && wantsHelp(args) {
		fmt.Println("Usage: idev project <type|info|apps|set-bundle-id|clean> [options]")
		return 0
	}
	fmt.Fprintln(os.Stderr, "Usage: idev project <type|info|apps|set-bundle-id|clean> [options]")
	return 2
}

func projectInfo(args []string) int {
	asJSON := false
	dir := manager.New().Config.Project
	appName := ""
	for index := 0; index < len(args); index++ {
		switch args[index] {
		case "--json":
			asJSON = true
		case "--project":
			index++
			if index >= len(args) {
				fmt.Fprintln(os.Stderr, "idev project info: --project needs a directory")
				return 2
			}
			dir = args[index]
		case "--app":
			index++
			if index >= len(args) {
				fmt.Fprintln(os.Stderr, "idev project info: --app needs a name or directory")
				return 2
			}
			appName = args[index]
		default:
			fmt.Fprintln(os.Stderr, "Usage: idev project info [--json] [--project DIR] [--app NAME|DIR]")
			return 2
		}
	}
	manager := manager.New()
	info := manager.Inspect(dir)
	workspace := manager.Workspace(dir)
	if appName != "" {
		selected, ok := workspace.App(appName)
		if !ok {
			fmt.Fprintf(os.Stderr, "idev project info: app %q not found\n", appName)
			printWorkspaceApps(os.Stderr, workspace)
			return 2
		}
		info = selected
	} else if !info.Known() && len(workspace.Apps) == 1 {
		info = workspace.Apps[0]
	} else if !info.Known() && len(workspace.Apps) > 1 {
		fmt.Fprintf(os.Stderr, "idev project info: %d apps found; use --app\n", len(workspace.Apps))
		printWorkspaceApps(os.Stderr, workspace)
		return 2
	}
	if asJSON {
		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "idev project info: %v\n", err)
			return 2
		}
		fmt.Println(string(data))
		return 0
	}
	fmt.Printf("kind: %s\nlabel: %s\nname: %s\ndir: %s\nbundle: %s\nproduct: %s\ndeployment: %s\nworkflow: %s\nscheme: %s\napp-entry: %t\n",
		info.Kind, info.Label, info.Name, info.Dir, info.BundleID, info.Product,
		info.DeploymentTarget, info.Workflow, info.Scheme, info.AppEntry)
	return 0
}
