package help

func agentsCommand() *Command {
	return &Command{
		Name:        "agents",
		Group:       "Setup",
		Summary:     "Install the idev skill for coding agents.",
		Usage:       "idev agents <install|uninstall|status> [--json]",
		Description: "Writes the embedded idev skill under IDEV_HOME, links it into the agent skill directories, and adds the idev section to the global instructions.",
		Children: []*Command{
			{
				Name:    "install",
				Summary: "Write the skill, link it, and add the instruction section.",
				Usage:   "idev agents install [--json]",
				Options: []Option{opt("--json", "Emit the installed paths as JSON.")},
			},
			{
				Name:    "uninstall",
				Summary: "Remove the links, the managed skill, and the instruction section.",
				Usage:   "idev agents uninstall [--json]",
				Options: []Option{opt("--json", "Emit the removed paths as JSON.")},
			},
			{
				Name:    "status",
				Summary: "Show the skill and link state.",
				Usage:   "idev agents status [--json]",
				Options: []Option{opt("--json", "Emit the state as JSON.")},
			},
		},
		Examples: []Example{
			example("Install the skill", "idev agents install"),
			example("Check the links", "idev agents status --json"),
		},
	}
}
