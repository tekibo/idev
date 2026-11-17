package cli

import (
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

// LaunchApp launches an installed app on one device.
func LaunchApp(args []string) int {
	bundleID, udid, code := bundleTarget(args, "Usage: idev launch-app --bundle-id <id> [--udid <udid>]")
	if code >= 0 {
		return code
	}
	if bundleID == "" {
		fmt.Fprintln(os.Stderr, "idev launch-app: --bundle-id is required")
		return 2
	}
	manager := manager.New()
	if err := manager.LaunchApp(bundleID, udid, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "idev launch-app: %v\n", err)
		return manager.ExitCode(err)
	}
	return 0
}

// UninstallApp removes an installed app from one device.
func UninstallApp(args []string) int {
	bundleID, udid, code := bundleTarget(args, "Usage: idev uninstall-app --bundle-id <id> [--udid <udid>]")
	if code >= 0 {
		return code
	}
	if bundleID == "" {
		fmt.Fprintln(os.Stderr, "idev uninstall-app: --bundle-id is required")
		return 2
	}
	manager := manager.New()
	if err := manager.UninstallApp(bundleID, udid, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "idev uninstall-app: %v\n", err)
		return manager.ExitCode(err)
	}
	return 0
}

func bundleTarget(args []string, usage string) (string, string, int) {
	bundleID, udid := "", ""
	for index := 0; index < len(args); index++ {
		switch args[index] {
		case "--bundle-id":
			index++
			if index >= len(args) {
				fmt.Fprintln(os.Stderr, usage)
				return "", "", 2
			}
			bundleID = args[index]
		case "--udid":
			index++
			if index >= len(args) {
				fmt.Fprintln(os.Stderr, usage)
				return "", "", 2
			}
			udid = args[index]
		case "-h", "--help":
			fmt.Println(usage)
			return "", "", 0
		default:
			fmt.Fprintf(os.Stderr, "unknown option %q\n", args[index])
			return "", "", 2
		}
	}
	return bundleID, udid, -1
}
