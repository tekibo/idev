package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tekibo/idev/internal/manager"
)

func setupXtool(force, dryRun bool) int {
	manager := manager.New()
	spec, err := manager.XtoolPlan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
		return 1
	}
	if dryRun {
		fmt.Printf("xtool:  %s\n", spec.Pin)
		fmt.Printf("Repo:   %s\n", spec.Repo)
		fmt.Printf("Archive: %s\n", spec.Archive)
		fmt.Printf("SHA:    %s\n", spec.SHA256)
		fmt.Printf("Source: %s\n", spec.Source)
		fmt.Printf("Root:   %s\n", spec.Root)
		fmt.Printf("Link:   %s\n", spec.Link)
		return 0
	}
	result, err := manager.InstallXtool(context.Background(), force, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
		return 1
	}
	fmt.Printf("xtool: %s\n", result.Pin)
	fmt.Printf("Root:  %s\n", result.Root)
	fmt.Printf("Link:  %s\n", result.Link)
	if result.Already {
		fmt.Println("State: already installed")
	}
	printXtoolShim(result.Shim, result.ShimOnPath)
	reportSDK(manager)
	return 0
}

func printXtoolShim(shim string, onPath bool) {
	if shim == "" {
		return
	}
	fmt.Printf("Shim:  %s\n", shim)
	if !onPath {
		fmt.Printf("Path:  add %s to PATH for raw xtool\n", filepath.Dir(shim))
	}
}

func reportSDK(manager *manager.App) {
	status, err := manager.SDKStatus()
	if err == nil && status.Installed {
		fmt.Printf("SDK:   %s\n", status.Path)
		return
	}
	fmt.Println("SDK:   not installed; run idev sdk install <Xcode.xip>")
}
