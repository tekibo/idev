package help

func rnCommand() *Command {
	return &Command{
		Name:        "rn",
		Aliases:     []string{"react-native"},
		Group:       "Build & Run",
		Summary:     "Run the project React Native CLI with the idev adapter.",
		Usage:       "idev rn <command> [args...]",
		Passthrough: true,
		Description: "The manager sets the private environment and runs the project's local React Native CLI.",
		Children: []*Command{
			rnRunIOS(),
			rnBuildIOS(),
			rnStart(),
			rnLogIOS(),
		},
	}
}

func rnRunIOS() *Command {
	return &Command{
		Name:        "run-ios",
		Summary:     "Build, install, and launch the iOS app.",
		Usage:       "idev rn run-ios [options]",
		Passthrough: true,
		Options: []Option{
			def("--mode MODE", "Xcode configuration. Debug or Release.", "Debug"),
			opt("--scheme NAME", "Xcode scheme to build."),
			opt("--device NAME", "Physical device name."),
			opt("--udid UDID", "Device UDID."),
			opt("--binary-path PATH", "Use a pre-built .app."),
			opt("--no-packager", "Do not start Metro."),
			def("--port PORT", "Metro port.", "8081"),
			opt("--buildFolder DIR", "Derived data output directory."),
			opt("--only-pods", "Run CocoaPods only, without building."),
			opt("--force-pods", "Force CocoaPods installation."),
		},
	}
}

func rnBuildIOS() *Command {
	return &Command{
		Name:        "build-ios",
		Summary:     "Build the iOS app without installing or launching it.",
		Usage:       "idev rn build-ios [options]",
		Passthrough: true,
		Options: []Option{
			def("--mode MODE", "Xcode configuration. Debug or Release.", "Debug"),
			opt("--scheme NAME", "Xcode scheme to build."),
			opt("--device NAME", "Physical device name."),
			opt("--buildFolder DIR", "Derived data output directory."),
			opt("--only-pods", "Run CocoaPods only, without building."),
			opt("--force-pods", "Force CocoaPods installation."),
		},
	}
}

func rnStart() *Command {
	return &Command{
		Name:        "start",
		Summary:     "Start the Metro dev server.",
		Usage:       "idev rn start [options]",
		Passthrough: true,
		Options: []Option{
			def("--port PORT", "Metro port.", "8081"),
			opt("--reset-cache", "Reset the Metro cache."),
		},
	}
}

func rnLogIOS() *Command {
	return &Command{
		Name:        "log-ios",
		Summary:     "Stream the iOS device syslog.",
		Usage:       "idev rn log-ios [options]",
		Passthrough: true,
	}
}
