package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

func projectApps(args []string) int {
	asJSON := false
	dir := manager.New().Config.Project
	for index := 0; index < len(args); index++ {
		switch args[index] {
		case "--json":
			asJSON = true
		case "--project":
			index++
			if index >= len(args) {
				fmt.Fprintln(os.Stderr, "idev project apps: --project needs a directory")
				return 2
			}
			dir = args[index]
		default:
			fmt.Fprintln(os.Stderr, "Usage: idev project apps [--json] [--project DIR]")
			return 2
		}
	}
	workspace := manager.New().Workspace(dir)
	if asJSON {
		data, err := json.MarshalIndent(workspace, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "idev project apps: %v\n", err)
			return 2
		}
		fmt.Println(string(data))
		return 0
	}
	printWorkspaceApps(os.Stdout, workspace)
	return 0
}

func printWorkspaceApps(w io.Writer, workspace manager.Workspace) {
	fmt.Fprintf(w, "workspace: %s\n", workspace.Dir)
	fmt.Fprintf(w, "apps: %d\n", len(workspace.Apps))
	for _, item := range workspace.Apps {
		label := item.Label
		if item.Workflow != "" {
			label += ", " + item.Workflow
		}
		fmt.Fprintf(w, "  - %s (%s)\n", item.Name, label)
		fmt.Fprintf(w, "    dir: %s\n", item.Dir)
		if item.BundleID != "" {
			if item.BundleIDSource != "" {
				fmt.Fprintf(w, "    bundle: %s (%s)\n", item.BundleID, item.BundleIDSource)
			} else {
				fmt.Fprintf(w, "    bundle: %s\n", item.BundleID)
			}
		}
		if item.Scheme != "" {
			fmt.Fprintf(w, "    scheme: %s\n", item.Scheme)
		}
	}
}
