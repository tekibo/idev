package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

// InstallApp signs, installs, and launches one app bundle.
func InstallApp(args []string) int {
	var bundle, bundleID, udid, devServerURL string
	for index := 0; index < len(args); index++ {
		switch args[index] {
		case "--bundle":
			index++
			if index < len(args) {
				bundle = args[index]
			}
		case "--bundle-id":
			index++
			if index < len(args) {
				bundleID = args[index]
			}
		case "--udid":
			index++
			if index < len(args) {
				udid = args[index]
			}
		case "--dev-server-url":
			index++
			if index < len(args) {
				devServerURL = args[index]
			}
		case "-h", "--help":
			fmt.Println("Usage: idev install-app --bundle <app> --bundle-id <id> [--udid <udid>] [--dev-server-url <url>]")
			return 0
		default:
			fmt.Fprintf(os.Stderr, "idev install-app: unknown option %q\n", args[index])
			return 2
		}
	}
	if bundle == "" || bundleID == "" {
		fmt.Fprintln(os.Stderr, "idev install-app: --bundle and --bundle-id are required")
		return 2
	}
	manager := manager.New()
	plan, err := manager.InstallPlan(bundle, bundleID, udid, devServerURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev install-app: %v\n", err)
		return 1
	}
	if plan.Terminate.Name != "" {
		_ = manager.RunCommand(context.Background(), plan.Terminate.Name, plan.Terminate.Args,
			plan.Terminate.Env, plan.Terminate.Dir, os.Stdin, os.Stdout, os.Stderr, JSON)
	}
	restore, err := manager.PromptTerminal(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev install-app: %v\n", err)
		return 1
	}
	err = manager.RunCommand(context.Background(), plan.Install.Name, plan.Install.Args,
		plan.Install.Env, plan.Install.Dir, os.Stdin, os.Stdout, os.Stderr, JSON)
	restoreErr := restore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev install-app: %v\n", err)
		return manager.ExitCode(err)
	}
	if restoreErr != nil {
		fmt.Fprintf(os.Stderr, "idev install-app: restore terminal: %v\n", restoreErr)
		return 1
	}
	launch, err := plan.Launch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev install-app: %v\n", err)
		return 1
	}
	err = manager.RunCommand(context.Background(), launch.Name, launch.Args,
		launch.Env, launch.Dir, os.Stdin, os.Stdout, os.Stderr, JSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev install-app: %v\n", err)
		return manager.ExitCode(err)
	}
	return 0
}
