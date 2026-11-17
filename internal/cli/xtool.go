package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	helptree "github.com/tekibo/idev/internal/cli/help"
	"github.com/tekibo/idev/internal/manager"
)

// Xtool runs any xtool command with the private environment.
func Xtool(args []string) int {
	return runXtool(args)
}

// DS runs xtool ds.
func DS(args []string) int {
	return runXtool(append([]string{"ds"}, args...))
}

// SDK runs xtool sdk, with a parsed status read.
func SDK(args []string) int {
	if len(args) == 2 && args[0] == "status" && wantsHelp(args[1:]) {
		return helptree.Show([]string{"sdk", "status"}, JSON)
	}
	if len(args) == 2 && args[0] == "status" && args[1] == "--json" {
		manager := manager.New()
		status, err := manager.SDKStatus()
		if err != nil {
			fmt.Fprintf(os.Stderr, "idev sdk: %v\n", err)
			return manager.ExitCode(err)
		}
		data, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "idev sdk: %v\n", err)
			return 2
		}
		fmt.Println(string(data))
		return 0
	}
	return runXtool(append([]string{"sdk"}, args...))
}

func runXtool(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: idev xtool <args...>")
		return 2
	}
	manager := manager.New()
	if !manager.Xtool.Available() {
		fmt.Fprintln(os.Stderr, "idev: xtool not found; run idev setup")
		return 127
	}
	env := manager.Config.SwiftEnv()
	if xtoolSDKCommand(args) {
		env = manager.Config.BaseEnv()
	}
	if xtoolBuildCommand(args) {
		env = manager.Config.SwiftBuildEnv()
	}
	if err := manager.RunCommand(context.Background(), manager.Xtool.Path, args, env, "",
		os.Stdin, os.Stdout, os.Stderr, JSON); err != nil {
		fmt.Fprintf(os.Stderr, "idev xtool: %v\n", err)
		return manager.ExitCode(err)
	}
	return 0
}

// xtoolSDKCommand reports the commands that manage the SDK, not consume it.
func xtoolSDKCommand(args []string) bool {
	return len(args) > 0 && args[0] == "sdk"
}

// xtoolBuildCommand reports the commands that need SwiftBuildEnv.
func xtoolBuildCommand(args []string) bool {
	if len(args) < 2 || args[0] != "dev" {
		return false
	}
	switch args[1] {
	case "build", "run", "generate-xcode-project", "build-server":
		return true
	default:
		return false
	}
}
