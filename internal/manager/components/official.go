package components

import (
	"context"
	"fmt"
	"github.com/tekibo/idev/internal/manager/core"
	"runtime"

	"github.com/tekibo/idev/internal/platform/release"
	"github.com/tekibo/idev/internal/tools/installers/official"
)

// OfficialSpec is one official binary install plan.
type OfficialSpec = official.Spec

// OfficialResult is one official binary install result.
type OfficialResult = official.Result

// RcodesignPlan returns the pinned rcodesign install plan.
func (s *Service) RcodesignPlan() (OfficialSpec, error) {
	return s.officialPlan("rcodesign", []string{"--version"})
}

// IpswPlan returns the pinned ipsw install plan.
func (s *Service) IpswPlan() (OfficialSpec, error) {
	return s.officialPlan("ipsw", []string{"version"})
}

// InstallRcodesign downloads and verifies rcodesign.
func (s *Service) InstallRcodesign(ctx context.Context, force bool) (OfficialResult, error) {
	spec, err := s.RcodesignPlan()
	if err != nil {
		return OfficialResult{}, err
	}
	result, err := official.Install(ctx, s.Config, s.Runner, spec, force)
	if err != nil {
		return result, err
	}
	if err := core.RecordComponents(s.Config.Home, map[string]string{"rcodesign": result.Version}); err != nil {
		return result, err
	}
	return result, nil
}

// InstallIpsw downloads and verifies ipsw.
func (s *Service) InstallIpsw(ctx context.Context, force bool) (OfficialResult, error) {
	spec, err := s.IpswPlan()
	if err != nil {
		return OfficialResult{}, err
	}
	result, err := official.Install(ctx, s.Config, s.Runner, spec, force)
	if err != nil {
		return result, err
	}
	if err := core.RecordComponents(s.Config.Home, map[string]string{"ipsw": result.Version}); err != nil {
		return result, err
	}
	return result, nil
}

func (s *Service) officialPlan(name string, versionArgs []string) (OfficialSpec, error) {
	if s.Config.SourcesErr != nil {
		return OfficialSpec{}, s.Config.SourcesErr
	}
	manifest, err := release.Load(s.Config.Sources.Release.Manifest)
	if err != nil {
		return OfficialSpec{}, err
	}
	component, ok := manifest.Components[name+"-"+runtime.GOARCH]
	if !ok {
		return OfficialSpec{}, fmt.Errorf("release: %s component missing", name)
	}
	return official.NewSpec(s.Config, name, s.Config.Sources.Release.Manifest, component, versionArgs), nil
}
