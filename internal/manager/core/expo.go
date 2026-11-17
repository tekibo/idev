package core

import (
	"strings"

	"github.com/tekibo/idev/internal/platform/execx"
	"github.com/tekibo/idev/internal/tools/reactnative/expo"
)

// ExpoCommand builds the project Expo CLI command.
func (s *Services) ExpoCommand(args []string) (execx.Command, error) {
	return expo.Command(s.Config, s.expoDeviceArgs(args))
}

// expoDeviceArgs selects the only connected physical device for run:ios.
// Expo resolves its own destination and otherwise falls back to a simulator
// even when xtool already sees one physical device.
func (s *Services) expoDeviceArgs(args []string) []string {
	if len(args) == 0 || args[0] != "run:ios" {
		return args
	}
	for _, arg := range args {
		if arg == "--device" || strings.HasPrefix(arg, "--device=") {
			return args
		}
	}
	devices, err := s.ConnectedDevices()
	if err != nil || len(devices) != 1 || devices[0].UDID == "" {
		return args
	}
	next := append([]string{}, args...)
	return append(next, "--device", devices[0].UDID)
}
