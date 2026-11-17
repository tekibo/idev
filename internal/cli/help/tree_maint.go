package help

func maintenanceCommands() []*Command {
	return []*Command{
		configCommand(),
		updateCommand(),
		uninstallCommand(),
		versionCommand(),
	}
}

func updateCommand() *Command {
	return &Command{
		Name:    "update",
		Group:   "Maintenance",
		Summary: "Check, apply, and clean manager and component updates.",
		Usage:   "idev update <check|apply|clean> [--all] [component]",
		Children: []*Command{
			{
				Name:    "check",
				Summary: "Show manager and component updates.",
				Usage:   "idev update check [--json]",
				Options: []Option{opt("--json", "Emit the update list as JSON.")},
			},
			{
				Name:    "apply",
				Summary: "Install one component or every outdated component.",
				Usage:   "idev update apply [--all] [component]",
				Options: []Option{
					opt("--all", "Install every outdated component and reconcile the manager."),
				},
			},
			{
				Name:    "clean",
				Summary: "Remove update downloads, caches, and stale files.",
				Usage:   "idev update clean [--dry-run] [--yes] [--json]",
				Options: []Option{
					opt("--dry-run", "Show the cleanup plan without removing files."),
					opt("--yes", "Skip the confirmation prompt."),
					opt("--json", "Emit the plan or result as JSON."),
				},
				JSON: []string{"targets", "bytes", "removed"},
			},
		},
	}
}

func uninstallCommand() *Command {
	return &Command{
		Name:    "uninstall",
		Group:   "Maintenance",
		Summary: "Remove idev-managed paths.",
		Usage:   "idev uninstall [--dry-run|--latest] [--json] [--keep <id>]... [--confirm \"remove idev\"]",
		Options: []Option{
			opt("--dry-run", "Show what would be removed."),
			opt("--latest", "Use the latest saved uninstall report."),
			opt("--json", "Emit the report as JSON."),
			opt("--keep ID", "Keep one component. Repeat for more."),
			opt("--confirm PHRASE", "Confirm removal with the exact phrase."),
		},
	}
}

func versionCommand() *Command {
	return &Command{
		Name:    "version",
		Aliases: []string{"-v", "--version"},
		Group:   "Maintenance",
		Summary: "Show the manager version.",
		Usage:   "idev version",
	}
}
