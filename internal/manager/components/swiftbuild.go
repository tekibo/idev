package components

import (
	"context"
	"fmt"
	"github.com/tekibo/idev/internal/manager/core"
	"io"

	"github.com/tekibo/idev/internal/platform/release"
	"github.com/tekibo/idev/internal/tools/installers/swiftbuild"
)

// SwiftBuildSpec is the Swift Build and Expo engine install plan.
type SwiftBuildSpec = swiftbuild.Spec

// SwiftBuildResult is one Swift Build and Expo engine install result.
type SwiftBuildResult = swiftbuild.Result

// SwiftBuildPlan returns the pinned Swift Build source plan.
func (s *Service) SwiftBuildPlan() (SwiftBuildSpec, error) {
	if s.Config.SourcesErr != nil {
		return SwiftBuildSpec{}, s.Config.SourcesErr
	}
	manifest, err := release.Load(s.Config.Sources.Release.Manifest)
	if err != nil {
		return SwiftBuildSpec{}, err
	}
	component, ok := manifest.Components["swift-build"]
	if !ok {
		return SwiftBuildSpec{}, fmt.Errorf("release: swift-build component missing")
	}
	return swiftbuild.NewSpec(s.Config, component), nil
}

// InstallSwiftBuild builds Swift Build and the Expo engine.
func (s *Service) InstallSwiftBuild(ctx context.Context, force bool, stdout, stderr io.Writer) (SwiftBuildResult, error) {
	spec, err := s.SwiftBuildPlan()
	if err != nil {
		return SwiftBuildResult{}, err
	}
	result, err := swiftbuild.Install(ctx, s.Config, s.Runner, spec, force, stdout, stderr)
	if err != nil {
		return result, err
	}
	if err := core.RecordComponents(s.Config.Home, map[string]string{"swift-build": result.Version}); err != nil {
		return result, err
	}
	return result, nil
}
