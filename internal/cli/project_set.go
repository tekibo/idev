package cli

import (
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

func projectSetBundleID(args []string) int {
	manager := manager.New()
	dir := manager.Config.Project
	id := ""
	for index := 0; index < len(args); index++ {
		switch args[index] {
		case "--project":
			index++
			if index >= len(args) {
				fmt.Fprintln(os.Stderr, "idev project set-bundle-id: --project needs a directory")
				return 2
			}
			dir = args[index]
		case "--value":
			index++
			if index >= len(args) {
				fmt.Fprintln(os.Stderr, "idev project set-bundle-id: --value needs an identifier")
				return 2
			}
			id = args[index]
		default:
			fmt.Fprintln(os.Stderr, "Usage: idev project set-bundle-id [--project DIR] --value ID")
			return 2
		}
	}
	if id == "" {
		fmt.Fprintln(os.Stderr, "idev project set-bundle-id: --value is required")
		return 2
	}
	if err := manager.SetBundleID(dir, id); err != nil {
		fmt.Fprintf(os.Stderr, "idev project set-bundle-id: %v\n", err)
		return 1
	}
	fmt.Println("bundle id set to " + id)
	return 0
}
