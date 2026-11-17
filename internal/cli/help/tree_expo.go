package help

func expoCommand() *Command {
	return &Command{
		Name:        "expo",
		Group:       "Build & Run",
		Summary:     "Run the project Expo CLI with the idev adapter.",
		Usage:       "idev expo <command> [args...]",
		Passthrough: true,
		Description: "The manager sets the private environment and runs the project's local Expo CLI.",
		Children: []*Command{
			expoRunIOS(),
			expoStart(),
			expoPrebuild(),
			expoExport(),
		},
	}
}

func expoRunIOS() *Command {
	return &Command{
		Name:        "run:ios",
		Summary:     "Build, install, and launch the iOS app.",
		Usage:       "idev expo run:ios [options]",
		Passthrough: true,
		Options: []Option{
			def("--configuration MODE", "Xcode configuration. Debug or Release.", "Debug"),
			opt("--device DEVICE", "Device name, UDID, or 'generic' for build-only."),
			opt("--output PATH", "Directory for the built app binary."),
			opt("--binary PATH", "Install an existing .app or .ipa."),
			opt("--no-bundler", "Do not start Metro."),
			opt("--no-build-cache", "Clear native derived data before building."),
			opt("--no-install", "Skip dependency installation."),
			opt("--scheme NAME", "Xcode scheme to build."),
			def("--port PORT", "Metro port.", "8081"),
		},
		Examples: []Example{
			example("Build only", "idev expo run:ios --device generic --output ./build"),
			example("Release build and run", "idev expo run:ios --configuration Release"),
		},
	}
}

func expoStart() *Command {
	return &Command{
		Name:        "start",
		Summary:     "Start the Metro dev server.",
		Usage:       "idev expo start [options]",
		Passthrough: true,
		Options: []Option{
			opt("--go", "Launch Expo Go."),
			opt("--dev-client", "Launch a custom dev client."),
			opt("--host MODE", "lan, tunnel, or localhost."),
			def("--port PORT", "Dev-server port.", "8081"),
			opt("--scheme NAME", "Custom URI protocol."),
			opt("--clear", "Clear the bundler cache."),
			opt("--tunnel", "Use an Expo tunnel."),
			opt("--offline", "Skip network requests."),
			opt("--no-dev", "Bundle in production mode."),
			opt("--minify", "Minify JavaScript."),
		},
	}
}

func expoPrebuild() *Command {
	return &Command{
		Name:        "prebuild",
		Summary:     "Create native iOS and Android project files.",
		Usage:       "idev expo prebuild [options]",
		Passthrough: true,
		Options: []Option{
			opt("--no-install", "Skip npm and CocoaPods installation."),
			opt("--no-clean", "Update existing native folders instead of recreating them."),
			opt("--npm", "Use npm."),
			opt("--yarn", "Use Yarn."),
			opt("--pnpm", "Use pnpm."),
			opt("--bun", "Use Bun."),
			opt("--platform PLATFORM", "ios, android, or all."),
		},
	}
}

func expoExport() *Command {
	return &Command{
		Name:        "export",
		Summary:     "Export static app files.",
		Usage:       "idev expo export [options]",
		Passthrough: true,
		Options: []Option{
			def("--output-dir DIR", "Export directory.", "dist"),
			opt("--platform PLATFORM", "android, ios, web, or all."),
			opt("--clear", "Clear the bundler cache."),
			opt("--no-minify", "Disable minification."),
		},
	}
}

func iosCommand() *Command {
	return &Command{
		Name:        "ios",
		Group:       "Build & Run",
		Summary:     "Shorthand for 'idev expo run:ios'.",
		Usage:       "idev ios [args...]",
		Passthrough: true,
	}
}

func startCommand() *Command {
	return &Command{
		Name:        "start",
		Group:       "Build & Run",
		Summary:     "Shorthand for 'idev expo start'.",
		Usage:       "idev start [args...]",
		Passthrough: true,
	}
}
