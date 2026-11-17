package components

import (
	"context"
	"fmt"
	"github.com/tekibo/idev/internal/manager/core"
	"io"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/release"
	xtoolcomponent "github.com/tekibo/idev/internal/tools/installers/xtool"
)

// XtoolSpec is the pinned xtool install plan.
type XtoolSpec = xtoolcomponent.Spec

// XtoolResult is the pinned xtool install result.
type XtoolResult = xtoolcomponent.Result

// XtoolPlan returns the pinned xtool install plan.
func (s *Service) XtoolPlan() (XtoolSpec, error) {
	if s.Config.SourcesErr != nil {
		return XtoolSpec{}, s.Config.SourcesErr
	}
	spec := xtoolcomponent.NewSpec(s.Config)
	manifest, err := release.Load(s.Config.Sources.Release.Manifest)
	if err != nil {
		return XtoolSpec{}, err
	}
	component, ok := manifest.Components["xtool"]
	if !ok {
		return XtoolSpec{}, fmt.Errorf("release: xtool component missing")
	}
	if component.Version != spec.Pin {
		return XtoolSpec{}, fmt.Errorf("release: xtool %s does not match lock %s", component.Version, spec.Pin)
	}
	spec.Archive = release.Resolve(s.Config.Sources.Release.Manifest, component.URL)
	spec.SHA256 = component.SHA256
	return spec, nil
}

// InstallXtool builds and links the pinned xtool.
func (s *Service) InstallXtool(ctx context.Context, force bool, stdout, stderr io.Writer) (XtoolResult, error) {
	if s.Config.SourcesErr != nil {
		return XtoolResult{}, s.Config.SourcesErr
	}
	if err := config.EnsureSources(s.Config.Home); err != nil {
		return XtoolResult{}, err
	}
	result, err := xtoolcomponent.Install(ctx, s.Config, s.Runner, force, stdout, stderr)
	if err != nil {
		return result, err
	}
	recorded := map[string]string{"xtool": result.Pin}
	if result.Shim != "" {
		recorded["xtool-shim"] = result.Shim
	}
	if err := core.RecordComponents(s.Config.Home, recorded); err != nil {
		return result, err
	}
	return result, nil
}
