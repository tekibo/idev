package help

func setupCommands() []*Command {
	return []*Command{
		setupCommand(),
		nextCommand(),
		agentsCommand(),
		doctorCommand(),
		systemCommand(),
		devicesCommand(),
		authCommand(),
		xcodeCommand(),
	}
}

func setupCommand() *Command {
	return &Command{
		Name:    "setup",
		Group:   "Setup",
		Summary: "Install the selected setup groups and components.",
		Usage:   "idev setup [--flow NAME] [--no-expo] [options]",
		Options: []Option{
			opt("--flow NAME", "Install one named setup flow."),
			opt("--no-expo", "Disable the Expo arm for the setup flow."),
			opt("--force", "Reinstall the selected components."),
			opt("--install-self", "Install the manager itself."),
			opt("--base", "Install the base system packages."),
			opt("--swift", "Install private Swift."),
			opt("--xtool", "Install xtool."),
			opt("--apple-tools", "Install the Darwin tools and OAM."),
			opt("--pymobile", "Create the pymobiledevice3 venv."),
			opt("--cocoapods", "Create the private CocoaPods tree."),
			opt("--payloads", "Install the release payloads."),
			opt("--rcodesign", "Install rcodesign."),
			opt("--ipsw", "Install ipsw."),
			opt("--swift-build", "Install Swift Build."),
			opt("--dry-run", "Print the plan without installing."),
			opt("--json", "Emit the base plan as JSON. Requires --base --dry-run."),
		},
		Examples: []Example{
			example("Preview Swift", "idev setup --swift --dry-run"),
			example("Install the base system", "idev setup --base"),
		},
	}
}

func doctorCommand() *Command {
	return &Command{
		Name:    "doctor",
		Group:   "Setup",
		Summary: "Check tools, SDK, devices, and the current project.",
		Usage:   "idev doctor [--json] [--private] [--project DIR]",
		Options: []Option{
			opt("--json", "Emit the report as JSON."),
			opt("--private", "Check the private toolchain and SDK."),
			opt("--project DIR", "Project directory to include project checks."),
		},
	}
}

func systemCommand() *Command {
	return &Command{
		Name:    "system",
		Group:   "System",
		Summary: "Plan or install the missing system packages.",
		Usage:   "idev system <plan|apply> [--all] [--json]",
		Options: []Option{
			opt("--all", "Include optional packages."),
			opt("--json", "Emit JSON. plan only."),
		},
		Children: []*Command{
			{Name: "plan", Summary: "Show the exact package command.", Usage: "idev system plan [--all] [--json]", Options: []Option{
				opt("--all", "Include optional packages."),
				opt("--json", "Emit the plan as JSON."),
			}},
			{Name: "apply", Summary: "Install the exact package command.", Usage: "idev system apply [--all]", Options: []Option{
				opt("--all", "Include optional packages."),
			}},
		},
	}
}

func devicesCommand() *Command {
	return &Command{
		Name:    "devices",
		Group:   "Devices",
		Summary: "List connected iOS devices.",
		Usage:   "idev devices [--json]",
		Options: []Option{opt("--json", "Emit the device list as JSON.")},
	}
}

func authCommand() *Command {
	return &Command{
		Name:        "auth",
		Group:       "Setup",
		Summary:     "Run xtool auth and read the account status.",
		Usage:       "idev auth <args...>",
		Passthrough: true,
		Description: "Other auth arguments are passed to xtool. Use 'idev auth status --json' for the parsed account status.",
		Children: []*Command{
			{Name: "status", Summary: "Show the account status as JSON.", Usage: "idev auth status --json", Options: []Option{
				opt("--json", "Emit the parsed account status."),
			}},
		},
	}
}
