package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

func setupPymobile(force, dryRun bool) int {
	manager := manager.New()
	spec := manager.PymobilePlan()
	if dryRun {
		fmt.Printf("pymobiledevice3: %s\n", spec.Version)
		fmt.Printf("Root:  %s\n", spec.Root)
		fmt.Printf("Bin:   %s\n", spec.Bin)
		fmt.Printf("Link:  %s\n", spec.Link)
		fmt.Printf("Cache: %s\n", spec.Cache)
		return 0
	}
	result, err := manager.InstallPymobile(context.Background(), force, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
		return 1
	}
	fmt.Printf("pymobiledevice3: %s\n", result.Version)
	fmt.Printf("Root:  %s\n", result.Root)
	fmt.Printf("Link:  %s\n", result.Link)
	if result.Already {
		fmt.Println("State: already installed")
	}
	return 0
}
