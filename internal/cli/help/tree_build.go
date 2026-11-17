package help

func buildCommands() []*Command {
	commands := []*Command{
		expoCommand(),
		rnCommand(),
		iosCommand(),
		startCommand(),
	}
	return append(commands, nativeCommands()...)
}
