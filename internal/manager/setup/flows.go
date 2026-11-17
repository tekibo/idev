// Package setup builds the flow-aware setup plan.
package setup

import (
	"fmt"

	"github.com/tekibo/idev/internal/platform/config"
)

// Plan is the resolved setup plan.
type Plan struct {
	Flows      []string `json:"flows,omitempty"`
	Expo       bool     `json:"expo"`
	Components []string `json:"components"`
	Payloads   []string `json:"payloads"`
}

var baseComponents = []string{
	"swift", "xtool", "pymobiledevice3", "rcodesign", "ipsw",
}
var basePayloads = []string{"apple-tools"}

var flows = map[string]struct {
	components []string
	payloads   []string
}{
	"flutter":     {components: []string{"cocoapods"}, payloads: []string{"flutter"}},
	"kmp":         {payloads: []string{"kmp"}},
	"reactnative": {components: []string{"cocoapods", "swift-build"}, payloads: []string{"react-native"}},
}

// BasePlan returns the shared base without a flow.
func BasePlan() Plan {
	return Plan{
		Payloads:   append([]string{}, basePayloads...),
		Components: append([]string{}, baseComponents...),
	}
}

// Resolve validates the selection and returns the setup plan.
func Resolve(selection config.FlowSelection) (Plan, error) {
	if len(selection.Flows) == 0 {
		return Plan{}, fmt.Errorf("setup: no flows selected")
	}
	plan := BasePlan()
	plan.Flows = append([]string{}, selection.Flows...)
	plan.Expo = selection.Expo
	seen := map[string]bool{}
	for _, component := range plan.Components {
		seen[component] = true
	}
	hasRN := false
	for _, name := range selection.Flows {
		flow, ok := flows[name]
		if !ok {
			return Plan{}, fmt.Errorf("setup: unknown flow %q", name)
		}
		if name == "reactnative" {
			hasRN = true
		}
		for _, component := range flow.components {
			if !seen[component] {
				seen[component] = true
				plan.Components = append(plan.Components, component)
			}
		}
		for _, payload := range flow.payloads {
			if !contains(plan.Payloads, payload) {
				plan.Payloads = append(plan.Payloads, payload)
			}
		}
	}
	if !hasRN {
		plan.Expo = false
	}
	return plan, nil
}

func contains(values []string, value string) bool {
	for _, candidate := range values {
		if candidate == value {
			return true
		}
	}
	return false
}
