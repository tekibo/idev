package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

func setupAppleTools(force, dryRun bool) int {
	manager := manager.New()
	spec, err := manager.AppleToolsPlan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
		return 1
	}
	if spec.SDKRoot == "" {
		fmt.Fprintln(os.Stderr, "idev setup: Darwin SDK not found; run idev sdk install <Xcode.xip>")
		return 1
	}
	if dryRun {
		fmt.Printf("Payload: %s\n", spec.Payload)
		fmt.Printf("SDK:     %s\n", spec.SDKRoot)
		fmt.Printf("Tools:   %s\n", spec.ToolsBin)
		fmt.Printf("OAM:     %s\n", spec.OAMServer)
		return 0
	}
	result, err := manager.InstallAppleTools(context.Background(), force, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
		return 1
	}
	fmt.Printf("Payload:   %s\n", result.Payload)
	fmt.Printf("SDK:       %s\n", result.SDKRoot)
	fmt.Printf("Tools:     %s\n", result.ToolsBin)
	fmt.Printf("Mac tools: %s\n", result.MacToolsBin)
	fmt.Printf("OAM:       %s\n", result.OAMServer)
	fmt.Printf("Plugins:   %s\n", result.Plugins)
	return 0
}
