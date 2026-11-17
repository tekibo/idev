package cli

import "github.com/tekibo/idev/internal/manager"

func setupTail(force, dryRun, installSelf, swiftOnly, xtoolOnly, expoSet bool, selection manager.FlowSelection) int {
	if swiftOnly && xtoolOnly {
		return setupBoth(force, dryRun)
	}
	if swiftOnly {
		return setupSwift(force, dryRun)
	}
	if xtoolOnly {
		return setupXtool(force, dryRun)
	}
	if installSelf {
		return setupPayload(force, true)
	}
	return setupSaved(force, expoSet, selection, dryRun)
}
