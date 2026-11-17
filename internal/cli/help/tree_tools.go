package help

func toolCommands() []*Command {
	return append([]*Command{
		xtoolCommand(),
		dsCommand(),
		sdkCommand(),
		helpCommand(),
	}, append(upstreamCommands(), upstreamExtraCommands()...)...)
}

func xtoolCommand() *Command {
	return &Command{
		Name:        "xtool",
		Group:       "Tools",
		Summary:     "Run any xtool command with the private environment.",
		Usage:       "idev xtool <args...>",
		Passthrough: true,
	}
}

func dsCommand() *Command {
	return &Command{
		Name:        "ds",
		Group:       "Tools",
		Summary:     "Run xtool ds with the private environment.",
		Usage:       "idev ds <args...>",
		Passthrough: true,
	}
}

func sdkCommand() *Command {
	return &Command{
		Name:        "sdk",
		Group:       "Tools",
		Summary:     "Run xtool sdk and read the parsed status.",
		Usage:       "idev sdk <args...>",
		Passthrough: true,
		Description: "Use 'idev sdk status --json' for the parsed SDK status.",
		Children: []*Command{
			{Name: "status", Summary: "Show the SDK status.", Usage: "idev sdk status [--json]", Options: []Option{
				opt("--json", "Emit the parsed SDK status."),
			}},
		},
	}
}

func helpCommand() *Command {
	return &Command{
		Name:    "help",
		Group:   "Tools",
		Summary: "Show focused help for a command or subcommand.",
		Usage:   "idev help [command] [subcommand] [--json]",
		Options: []Option{opt("--json", "Emit the help tree as JSON.")},
		Examples: []Example{
			example("Project info help", "idev help project info"),
			example("Expo run:ios help", "idev help expo run:ios"),
		},
	}
}
