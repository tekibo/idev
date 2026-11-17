package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/tekibo/idev/internal/manager"
)

func setupBase(force, dryRun, asJSON bool) int {
	manager := manager.New()
	plan := manager.BasePlan()
	if dryRun {
		if asJSON {
			data, err := json.MarshalIndent(plan, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
				return 2
			}
			fmt.Println(string(data))
			return 0
		}
		fmt.Printf("Payloads: %s\n", strings.Join(plan.Payloads, ", "))
		fmt.Printf("Components: %s\n", strings.Join(plan.Components, ", "))
		return 0
	}
	ctx := context.Background()
	if _, err := manager.SetupPayloads(ctx, plan.Payloads, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
		return 1
	}
	for _, id := range plan.Components {
		var err error
		switch id {
		case "swift":
			_, err = manager.InstallSwift(ctx, force, os.Stdout)
		case "xtool":
			result, xtoolErr := manager.InstallXtool(ctx, force, os.Stdout, os.Stderr)
			err = xtoolErr
			if err == nil && !asJSON {
				printXtoolShim(result.Shim, result.ShimOnPath)
			}
		case "pymobiledevice3":
			_, err = manager.InstallPymobile(ctx, force, os.Stdout, os.Stderr)
		case "rcodesign":
			_, err = manager.InstallRcodesign(ctx, force)
		case "ipsw":
			_, err = manager.InstallIpsw(ctx, force)
		default:
			err = fmt.Errorf("unknown base component %q", id)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
			return 1
		}
	}
	if !asJSON {
		reportSDK(manager)
	}
	return 0
}
