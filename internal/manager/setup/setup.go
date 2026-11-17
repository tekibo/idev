package setup

import (
	"context"
	"github.com/tekibo/idev/internal/manager/core"
	"io"
	"os"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/download"
	"github.com/tekibo/idev/internal/platform/release"
)

// SetupResult reports what the setup wrote.
type SetupResult struct {
	Home             string
	ReactNative      string
	AppleTools       string
	Written          int
	ManagerPath      string
	ManagerInstalled bool
}

// Setup downloads the release payloads and optionally the manager binary.
func (s *Service) Setup(ctx context.Context, force, installSelf bool, stdout io.Writer) (SetupResult, error) {
	result, err := s.SetupPayloads(ctx, []string{"react-native", "apple-tools"}, stdout)
	if err != nil {
		return result, err
	}
	if installSelf {
		path, installed, err := s.InstallSelf()
		if err != nil {
			return result, err
		}
		result.ManagerPath = path
		result.ManagerInstalled = installed
	}
	return result, nil
}

// SetupPayloads downloads the selected release payloads.
func (s *Service) SetupPayloads(ctx context.Context, names []string, stdout io.Writer) (SetupResult, error) {
	result := SetupResult{Home: s.Config.Home}
	if s.Config.Home == "" {
		return result, os.ErrNotExist
	}
	if s.Config.SourcesErr != nil {
		return result, s.Config.SourcesErr
	}
	if err := config.EnsureSources(s.Config.Home); err != nil {
		return result, err
	}
	manifest, err := release.Load(s.Config.Sources.Release.Manifest)
	if err != nil {
		return result, err
	}
	installed, _ := core.LoadComponents(s.Config.Home)
	versions := map[string]string{}
	for _, name := range names {
		component, ok := manifest.Components[name]
		if !ok {
			return result, os.ErrNotExist
		}
		target := s.Config.Path(name)
		switch name {
		case "react-native":
			result.ReactNative = target
		case "apple-tools":
			result.AppleTools = target
		}
		if installed[name] == component.Version {
			if info, statErr := os.Stat(target); statErr == nil && info.IsDir() {
				versions[name] = component.Version
				continue
			}
		}
		count, err := release.FetchProgress(ctx, s.Config.Sources.Release.Manifest, component, target, s.Config.Path("cache", "downloads"), download.Log(stdout, name))
		if err != nil {
			return result, err
		}
		result.Written += count
		versions[name] = component.Version
	}
	if err := core.RecordComponents(s.Config.Home, versions); err != nil {
		return result, err
	}
	return result, nil
}
