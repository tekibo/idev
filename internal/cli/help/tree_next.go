package help

func nextCommand() *Command {
	return &Command{
		Name:        "next",
		Group:       "Setup",
		Summary:     "Show the setup state and the exact next command.",
		Usage:       "idev next [--json] [--project DIR]",
		Description: "Reads the system plan, source pins, private toolchain, Darwin SDK, flows, and the project. It never changes the machine.",
		Options: []Option{
			opt("--json", "Emit the steps and the next command as JSON."),
			opt("--project DIR", "Project directory. Default: IDEV_PROJECT or the current directory."),
		},
		Examples: []Example{
			example("Show the next setup step", "idev next"),
			example("Machine-readable state", "idev next --json"),
		},
	}
}
