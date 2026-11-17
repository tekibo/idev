package help

func installAppCommand() *Command {
	return &Command{
		Name:    "install-app",
		Group:   "Devices",
		Summary: "Sign, install, and launch one app bundle.",
		Usage:   "idev install-app --bundle <app> --bundle-id <id> [--udid <udid>] [--dev-server-url <url>]",
		Options: []Option{
			req("--bundle PATH", "Path to the .app or .ipa bundle."),
			req("--bundle-id ID", "Installed bundle identifier."),
			opt("--udid UDID", "Device UDID. Default: the only connected device."),
			opt("--dev-server-url URL", "Initial dev-server URL for a dev client."),
		},
	}
}

func launchAppCommand() *Command {
	return &Command{
		Name:    "launch-app",
		Group:   "Devices",
		Summary: "Launch an installed app on one device.",
		Usage:   "idev launch-app --bundle-id <id> [--udid <udid>]",
		Options: []Option{
			req("--bundle-id ID", "Installed bundle identifier."),
			opt("--udid UDID", "Device UDID. Default: the only connected device."),
		},
	}
}

func uninstallAppCommand() *Command {
	return &Command{
		Name:    "uninstall-app",
		Group:   "Devices",
		Summary: "Remove an installed app from one device.",
		Usage:   "idev uninstall-app --bundle-id <id> [--udid <udid>]",
		Options: []Option{
			req("--bundle-id ID", "Installed bundle identifier."),
			opt("--udid UDID", "Device UDID. Default: the only connected device."),
		},
	}
}
