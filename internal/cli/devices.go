package cli

import (
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

// Devices prints connected devices.
func Devices(args []string) int {
	manager := manager.New()
	if len(args) > 0 && args[0] == "--json" {
		data, err := manager.DevicesJSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "idev devices: %v\n", err)
			return manager.ExitCode(err)
		}
		fmt.Println(string(data))
		return 0
	}
	if err := manager.Devices(os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "idev devices: %v\n", err)
		return manager.ExitCode(err)
	}
	return 0
}
