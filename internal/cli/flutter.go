package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

// Flutter runs the pinned Flutter surface.
func Flutter(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: idev flutter <setup|build|ship> [args...] [--dry-run]")
		return 2
	}
	command := args[0]
	switch command {
	case "setup", "build", "ship":
	case "-h", "--help":
		fmt.Println("Usage: idev flutter <setup|build|ship> [args...] [--dry-run]")
		return 0
	default:
		fmt.Fprintf(os.Stderr, "idev flutter: unknown command %q\n", command)
		return 2
	}
	rest := []string{}
	dryRun := false
	for _, arg := range args[1:] {
		if arg == "--dry-run" {
			dryRun = true
			continue
		}
		rest = append(rest, arg)
	}
	manager := manager.New()
	if dryRun {
		fmt.Println(manager.FlutterPreview(command, rest))
		return 0
	}
	script, err := manager.FlutterCommand(command, rest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev flutter: %v\n", err)
		return manager.ExitCode(err)
	}
	err = manager.RunCommand(context.Background(), script.Name, script.Args, script.Env, script.Dir,
		os.Stdin, os.Stdout, os.Stderr, JSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev flutter: %v\n", err)
		return manager.ExitCode(err)
	}
	return 0
}
