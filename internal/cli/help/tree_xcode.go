package help

func xcodeCommand() *Command {
	return &Command{
		Name:    "xcode",
		Group:   "Setup",
		Summary: "Manage the idev-owned Xcode.xip files.",
		Usage:   "idev xcode <command> [options]",
		Children: []*Command{
			{
				Name:    "list",
				Summary: "List the available Xcode.xip files.",
				Usage:   "idev xcode list [--json]",
				Options: []Option{opt("--json", "Emit the list as JSON.")},
			},
			{
				Name:    "import",
				Summary: "Move an Xcode.xip into the idev folder.",
				Usage:   "idev xcode import <path> [--replace] [--json]",
				Options: []Option{
					opt("--replace", "Replace an existing Xcode.xip."),
					opt("--json", "Emit the list as JSON."),
				},
			},
			{Name: "select", Summary: "Select one Xcode.xip.", Usage: "idev xcode select <path>"},
			{Name: "remove", Summary: "Remove one Xcode.xip.", Usage: "idev xcode remove <path>"},
		},
	}
}
