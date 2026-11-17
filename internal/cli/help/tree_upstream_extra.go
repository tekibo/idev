package help

func xcodeProjectCommand() *Command {
	return &Command{
		Name:        "xcodeproj2xtool",
		Group:       "Tools",
		Summary:     "Convert an Xcode project into an xtool package.",
		Usage:       "idev xcodeproj2xtool [project] [--target NAME] [--out DIR] [--configuration NAME] [--bundle-id ID] [--self-test]",
		Passthrough: true,
		Options: []Option{
			opt("--target NAME", "Application target name when the project has several."),
			opt("--out DIR", "Output directory. Default: <project root>/omarchy-xtool."),
			def("--configuration NAME", "Build configuration used by the adapter.", "Release"),
			opt("--bundle-id ID", "Set the app bundle ID and rebase extension IDs."),
			opt("--self-test", "Run the built-in round-trip check."),
		},
	}
}

func ascCommand() *Command {
	return &Command{
		Name:        "asc",
		Group:       "Tools",
		Summary:     "Run the App Store Connect helper.",
		Usage:       "idev asc <stamp|frameworks|identity|given-identity|test-identity|validate|upload|publish> [args...]",
		Passthrough: true,
	}
}

func provisionDevCommand() *Command {
	return &Command{
		Name:        "provision-dev",
		Group:       "Tools",
		Summary:     "Create a development provisioning profile.",
		Usage:       "idev provision-dev --bundle-id ID [options]",
		Passthrough: true,
		Options: []Option{
			req("--bundle-id ID", "Bundle identifier. Repeat for more than one."),
			opt("--udid UDID", "Device UDID. Repeat for more than one."),
			def("--name NAME", "Profile name.", "iPhone"),
			opt("--out DIR", "Output directory."),
			opt("--expect-group GROUP", "Require one App Group."),
		},
	}
}

func signDevCommand() *Command {
	return &Command{
		Name:        "sign-dev",
		Group:       "Tools",
		Summary:     "Sign an app with a development identity.",
		Usage:       "idev sign-dev <app> [team] [dev dir]",
		Passthrough: true,
	}
}

func sdkFreeCommand() *Command {
	children := []*Command{
		{Name: "setup", Summary: "Build the no-Xcode sysroot.", Usage: "idev sdk-free setup [--ipsw] [--repair]", Passthrough: true, Options: []Option{
			opt("--ipsw", "Download the cache from an iOS update instead of a device."),
			opt("--repair", "Repair an existing sysroot."),
		}},
		{Name: "flutter-build", Summary: "Build a Flutter app without Xcode.", Usage: "idev sdk-free flutter-build [--debug] <app>", Passthrough: true, Options: []Option{
			opt("--debug", "Build the debug configuration."),
		}},
		{Name: "ship", Summary: "Ship a no-Xcode app.", Usage: "idev sdk-free ship [--upload] <app>", Passthrough: true, Options: []Option{
			opt("--upload", "Upload to App Store Connect."),
		}},
		{Name: "cc", Summary: "Compile and link with the sdk-free sysroot.", Usage: "idev sdk-free cc [args...]", Passthrough: true},
		{Name: "swiftc", Summary: "Compile Swift with the sdk-free sysroot.", Usage: "idev sdk-free swiftc [args...]", Passthrough: true},
		{Name: "swift-build-stdlib", Summary: "Build the sdk-free Swift stdlib.", Usage: "idev sdk-free swift-build-stdlib [args...]", Passthrough: true},
		{Name: "swift-overlays", Summary: "Build the sdk-free Swift overlays.", Usage: "idev sdk-free swift-overlays [args...]", Passthrough: true},
		{Name: "swiftui-build", Summary: "Build the sdk-free SwiftUI module.", Usage: "idev sdk-free swiftui-build [args...]", Passthrough: true},
	}
	return &Command{
		Name:        "sdk-free",
		Group:       "Tools",
		Summary:     "Run the experimental no-Xcode scripts.",
		Usage:       "idev sdk-free <setup|flutter-build|ship|cc|swiftc|swift-build-stdlib|swift-overlays|swiftui-build> [args...]",
		Passthrough: true,
		Children:    children,
	}
}
