package core

import (
	"io"

	"github.com/tekibo/idev/internal/platform/execx"
)

// PromptTerminal restores sane terminal modes around a privileged prompt.
func (s *Services) PromptTerminal(in io.Reader) (func() error, error) {
	return execx.PromptTerminal(in)
}
