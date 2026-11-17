package cli

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/tekibo/idev/internal/manager"
)

// ReactNative runs the project React Native CLI.
func ReactNative(args []string) int {
	manager := manager.New()
	command, err := manager.ReactNativeCommand(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev rn: %v\n", err)
		return 2
	}
	if !JSON {
		if err := os.Chdir(command.Dir); err != nil {
			fmt.Fprintf(os.Stderr, "idev rn: %v\n", err)
			return 2
		}
		argv := append([]string{command.Name}, command.Args...)
		if err := syscall.Exec(command.Name, argv, command.Env); err != nil {
			fmt.Fprintf(os.Stderr, "idev rn: %v\n", err)
			return 1
		}
	}
	err = manager.RunCommand(context.Background(), command.Name, command.Args, command.Env, command.Dir,
		os.Stdin, os.Stdout, os.Stderr, JSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev rn: %v\n", err)
		return manager.ExitCode(err)
	}
	return 0
}
