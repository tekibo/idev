package auth

import (
	"io"

	"github.com/tekibo/idev/internal/platform/redact"
	"github.com/tekibo/idev/internal/tools/xtool"
)

// Auth runs xtool auth.
func (s *Service) Auth(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	return s.Xtool.Auth(args, stdin, stdout, stderr)
}

// AuthArgs validates and builds the xtool auth command.
func (s *Service) AuthArgs(args []string) ([]string, error) {
	if redact.HasSecretFlag(args) {
		return nil, xtool.ErrSecretArg
	}
	return append([]string{"auth"}, args...), nil
}
