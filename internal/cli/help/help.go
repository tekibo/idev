// Package help owns the CLI command tree and focused help output.
package help

import (
	"fmt"
	"os"
	"strings"
)

// Command is one node in the CLI help tree.
type Command struct {
	Name        string
	Aliases     []string
	Group       string
	Summary     string
	Usage       string
	Description string
	Passthrough bool
	Options     []Option
	Env         []Env
	Examples    []Example
	JSON        []string
	Children    []*Command
}

// Option is one command flag.
type Option struct {
	Flags       string `json:"flags"`
	Description string `json:"description"`
	Default     string `json:"default,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// Env is one environment variable.
type Env struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Example is one usage example.
type Example struct {
	Description string `json:"description,omitempty"`
	Command     string `json:"command"`
}

// Root builds the command tree.
func Root() *Command {
	return &Command{
		Name:     "idev",
		Summary:  "Native iOS build manager for xtool, Expo, React Native, Flutter, and KMP on Linux.",
		Usage:    "idev [--json] <command> [options]",
		Options:  globalOptions(),
		Env:      globalEnv(),
		Children: allCommands(),
	}
}

// Show prints help for one command path.
func Show(path []string, asJSON bool) int {
	filtered := make([]string, 0, len(path))
	for _, part := range path {
		switch part {
		case "--json":
			asJSON = true
		case "-h", "--help":
		default:
			filtered = append(filtered, part)
		}
	}
	node, ok := findPath(Root(), filtered)
	if !ok {
		fmt.Fprintf(os.Stderr, "idev help: unknown command %q\n", strings.Join(filtered, " "))
		return 2
	}
	if asJSON {
		return printJSON(node)
	}
	printText(node, filtered)
	return 0
}

func allCommands() []*Command {
	commands := []*Command{}
	commands = append(commands, projectCommands()...)
	commands = append(commands, buildCommands()...)
	commands = append(commands, setupCommands()...)
	commands = append(commands, maintenanceCommands()...)
	commands = append(commands, toolCommands()...)
	return commands
}

func globalOptions() []Option {
	return []Option{
		opt("--json", "Emit machine-readable output where the command supports it."),
		opt("-h, --help", "Show help for the command or subcommand."),
		opt("-v, --version", "Show the manager version."),
	}
}

func globalEnv() []Env {
	return []Env{
		{Name: "IDEV_HOME", Description: "Manager prefix. Default: ~/.local/share/idev."},
		{Name: "IDEV_PROJECT", Description: "Project directory. Default: the current directory."},
	}
}
