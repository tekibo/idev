package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

func setupSwift(force, dryRun bool) int {
	manager := manager.New()
	spec, err := manager.SwiftPlan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
		return 1
	}
	if dryRun {
		fmt.Printf("Swift: %s\n", spec.Version)
		fmt.Printf("URL:   %s\n", spec.URL)
		fmt.Printf("SHA:   %s\n", spec.SHA256)
		fmt.Printf("Root:  %s\n", spec.Root)
		fmt.Printf("Cache: %s\n", spec.Cache)
		return 0
	}
	result, err := manager.InstallSwift(context.Background(), force, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
		return 1
	}
	fmt.Printf("Swift: %s\n", result.Version)
	fmt.Printf("Root:  %s\n", result.Root)
	fmt.Printf("Cache: %s\n", result.Cache)
	if result.Already {
		fmt.Println("State: already installed")
	}
	return 0
}
