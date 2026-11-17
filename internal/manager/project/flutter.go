package project

import (
	"github.com/tekibo/idev/internal/platform/execx"
	"github.com/tekibo/idev/internal/tools/flutter"
)

// FlutterCommand builds one pinned Flutter script command.
func (s *Service) FlutterCommand(command string, args []string) (execx.Command, error) {
	return flutter.New(s.Config, s.Runner).Command(command+".sh", args)
}

// FlutterPreview returns the command preview for one Flutter script.
func (s *Service) FlutterPreview(command string, args []string) string {
	client := flutter.New(s.Config, s.Runner)
	return client.Preview(command+".sh", args)
}
