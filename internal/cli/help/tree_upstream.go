package help

func upstreamCommands() []*Command {
	return []*Command{
		kmpCommand(),
		hostManifestsCommand(),
		simViewCommand(),
		machoLoadCommand(),
		compileAppIconCommand(),
	}
}

func upstreamExtraCommands() []*Command {
	return []*Command{
		xcodeProjectCommand(),
		ascCommand(),
		provisionDevCommand(),
		signDevCommand(),
		sdkFreeCommand(),
	}
}

func kmpCommand() *Command {
	return &Command{
		Name:        "kmp",
		Group:       "Tools",
		Summary:     "Link a Kotlin Multiplatform iOS framework on Linux.",
		Usage:       "idev kmp [project dir] [--module :path] [--target ios_arm64|ios_simulator_arm64] [--out dir] [--dynamic] [--xcframework]",
		Passthrough: true,
		Options: []Option{
			opt("--module :path", "One Gradle module. Default: every framework project."),
			def("--target TARGET", "ios_arm64 or ios_simulator_arm64.", "ios_arm64"),
			opt("--out DIR", "Output directory. Default: <module>/build/omarchy-kmp/<target>."),
			opt("--dynamic", "Link a dynamic framework."),
			opt("--xcframework", "Wrap the framework as an xcframework."),
		},
	}
}

func hostManifestsCommand() *Command {
	return &Command{
		Name:        "host-manifests",
		Group:       "Tools",
		Summary:     "Answer SwiftPM host manifest checks as on a Mac.",
		Usage:       "idev host-manifests <package dir>",
		Passthrough: true,
		Options: []Option{
			opt("--self-test", "Run the built-in manifest checks."),
		},
	}
}

func simViewCommand() *Command {
	children := []*Command{}
	for _, name := range []string{"build", "test", "test-m3", "test-m4", "test-m5", "test-m6", "test-m7"} {
		children = append(children, &Command{
			Name: name, Summary: "Run simview " + name + ".", Usage: "idev simview " + name + " [args...]", Passthrough: true,
		})
	}
	return &Command{
		Name:        "simview",
		Group:       "Tools",
		Summary:     "Run the simview scripts.",
		Usage:       "idev simview <build|test|test-m3|test-m4|test-m5|test-m6|test-m7> [args...]",
		Passthrough: true,
		Children:    children,
	}
}

func machoLoadCommand() *Command {
	return &Command{
		Name:        "machoload",
		Group:       "Tools",
		Summary:     "Build the Mach-O loader experiment.",
		Usage:       "idev machoload [workdir]",
		Passthrough: true,
	}
}

func compileAppIconCommand() *Command {
	return &Command{
		Name:        "compile-app-icon",
		Group:       "Tools",
		Summary:     "Compile the app icon catalog into an app.",
		Usage:       "idev compile-app-icon <app>",
		Passthrough: true,
	}
}
