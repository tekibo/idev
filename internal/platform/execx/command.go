// Package execx owns external command construction and execution.
package execx

import (
	"io"
	"strings"
)

// Command is one external command.
type Command struct {
	Name   string
	Args   []string
	Dir    string
	Env    []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Preview returns a shell-like display of the command.
func (c Command) Preview() string {
	parts := append([]string{c.Name}, c.Args...)
	for i, part := range parts {
		if strings.ContainsAny(part, " \t\n\"'") {
			parts[i] = "'" + strings.ReplaceAll(part, "'", "'\\''") + "'"
		}
	}
	return strings.Join(parts, " ")
}
