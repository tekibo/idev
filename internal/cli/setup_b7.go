package cli

func setupB7(force, dryRun bool) int {
	if code := setupPymobile(force, dryRun); code != 0 {
		return code
	}
	return setupCocoaPods(force, dryRun)
}
