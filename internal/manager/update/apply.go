package update

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/tekibo/idev/internal/manager/core"
	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/release"
)

// UpdateApply installs one component from the release manifest.
func (s *Service) UpdateApply(ctx context.Context, component string, stdout, stderr io.Writer) error {
	if s.Config.SourcesErr != nil {
		return s.Config.SourcesErr
	}
	if err := config.EnsureSources(s.Config.Home); err != nil {
		return err
	}
	manifestFile, err := release.Load(s.Config.Sources.Release.Manifest)
	if err != nil {
		return err
	}
	switch component {
	case "react-native", "apple-tools", "flutter", "kmp":
		return s.applyPayload(ctx, manifestFile, component, stdout, stderr)
	case "xtool":
		_, err := s.Components.InstallXtool(ctx, true, stdout, stderr)
		return err
	case "swift":
		_, err := s.Components.InstallSwift(ctx, true, stdout)
		return err
	case "rcodesign":
		_, err := s.Components.InstallRcodesign(ctx, true)
		return err
	case "ipsw":
		_, err := s.Components.InstallIpsw(ctx, true)
		return err
	case "swift-build":
		_, err := s.Components.InstallSwiftBuild(ctx, true, stdout, stderr)
		return err
	case "pymobiledevice3":
		_, err := s.Components.InstallPymobile(ctx, true, stdout, stderr)
		return err
	case "cocoapods":
		_, err := s.Components.InstallCocoaPods(ctx, true, stdout, stderr)
		return err
	case "idev":
		return s.applyManager(ctx, manifestFile, stdout, stderr)
	default:
		return fmt.Errorf("unknown component %q", component)
	}
}

// UpdateApplyAll installs every outdated component, then reconciles the manager.
func (s *Service) UpdateApplyAll(ctx context.Context, stdout, stderr io.Writer) error {
	updates, err := s.UpdateCheck()
	if err != nil {
		return err
	}
	for _, update := range updates {
		if update.Current || update.Component == "idev" {
			continue
		}
		if err := s.UpdateApply(ctx, update.Component, stdout, stderr); err != nil {
			return err
		}
	}
	return s.UpdateApply(ctx, "idev", stdout, stderr)
}

func (s *Service) applyManager(ctx context.Context, manifestFile release.Manifest,
	stdout, stderr io.Writer) error {
	name := "idev-" + runtime.GOOS + "-" + runtime.GOARCH
	item, ok := manifestFile.Components[name]
	if !ok {
		return fmt.Errorf("release: component %s missing", name)
	}
	target, err := os.Executable()
	if err != nil {
		return err
	}
	if err := release.FetchBinary(ctx, s.Config.Sources.Release.Manifest, item, "idev", target,
		s.Config.Path("cache", "downloads")); err != nil {
		return err
	}
	s.refreshAgents(ctx, target, stdout, stderr)
	return nil
}

func (s *Service) applyPayload(ctx context.Context, manifestFile release.Manifest, component string, stdout, stderr io.Writer) error {
	item, ok := manifestFile.Components[component]
	if !ok {
		return fmt.Errorf("release: component %s missing", component)
	}
	if _, err := release.Fetch(ctx, s.Config.Sources.Release.Manifest, item, s.Config.Path(component), s.Config.Path("cache", "downloads")); err != nil {
		return err
	}
	if component == "apple-tools" {
		if _, err := s.Components.InstallAppleTools(ctx, true, stdout, stderr); err != nil {
			return err
		}
	}
	return core.RecordComponents(s.Config.Home, map[string]string{component: item.Version})
}
