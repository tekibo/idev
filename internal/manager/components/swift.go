package components

import (
	"context"
	"github.com/tekibo/idev/internal/manager/core"
	"io"
	"runtime"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/tools/installers/swift"
)

// SwiftSpec is the official Swift install plan.
type SwiftSpec = swift.Spec

// SwiftResult is the official Swift install result.
type SwiftResult = swift.Result

// SwiftPlan returns the official Swift install plan.
func (s *Service) SwiftPlan() (SwiftSpec, error) {
	if s.Config.SourcesErr != nil {
		return SwiftSpec{}, s.Config.SourcesErr
	}
	return swift.ForArch(s.Config, runtime.GOARCH)
}

// InstallSwift installs the official Swift toolchain.
func (s *Service) InstallSwift(ctx context.Context, force bool, stdout io.Writer) (SwiftResult, error) {
	if s.Config.SourcesErr != nil {
		return SwiftResult{}, s.Config.SourcesErr
	}
	if err := config.EnsureSources(s.Config.Home); err != nil {
		return SwiftResult{}, err
	}
	result, err := swift.Install(ctx, s.Config, s.Runner, force, stdout)
	if err != nil {
		return result, err
	}
	if err := core.RecordComponents(s.Config.Home, map[string]string{"swift": result.Version}); err != nil {
		return result, err
	}
	return result, nil
}
