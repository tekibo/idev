package doctor

import (
	"os"
	"path/filepath"

	proj "github.com/tekibo/idev/internal/platform/project"
	"github.com/tekibo/idev/internal/platform/system"
)

func projectChecks(dir string) []Check {
	if dir == "" {
		return nil
	}
	info := proj.Inspect(dir)
	switch info.Kind {
	case proj.SwiftUI:
		return []Check{
			projectPath("project.xtool_yml", info.XtoolFile, "add xtool.yml"),
			projectPath("project.package_swift", info.PackageFile, "add Package.swift"),
			projectFlag("project.app_entry", info.AppEntry, "add a SwiftUI @main App"),
			projectFlag("project.bundle_id", info.BundleID != "", "set bundleID in xtool.yml"),
			projectOptional("project.icon", info.IconPath != "", "add an AppIcon asset"),
			projectOptional("project.info_plist", info.InfoPath != "", "add a source Info.plist for versioning"),
		}
	case proj.Expo:
		expo := filepath.Join(dir, "node_modules", ".bin", "expo")
		return []Check{
			projectFile("project.expo_cli", expo, "install the project dependencies in the app directory"),
		}
	case proj.ReactNative:
		cli := filepath.Join(dir, "node_modules", ".bin", "react-native")
		return []Check{
			projectFile("project.reactnative_cli", cli, "install the project dependencies in the app directory"),
		}
	}
	return nil
}

func projectPath(id, path, hint string) Check {
	ok := path != "" && fileExists(path)
	detail := ""
	if ok {
		detail = path
	}
	return Check{
		ID:        id,
		Class:     "project",
		Status:    status(ok),
		FixableBy: system.FixUser,
		Detail:    detail,
		Hint:      hint,
		Critical:  true,
	}
}

func projectFlag(id string, ok bool, hint string) Check {
	detail := ""
	if ok {
		detail = "found"
	}
	return Check{
		ID:        id,
		Class:     "project",
		Status:    status(ok),
		FixableBy: system.FixUser,
		Detail:    detail,
		Hint:      hint,
		Critical:  true,
	}
}

func projectOptional(id string, ok bool, hint string) Check {
	check := projectFlag(id, ok, hint)
	check.Critical = false
	return check
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
