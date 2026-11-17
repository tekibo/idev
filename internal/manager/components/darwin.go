package components

import (
	"context"
	"io"

	"github.com/tekibo/idev/internal/tools/installers/darwin"
)

// AppleToolsSpec is the payload and SDK target for one install.
type AppleToolsSpec = darwin.Spec

// AppleToolsResult is one Darwin tools and OAM install result.
type AppleToolsResult = darwin.Result

// AppleToolsPlan returns the payload and SDK paths.
func (s *Service) AppleToolsPlan() (AppleToolsSpec, error) {
	if s.Config.SourcesErr != nil {
		return AppleToolsSpec{}, s.Config.SourcesErr
	}
	return darwin.NewSpec(s.Config), nil
}

// InstallAppleTools builds the vendored tools and copies them into the xtool SDK.
func (s *Service) InstallAppleTools(ctx context.Context, force bool, stdout, stderr io.Writer) (AppleToolsResult, error) {
	if s.Config.SourcesErr != nil {
		return AppleToolsResult{}, s.Config.SourcesErr
	}
	return darwin.Install(ctx, s.Config, s.Runner, force, stdout, stderr)
}
