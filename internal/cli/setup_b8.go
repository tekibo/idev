package cli

func setupB8(flags []string, force, dryRun bool) int {
	for _, flag := range flags {
		switch flag {
		case "--rcodesign":
			if code := setupRcodesign(force, dryRun); code != 0 {
				return code
			}
		case "--ipsw":
			if code := setupIpsw(force, dryRun); code != 0 {
				return code
			}
		case "--swift-build":
			if code := setupSwiftBuild(force, dryRun); code != 0 {
				return code
			}
		}
	}
	return 0
}
