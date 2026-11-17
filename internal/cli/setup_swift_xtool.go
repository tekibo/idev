package cli

func setupBoth(force, dryRun bool) int {
	if dryRun {
		if code := setupSwift(force, true); code != 0 {
			return code
		}
		return setupXtool(force, true)
	}
	if code := setupSwift(force, false); code != 0 {
		return code
	}
	return setupXtool(force, false)
}
