package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

func setupPayload(force, installSelf bool) int {
	manager := manager.New()
	result, err := manager.Setup(context.Background(), force, installSelf, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
		return 1
	}
	fmt.Printf("idev home: %s\n", result.Home)
	fmt.Printf("React Native: %s\n", result.ReactNative)
	fmt.Printf("Apple tools: %s\n", result.AppleTools)
	fmt.Printf("updated:   %d file(s)\n", result.Written)
	if installSelf {
		state := "already installed"
		if result.ManagerInstalled {
			state = "installed"
		}
		fmt.Printf("manager:   %s (%s)\n", result.ManagerPath, state)
	}
	fmt.Print(`
Next:
  1. Run the first-run wizard: idev
  2. Run: idev auth
  3. Connect and trust the device, then: idev doctor
  4. In an Expo app: idev expo run:ios

idev does not embed the Darwin SDK, Swift toolchain, Xcode.xip, and private components.
`)
	return 0
}
