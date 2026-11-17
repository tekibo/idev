// Package cli owns command parsing and text presentation.
package cli

import helptree "github.com/tekibo/idev/internal/cli/help"

// JSON emits machine-readable command events.
var JSON bool

// Run dispatches the CLI.
func Run(args []string) int {
	if len(args) > 0 && args[0] == "--json" {
		JSON = true
		args = args[1:]
	}
	if len(args) == 0 {
		return helptree.Show(nil, JSON)
	}
	if args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		return helptree.Show(args[1:], JSON)
	}
	if handled, code := helptree.PreDispatch(args, JSON); handled {
		return code
	}
	return dispatch(args)
}
