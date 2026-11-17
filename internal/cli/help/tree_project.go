package help

func projectCommands() []*Command {
	return []*Command{
		newCommand(),
		projectCommand(),
		installAppCommand(),
		launchAppCommand(),
		uninstallAppCommand(),
		ipadScreenCommand(),
	}
}

func projectCommand() *Command {
	return &Command{
		Name:    "project",
		Group:   "Project",
		Summary: "Inspect and modify the project directory.",
		Usage:   "idev project <command> [options]",
		Children: []*Command{
			{
				Name:    "type",
				Summary: "Show the detected project kind.",
				Usage:   "idev project type",
			},
			{
				Name:    "info",
				Summary: "Show the detected project fields.",
				Usage:   "idev project info [--json] [--project DIR] [--app NAME|DIR]",
				Options: []Option{
					opt("--json", "Emit the project descriptor as JSON."),
					opt("--project DIR", "Project directory. Default: IDEV_PROJECT or the current directory."),
					opt("--app NAME|DIR", "Select one app in a workspace."),
				},
				JSON: []string{"kind", "dir", "workspace_dir", "name", "label", "bundle_id", "workflow", "dev_client", "native_dir", "project_path", "scheme", "package_manager", "warnings"},
			},
			{
				Name:    "apps",
				Summary: "List the apps detected in a workspace.",
				Usage:   "idev project apps [--json] [--project DIR]",
				Options: []Option{
					opt("--json", "Emit the workspace and app list as JSON."),
					opt("--project DIR", "Project or workspace directory. Default: IDEV_PROJECT or the current directory."),
				},
				JSON: []string{"dir", "name", "apps", "warnings"},
			},
			{
				Name:    "set-bundle-id",
				Summary: "Set the bundle identifier in xtool.yml.",
				Usage:   "idev project set-bundle-id [--project DIR] --value ID",
				Options: []Option{
					opt("--project DIR", "Project directory. Default: IDEV_PROJECT or the current directory."),
					req("--value ID", "New reverse-DNS bundle identifier."),
				},
			},
			{
				Name:    "clean",
				Summary: "Remove build output for the current project.",
				Usage:   "idev project clean",
			},
		},
	}
}
