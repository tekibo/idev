package components

import (
	"context"
	"github.com/tekibo/idev/internal/manager/core"
	"io"

	"github.com/tekibo/idev/internal/tools/installers/cocoapods"
)

// CocoaPodsSpec is the private gem tree install plan.
type CocoaPodsSpec = cocoapods.Spec

// CocoaPodsResult is one CocoaPods install result.
type CocoaPodsResult = cocoapods.Result

// CocoaPodsPlan returns the private gem tree paths.
func (s *Service) CocoaPodsPlan() CocoaPodsSpec {
	return cocoapods.NewSpec(s.Config)
}

// InstallCocoaPods creates the private gem tree and pod wrapper.
func (s *Service) InstallCocoaPods(ctx context.Context, force bool, stdout, stderr io.Writer) (CocoaPodsResult, error) {
	result, err := cocoapods.Install(ctx, s.Config, s.Runner, force, stdout, stderr)
	if err != nil {
		return result, err
	}
	if err := core.RecordComponents(s.Config.Home, map[string]string{"cocoapods": result.Version}); err != nil {
		return result, err
	}
	return result, nil
}
