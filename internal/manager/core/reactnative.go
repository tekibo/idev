package core

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
	"github.com/tekibo/idev/internal/tools/reactnative/expo"
)

// ReactNativeCommand builds the project React Native CLI command with the
// same environment and shims the Expo arm uses.
func (s *Services) ReactNativeCommand(args []string) (execx.Command, error) {
	cli := filepath.Join(s.Config.Project, "node_modules", ".bin", "react-native")
	if !config.IsExec(cli) {
		return execx.Command{}, fmt.Errorf("no React Native CLI at %s; install project dependencies first", cli)
	}
	hook, err := expo.HookPath(s.Config.Home)
	if err != nil {
		return execx.Command{}, err
	}
	return execx.Command{
		Name: cli,
		Args: s.reactNativeDeviceArgs(args),
		Dir:  s.Config.Project,
		Env:  s.Config.ReactNativeEnv(hook),
	}, nil
}

// reactNativeDeviceArgs selects the only connected physical device for
// run-ios. Without a destination the React Native CLI starts a simulator,
// which does not exist on Linux, even when xtool already sees one device.
func (s *Services) reactNativeDeviceArgs(args []string) []string {
	if len(args) == 0 || args[0] != "run-ios" {
		return args
	}
	for _, arg := range args {
		if arg == "--device" || strings.HasPrefix(arg, "--device=") ||
			arg == "--udid" || strings.HasPrefix(arg, "--udid=") ||
			arg == "--simulator" || strings.HasPrefix(arg, "--simulator=") ||
			arg == "--list-devices" || arg == "-i" || arg == "--interactive" {
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
