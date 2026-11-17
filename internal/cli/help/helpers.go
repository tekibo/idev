package help

func opt(flags, description string) Option {
	return Option{Flags: flags, Description: description}
}

func def(flags, description, value string) Option {
	return Option{Flags: flags, Description: description, Default: value}
}

func req(flags, description string) Option {
	return Option{Flags: flags, Description: description, Required: true}
}

func example(description, command string) Example {
	return Example{Description: description, Command: command}
}
