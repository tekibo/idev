package core

import (
	"fmt"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/execx"
)

// PodCommand runs the private CocoaPods wrapper in the project. The wrapper
// owns GEM_HOME, the CocoaPods cache, and the Expo autolinking patch.
func (s *Services) PodCommand(args []string) (execx.Command, error) {
	pod := s.Config.Path("bin", "pod")
	if !config.IsExec(pod) {
		return execx.Command{}, fmt.Errorf("no private CocoaPods at %s; run idev setup --cocoapods", pod)
	}
	dir := s.Config.Project
	if !config.IsFile(filepath.Join(dir, "Podfile")) {
		if ios := filepath.Join(dir, "ios"); config.IsFile(filepath.Join(ios, "Podfile")) {
			dir = ios
		}
	}
	return execx.Command{
		Name: pod,
		Args: args,
		Dir:  dir,
		Env:  s.Config.ReactNativeEnv(""),
	}, nil
}
