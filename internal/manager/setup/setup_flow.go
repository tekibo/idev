package setup

import (
	"context"
	"fmt"
	"io"

	"github.com/tekibo/idev/internal/platform/config"
)

// SetupFlowPlan is the resolved setup plan.
type SetupFlowPlan = Plan

// LoadFlowSelection reads the saved setup selection.
func (s *Service) LoadFlowSelection() (FlowSelection, error) {
	selection, err := config.LoadFlows(s.Config.Home)
	if err != nil {
		return FlowSelection{}, err
	}
	return FlowSelection{Flows: normalizeFlows(selection.Flows), Expo: selection.Expo}, nil
}

func normalizeFlows(flows []string) []string {
	allowed := map[string]bool{"flutter": true, "kmp": true, "reactnative": true}
	result := make([]string, 0, len(flows))
	for _, flow := range flows {
		if allowed[flow] {
			result = append(result, flow)
		}
	}
	return result
}

// SaveFlowSelection writes the setup selection.
func (s *Service) SaveFlowSelection(selection FlowSelection) error {
	return config.SaveFlows(s.Config.Home, config.FlowSelection{Flows: selection.Flows, Expo: selection.Expo})
}

// SetupFlowPlan validates the selection and returns the plan.
func (s *Service) SetupFlowPlan(selection FlowSelection) (SetupFlowPlan, error) {
	return Resolve(config.FlowSelection{Flows: selection.Flows, Expo: selection.Expo})
}

// ApplySetupFlow installs the selected payloads and components.
func (s *Service) ApplySetupFlow(ctx context.Context, plan SetupFlowPlan, force bool, stdout, stderr io.Writer) error {
	if _, err := s.SetupPayloads(ctx, plan.Payloads, stdout); err != nil {
		return err
	}
	for _, component := range plan.Components {
		var err error
		switch component {
		case "swift":
			_, err = s.Components.InstallSwift(ctx, force, stdout)
		case "xtool":
			_, err = s.Components.InstallXtool(ctx, force, stdout, stderr)
		case "pymobiledevice3":
			_, err = s.Components.InstallPymobile(ctx, force, stdout, stderr)
		case "cocoapods":
			_, err = s.Components.InstallCocoaPods(ctx, force, stdout, stderr)
		case "rcodesign":
			_, err = s.Components.InstallRcodesign(ctx, force)
		case "ipsw":
			_, err = s.Components.InstallIpsw(ctx, force)
		case "darwin-tools":
			_, err = s.Components.InstallAppleTools(ctx, force, stdout, stderr)
		case "swift-build":
			_, err = s.Components.InstallSwiftBuild(ctx, force, stdout, stderr)
		default:
			return fmt.Errorf("setup: unknown component %q", component)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
