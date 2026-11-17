package next

import (
	"github.com/tekibo/idev/internal/platform/project"
)

func projectSteps(info project.Info) (Step, *Step) {
	if !info.Known() {
		return Step{
			ID:      "project",
			Label:   "Project",
			Detail:  "no supported project in " + info.Dir,
			Command: "idev new <name>",
		}, nil
	}
	projectStep := Step{
		ID:     "project",
		Label:  "Project: " + info.Label + " " + info.Name,
		Done:   true,
		Detail: info.Dir,
	}
	if info.BundleID != "" || !info.BundleIDEditable {
		return projectStep, nil
	}
	return projectStep, &Step{
		ID:      "bundle",
		Label:   "Bundle ID",
		Detail:  "xtool.yml has no bundle identifier",
		Command: "idev project set-bundle-id --bundle-id <reverse.dns>",
	}
}

func action(info project.Info) string {
	switch info.Kind {
	case project.SwiftUI:
		return "idev swift run"
	case project.Expo:
		return "idev expo run:ios"
	case project.ReactNative:
		return "idev rn run-ios"
	case project.Flutter:
		return "idev flutter build"
	case project.Xcode:
		return "idev xcodeproj2xtool"
	case project.SwiftPackage:
		return "idev swift build"
	}
	return "idev doctor"
}
