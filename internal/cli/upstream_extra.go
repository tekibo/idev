package cli

import (
	helptree "github.com/tekibo/idev/internal/cli/help"
	"github.com/tekibo/idev/internal/manager"
)

// XcodeProject runs the Xcode project converter.
func XcodeProject(args []string) int {
	if upstreamHelp(args) {
		return helptree.Show([]string{"xcodeproj2xtool"}, JSON)
	}
	manager := manager.New()
	command, err := manager.XcodeProjectCommand(args)
	return runUpstream("xcodeproj2xtool", manager, command, err)
}

// Asc runs the App Store Connect helper.
func Asc(args []string) int {
	if upstreamHelp(args) {
		return helptree.Show([]string{"asc"}, JSON)
	}
	manager := manager.New()
	command, err := manager.AscCommand(args)
	return runUpstream("asc", manager, command, err)
}

// ProvisionDev runs the development profile helper.
func ProvisionDev(args []string) int {
	if upstreamHelp(args) {
		return helptree.Show([]string{"provision-dev"}, JSON)
	}
	manager := manager.New()
	command, err := manager.ProvisionDevCommand(args)
	return runUpstream("provision-dev", manager, command, err)
}

// SignDev runs the development signing helper.
func SignDev(args []string) int {
	if upstreamHelp(args) {
		return helptree.Show([]string{"sign-dev"}, JSON)
	}
	manager := manager.New()
	command, err := manager.SignDevCommand(args)
	return runUpstream("sign-dev", manager, command, err)
}

// SDKFree runs one sdk-free script.
func SDKFree(args []string) int {
	if len(args) == 0 {
		return helptree.Show([]string{"sdk-free"}, JSON)
	}
	if args[0] == "-h" || args[0] == "--help" {
		return helptree.Show([]string{"sdk-free"}, JSON)
	}
	if len(args) == 2 && (args[1] == "-h" || args[1] == "--help") {
		return helptree.Show([]string{"sdk-free", args[0]}, JSON)
	}
	manager := manager.New()
	command, err := manager.SDKFreeCommand(args[0], args[1:])
	return runUpstream("sdk-free", manager, command, err)
}
