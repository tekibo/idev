package help

func configCommand() *Command {
	return &Command{
		Name:        "config",
		Group:       "Maintenance",
		Summary:     "Inspect or migrate the toolchain source config.",
		Usage:       "idev config <migrate|refresh>",
		Description: "Merges $IDEV_HOME/config/sources.json against the applied defaults. Conflicting or unsupported pins must be confirmed.",
		Children: []*Command{
			{
				Name:    "migrate",
				Summary: "Merge new defaults and confirm conflicting pins.",
				Usage:   "idev config migrate [--yes] [--dry-run]",
				Options: []Option{
					opt("--yes", "Accept required values and new defaults without prompting."),
					opt("--dry-run", "Print the conflicts and warnings without writing."),
				},
			},
			{
				Name:    "refresh",
				Summary: "Reset sources.json to the embedded defaults.",
				Usage:   "idev config refresh",
			},
		},
	}
}
