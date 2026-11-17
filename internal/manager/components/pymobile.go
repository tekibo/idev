package components

import (
	"context"
	"github.com/tekibo/idev/internal/manager/core"
	"io"

	"github.com/tekibo/idev/internal/tools/installers/pymobile"
)

// PymobileSpec is the private pymobiledevice3 install plan.
type PymobileSpec = pymobile.Spec

// PymobileResult is one pymobiledevice3 install result.
type PymobileResult = pymobile.Result

// PymobilePlan returns the private venv paths.
func (s *Service) PymobilePlan() PymobileSpec {
	return pymobile.NewSpec(s.Config)
}

// InstallPymobile creates the private venv and pins pymobiledevice3.
func (s *Service) InstallPymobile(ctx context.Context, force bool, stdout, stderr io.Writer) (PymobileResult, error) {
	result, err := pymobile.Install(ctx, s.Config, s.Runner, force, stdout, stderr)
	if err != nil {
		return result, err
	}
	if err := core.RecordComponents(s.Config.Home, map[string]string{"pymobiledevice3": result.Version}); err != nil {
		return result, err
	}
	return result, nil
}
