package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	helptree "github.com/tekibo/idev/internal/cli/help"
	"github.com/tekibo/idev/internal/manager"
)

// Auth runs xtool auth through the shared command session.
func Auth(args []string) int {
	if len(args) == 2 && args[0] == "status" && wantsHelp(args[1:]) {
		return helptree.Show([]string{"auth", "status"}, JSON)
	}
	manager := manager.New()
	if len(args) == 2 && args[0] == "status" && args[1] == "--json" {
		status, err := manager.AuthStatus()
		if err != nil {
			fmt.Fprintf(os.Stderr, "idev auth: %v\n", err)
			return manager.ExitCode(err)
		}
		data, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "idev auth: %v\n", err)
			return 2
		}
		fmt.Println(string(data))
		return 0
	}
	command, err := manager.AuthArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev auth: %v\n", err)
		return 2
	}
	if !manager.Xtool.Available() {
		fmt.Fprintln(os.Stderr, "idev auth: xtool not found; run idev setup")
		return 127
	}
	if err := manager.RunCommand(context.Background(), manager.Xtool.Path, command, manager.Xtool.Env, "",
		os.Stdin, os.Stdout, os.Stderr, JSON); err != nil {
		fmt.Fprintf(os.Stderr, "idev auth: %v\n", err)
		return manager.ExitCode(err)
	}
	return 0
}
