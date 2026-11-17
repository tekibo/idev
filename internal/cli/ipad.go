package cli

import (
	helptree "github.com/tekibo/idev/internal/cli/help"
	"github.com/tekibo/idev/internal/manager"
)

// IpadScreen runs the iPad viewer lifecycle.
func IpadScreen(args []string) int {
	if len(args) == 0 {
		return helptree.Show([]string{"ipad-screen"}, JSON)
	}
	if args[0] == "-h" || args[0] == "--help" {
		return helptree.Show([]string{"ipad-screen"}, JSON)
	}
	if len(args) == 2 && (args[1] == "-h" || args[1] == "--help") {
		return helptree.Show([]string{"ipad-screen", args[0]}, JSON)
	}
	manager := manager.New()
	command, err := manager.IpadScreenCommand(args)
	return runUpstream("ipad-screen", manager, command, err)
}
