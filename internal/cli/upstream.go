package cli

import (
	"context"
	"fmt"
	"os"

	helptree "github.com/tekibo/idev/internal/cli/help"
	"github.com/tekibo/idev/internal/manager"
)

// KMP runs the Kotlin Multiplatform framework script.
func KMP(args []string) int {
	if upstreamHelp(args) {
		return helptree.Show([]string{"kmp"}, JSON)
	}
	manager := manager.New()
	command, err := manager.KMPCommand(args)
	return runUpstream("kmp", manager, command, err)
}

// HostManifests runs the SwiftPM host-manifest helper.
func HostManifests(args []string) int {
	if upstreamHelp(args) {
		return helptree.Show([]string{"host-manifests"}, JSON)
	}
	manager := manager.New()
	command, err := manager.HostManifestsCommand(args)
	return runUpstream("host-manifests", manager, command, err)
}

// SimView runs one simview script.
func SimView(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: idev simview <build|test|test-m3|test-m4|test-m5|test-m6|test-m7> [args...]")
		return 2
	}
	if args[0] == "-h" || args[0] == "--help" {
		return helptree.Show([]string{"simview"}, JSON)
	}
	if len(args) == 2 && (args[1] == "-h" || args[1] == "--help") {
		return helptree.Show([]string{"simview", args[0]}, JSON)
	}
	manager := manager.New()
	command, err := manager.SimViewCommand(args[0], args[1:])
	return runUpstream("simview", manager, command, err)
}

// MachoLoad runs the machoload builder.
func MachoLoad(args []string) int {
	if upstreamHelp(args) {
		return helptree.Show([]string{"machoload"}, JSON)
	}
	manager := manager.New()
	command, err := manager.MachoLoadCommand(args)
	return runUpstream("machoload", manager, command, err)
}

// CompileAppIcon runs the app icon compiler.
func CompileAppIcon(args []string) int {
	if upstreamHelp(args) {
		return helptree.Show([]string{"compile-app-icon"}, JSON)
	}
	manager := manager.New()
	command, err := manager.CompileAppIconCommand(args)
	return runUpstream("compile-app-icon", manager, command, err)
}

func upstreamHelp(args []string) bool {
	return len(args) == 1 && (args[0] == "-h" || args[0] == "--help")
}

func runUpstream(label string, manager *manager.App, command manager.UpstreamCommand, err error) int {
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev %s: %v\n", label, err)
		return 2
	}
	err = manager.RunCommand(context.Background(), command.Name, command.Args, command.Env, "",
		os.Stdin, os.Stdout, os.Stderr, JSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev %s: %v\n", label, err)
		return manager.ExitCode(err)
	}
	return 0
}
