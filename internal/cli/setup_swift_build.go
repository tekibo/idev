package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

func setupSwiftBuild(force, dryRun bool) int {
	manager := manager.New()
	spec, err := manager.SwiftBuildPlan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
		return 1
	}
	if dryRun {
		fmt.Printf("Swift Build: %s\n", spec.Version)
		fmt.Printf("Source:     %s\n", spec.SourceRoot)
		fmt.Printf("Swift:      %s\n", spec.SwiftBin)
		fmt.Printf("Engine:     %s\n", spec.EngineSource)
		fmt.Printf("Service:    %s\n", spec.Bin)
		fmt.Printf("Engine bin: %s\n", spec.EngineBin)
		fmt.Printf("Bundles:    %s\n", spec.Bundles)
		return 0
	}
	result, err := manager.InstallSwiftBuild(context.Background(), force, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
		return 1
	}
	fmt.Printf("Swift Build: %s\n", result.Version)
	fmt.Printf("Service:     %s\n", result.Bin)
	fmt.Printf("Engine bin:  %s\n", result.EngineBin)
	if result.Already {
		fmt.Println("State: already installed")
	}
	return 0
}
