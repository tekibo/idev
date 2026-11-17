package cli

import (
	"fmt"
	"os"
)

// setupBody installs the selected setup groups and components.
func setupBody(args []string) int {
	selection, args, expoSet, err := takeFlowSelection(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "idev setup: %v\n", err)
		return 2
	}
	force, installSelf := false, false
	swiftOnly, xtoolOnly := false, false
	appleToolsOnly, pymobileOnly := false, false
	cocoapodsOnly, payloadsOnly := false, false
	baseOnly, asJSON := false, false
	var b8 []string
	dryRun := false
	for _, arg := range args {
		switch arg {
		case "-f", "--force":
			force = true
		case "--install-self":
			installSelf = true
		case "--swift":
			swiftOnly = true
		case "--xtool":
			xtoolOnly = true
		case "--apple-tools":
			appleToolsOnly = true
		case "--pymobile":
			pymobileOnly = true
		case "--cocoapods":
			cocoapodsOnly = true
		case "--payloads":
			payloadsOnly = true
		case "--base":
			baseOnly = true
		case "--rcodesign", "--ipsw", "--swift-build":
			b8 = append(b8, arg)
		case "--dry-run":
			dryRun = true
		case "--json":
			asJSON = true
		case "-h", "--help":
			fmt.Println("Usage: idev setup [--flow <name>] [--no-expo] [--force] [--install-self] [--base] [--swift] [--xtool] [--apple-tools] [--pymobile] [--cocoapods] [--payloads] [--rcodesign] [--ipsw] [--swift-build] [--dry-run]")
			return 0
		default:
			fmt.Fprintf(os.Stderr, "idev setup: unknown option %q\n", arg)
			return 2
		}
	}
	if asJSON && !(baseOnly && dryRun) {
		fmt.Fprintln(os.Stderr, "idev setup: --json requires --base --dry-run")
		return 2
	}
	if !setupMigrate(dryRun) {
		return 1
	}
	if expoSet && (installSelf || swiftOnly || xtoolOnly || appleToolsOnly || pymobileOnly || cocoapodsOnly || payloadsOnly || len(b8) > 0) {
		fmt.Fprintln(os.Stderr, "idev setup: --no-expo applies to the setup flow")
		return 2
	}
	if baseOnly {
		if installSelf || swiftOnly || xtoolOnly || appleToolsOnly || pymobileOnly || cocoapodsOnly || payloadsOnly || len(b8) > 0 || len(selection.Flows) > 0 || expoSet {
			fmt.Fprintln(os.Stderr, "idev setup: --base runs alone")
			return 2
		}
		return setupBase(force, dryRun, asJSON)
	}
	if len(selection.Flows) > 0 {
		if installSelf || swiftOnly || xtoolOnly || appleToolsOnly || pymobileOnly || cocoapodsOnly || payloadsOnly || len(b8) > 0 {
			fmt.Fprintln(os.Stderr, "idev setup: --flow runs as its own setup step")
			return 2
		}
		return setupFlow(selection, force, dryRun)
	}
	if appleToolsOnly {
		if swiftOnly || xtoolOnly {
			fmt.Fprintln(os.Stderr, "idev setup: --apple-tools runs after --swift and the xtool SDK install")
			return 2
		}
		return setupAppleTools(force, dryRun)
	}
	if pymobileOnly || cocoapodsOnly {
		if swiftOnly || xtoolOnly || appleToolsOnly {
			fmt.Fprintln(os.Stderr, "idev setup: --pymobile/--cocoapods run as their own setup step")
			return 2
		}
		if pymobileOnly && cocoapodsOnly {
			return setupB7(force, dryRun)
		}
		if pymobileOnly {
			return setupPymobile(force, dryRun)
		}
		return setupCocoaPods(force, dryRun)
	}
	if len(b8) > 0 {
		if swiftOnly || xtoolOnly || appleToolsOnly || pymobileOnly || cocoapodsOnly {
			fmt.Fprintln(os.Stderr, "idev setup: --rcodesign/--ipsw/--swift-build run as their own setup step")
			return 2
		}
		return setupB8(b8, force, dryRun)
	}
	if payloadsOnly {
		if dryRun {
			fmt.Fprintln(os.Stderr, "idev setup: --payloads does not support --dry-run")
			return 2
		}
		return setupPayload(force, installSelf)
	}
	return setupTail(force, dryRun, installSelf, swiftOnly, xtoolOnly, expoSet, selection)
}
