package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

func setupCocoaPods(force, dryRun bool) int {
	manager := manager.New()
	spec := manager.CocoaPodsPlan()
	if dryRun {
		fmt.Printf("CocoaPods: %s\n", spec.Version)
		fmt.Printf("Root:  %s\n", spec.Root)
		fmt.Printf("Bin:   %s\n", spec.Bin)
		fmt.Printf("Link:  %s\n", spec.Link)
		fmt.Printf("Cache: %s\n", spec.Cache)
		return 0
	}
	result, err := manager.InstallCocoaPods(context.Background(), force, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
		return 1
	}
	fmt.Printf("CocoaPods: %s\n", result.Version)
	fmt.Printf("Root:  %s\n", result.Root)
	fmt.Printf("Link:  %s\n", result.Link)
	if result.Already {
		fmt.Println("State: already installed")
	}
	return 0
}
