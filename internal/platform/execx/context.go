package execx

import "context"

// ContextRunner is a runner that cancels commands with a context.
type ContextRunner interface {
	Runner
	RunContext(context.Context, Command) error
	OutputContext(context.Context, Command) ([]byte, error)
}
