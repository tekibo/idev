package help

import "fmt"

func writeOptions(options []Option) {
	if len(options) == 0 {
		return
	}
	fmt.Println("Options:")
	for _, option := range options {
		description := option.Description
		if option.Required {
			description += " (required)"
		}
		if option.Default != "" {
			description += " (default: " + option.Default + ")"
		}
		fmt.Printf("  %-24s %s\n", option.Flags, description)
	}
	fmt.Println()
}

func writeEnv(env []Env) {
	if len(env) == 0 {
		return
	}
	fmt.Println("Environment:")
	for _, item := range env {
		fmt.Printf("  %-18s %s\n", item.Name, item.Description)
	}
	fmt.Println()
}

func writeExamples(examples []Example) {
	if len(examples) == 0 {
		return
	}
	fmt.Println("Examples:")
	for _, item := range examples {
		if item.Description != "" {
			fmt.Printf("  # %s\n", item.Description)
		}
		fmt.Printf("  $ %s\n", item.Command)
	}
	fmt.Println()
}
