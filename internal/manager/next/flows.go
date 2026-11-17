package next

import (
	"strings"

	"github.com/tekibo/idev/internal/platform/config"
	"github.com/tekibo/idev/internal/platform/project"
)

func flowsStep(cfg config.Config, info project.Info) Step {
	flows, err := config.LoadFlows(cfg.Home)
	if err != nil {
		return Step{ID: "flows", Label: "Flows", Detail: err.Error(), Command: "idev setup --flow <name>"}
	}
	if len(flows.Flows) == 0 {
		flow := flowFor(info.Kind)
		if flow == "" {
			return Step{ID: "flows", Label: "Flows", Done: true, Detail: "no optional flow needed"}
		}
		return Step{ID: "flows", Label: "Flows", Detail: "no flow selected", Command: "idev setup --flow " + flow}
	}
	missing := []string{}
	for _, flow := range flows.Flows {
		if !flowInstalled(cfg, flow) {
			missing = append(missing, flow)
		}
	}
	if len(missing) == 0 {
		return Step{ID: "flows", Label: "Flows", Done: true, Detail: strings.Join(flows.Flows, ", ")}
	}
	return Step{
		ID:      "flows",
		Label:   "Flows",
		Detail:  "missing: " + strings.Join(missing, ", "),
		Command: "idev setup --flow " + strings.Join(flows.Flows, " "),
	}
}

func flowFor(kind project.Kind) string {
	switch kind {
	case project.Expo, project.ReactNative:
		return "reactnative"
	case project.Flutter:
		return "flutter"
	}
	return ""
}

func flowInstalled(cfg config.Config, flow string) bool {
	switch flow {
	case "flutter":
		return config.IsFile(cfg.Path("flutter", "setup.sh"))
	case "kmp":
		return config.IsFile(cfg.Path("kmp", "build-framework.sh"))
	case "reactnative":
		return config.IsFile(cfg.Path("react-native", "tools", "node-darwin.js"))
	}
	return true
}
