package help

func newCommand() *Command {
	return &Command{
		Name:    "new",
		Group:   "Project",
		Summary: "Create a new SwiftUI project.",
		Usage:   "idev new <name> [--bundle-id <id>] [--skip-setup]",
		Options: []Option{
			opt("--bundle-id ID", "Set the project bundle identifier."),
			opt("--skip-setup", "Pass xtool's non-interactive setup flag."),
		},
		Examples: []Example{
			example("Create a project with a bundle ID", "idev new MyApp --bundle-id com.example.myapp"),
		},
	}
}
