package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

// DeviceRun runs the apple-dev device-run script.
func DeviceRun(args []string) int {
	if wantsHelp(args) {
		fmt.Println("Usage: idev device-run [--lldb|--attach|--console [SECONDS]] [args...]")
		return 0
	}
	return runApple("device-run.sh", "device-run", args)
}

// Ship runs the apple-dev ship script.
func Ship(args []string) int {
	script := []string{}
	for index := 0; index < len(args); index++ {
		switch args[index] {
		case "-h", "--help":
			fmt.Println("Usage: idev ship [--upload] [--device [TEAMID]] [--xcode-version V] [--xcode-build B]")
			return 0
		case "--xcode-version":
			index++
			if index >= len(args) {
				fmt.Fprintln(os.Stderr, "idev ship: --xcode-version needs a value")
				return 2
			}
			_ = os.Setenv("XCODE_VERSION", args[index])
		case "--xcode-build":
			index++
			if index >= len(args) {
				fmt.Fprintln(os.Stderr, "idev ship: --xcode-build needs a value")
				return 2
			}
			_ = os.Setenv("XCODE_BUILD", args[index])
		default:
			script = append(script, args[index])
		}
	}
	return runApple("ship.sh", "ship", script)
}

func wantsHelp(args []string) bool {
	return len(args) == 1 && (args[0] == "-h" || args[0] == "--help")
}

func runApple(script, label string, args []string) int {
	manager := manager.New()
	err := manager.RunApple(context.Background(), script, args, os.Stdin, os.Stdout, os.Stderr, JSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev %s: %v\n", label, err)
		return manager.ExitCode(err)
	}
	return 0
}
