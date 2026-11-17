package help

func nativeCommands() []*Command {
	return []*Command{
		swiftCommand(),
		podCommand(),
		lspCommand(),
		flutterCommand(),
		shipCommand(),
		deviceRunCommand(),
	}
}

func swiftCommand() *Command {
	return &Command{
		Name:        "swift",
		Group:       "Build & Run",
		Summary:     "Run the xtool SwiftPM flow.",
		Usage:       "idev swift <command> [args...]",
		Passthrough: true,
		Description: "The manager runs xtool with the private Swift, SDK, and Swift Build environment.",
		Children: []*Command{
			{Name: "build", Summary: "Build the SwiftPM app.", Passthrough: true, Usage: "idev swift build [args...]"},
			{Name: "run", Summary: "Build, install, and launch the SwiftPM app.", Passthrough: true, Usage: "idev swift run [args...]"},
			{Name: "new", Summary: "Create a new xtool SwiftPM project.", Passthrough: true, Usage: "idev swift new [args...]"},
			{Name: "generate-xcode-project", Summary: "Generate the Xcode project.", Passthrough: true, Usage: "idev swift generate-xcode-project [args...]"},
			{Name: "build-server", Summary: "Run the Swift Build server.", Passthrough: true, Usage: "idev swift build-server [args...]"},
		},
	}
}

func podCommand() *Command {
	return &Command{
		Name:        "pod",
		Group:       "Build & Run",
		Summary:     "Run the private CocoaPods wrapper.",
		Usage:       "idev pod [args...]",
		Passthrough: true,
	}
}

func lspCommand() *Command {
	return &Command{
		Name:        "lsp",
		Group:       "Tools",
		Summary:     "Run the private SourceKit-LSP over stdio.",
		Usage:       "idev lsp [args...]",
		Passthrough: true,
	}
}

func flutterCommand() *Command {
	return &Command{
		Name:    "flutter",
		Group:   "Build & Run",
		Summary: "Run the pinned Flutter setup, build, and ship scripts.",
		Usage:   "idev flutter <setup|build|ship> [args...] [--dry-run]",
		Options: []Option{
			opt("--dry-run", "Print the exact command without running it."),
		},
		Children: []*Command{
			{Name: "setup", Summary: "Install the pinned Flutter toolchain.", Usage: "idev flutter setup [args...] [--dry-run]", Options: []Option{
				opt("--dry-run", "Print the exact command without running it."),
			}},
			{Name: "build", Summary: "Build the Flutter iOS app.", Usage: "idev flutter build [args...] [--dry-run]", Options: []Option{
				opt("--dry-run", "Print the exact command without running it."),
			}},
			{Name: "ship", Summary: "Build and ship the Flutter app.", Usage: "idev flutter ship [args...] [--dry-run]", Options: []Option{
				opt("--dry-run", "Print the exact command without running it."),
			}},
		},
	}
}

func shipCommand() *Command {
	return &Command{
		Name:    "ship",
		Group:   "Build & Run",
		Summary: "Build and ship an xtool SwiftUI app.",
		Usage:   "idev ship [--upload] [--device [TEAMID]] [--xcode-version V] [--xcode-build B]",
		Options: []Option{
			opt("--upload", "Upload the build."),
			opt("--device [TEAMID]", "Build for a device, optionally with a team ID."),
			opt("--xcode-version V", "Xcode version for the build."),
			opt("--xcode-build B", "Xcode build number for the build."),
		},
	}
}

func deviceRunCommand() *Command {
	return &Command{
		Name:    "device-run",
		Group:   "Devices",
		Summary: "Run the apple-dev device helper.",
		Usage:   "idev device-run [--lldb|--attach|--console [SECONDS]] [args...]",
		Options: []Option{
			opt("--lldb", "Start an LLDB session."),
			opt("--attach", "Attach to a running process."),
			opt("--console [SECONDS]", "Stream the device console."),
		},
	}
}
