package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/tekibo/idev/internal/manager"
)

func setupIpsw(force, dryRun bool) int {
	manager := manager.New()
	spec, err := manager.IpswPlan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
		return 1
	}
	if dryRun {
		fmt.Printf("ipsw: %s\n", spec.Version)
		fmt.Printf("URL:  %s\n", spec.Component.URL)
		fmt.Printf("SHA:  %s\n", spec.Component.SHA256)
		fmt.Printf("Bin:  %s\n", spec.Bin)
		return 0
	}
	result, err := manager.InstallIpsw(context.Background(), force)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
		return 1
	}
	fmt.Printf("ipsw: %s\n", result.Version)
	fmt.Printf("Bin:  %s\n", result.Bin)
	if result.Already {
		fmt.Println("State: already installed")
	}
	return 0
}
