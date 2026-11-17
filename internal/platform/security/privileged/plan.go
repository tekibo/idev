// Package privileged owns sudo plans and the masked password runner.
package privileged

import "github.com/tekibo/idev/internal/platform/execx"

// Plan describes one exact privileged command.
type Plan struct {
	Command  execx.Command
	Reason   string
	Packages []string
}

// Display returns the user-facing command.
func (p Plan) Display() string {
	return "sudo " + p.Command.Preview()
}
