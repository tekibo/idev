package project

import (
	"os"
	"path/filepath"
	"strings"
)

func inspectFlutter(dir string) Info {
	info := Info{
		Kind:  Flutter,
		Dir:   dir,
		Name:  filepath.Base(dir),
		Label: LabelFor(Flutter),
	}
	if data, err := os.ReadFile(filepath.Join(dir, "pubspec.yaml")); err == nil {
		for _, raw := range strings.Split(string(data), "\n") {
			line := strings.TrimSpace(raw)
			if strings.HasPrefix(raw, " ") || strings.HasPrefix(raw, "\t") {
				continue
			}
			key, value, ok := strings.Cut(line, ":")
			if !ok {
				continue
			}
			if strings.TrimSpace(key) == "name" && strings.TrimSpace(value) != "" {
				info.Name = scalar(value)
			}
		}
	}
	info.BundleID = pbxBundleID(filepath.Join(dir, "ios", "Runner.xcodeproj", "project.pbxproj"))
	if info.BundleID == "" {
		addWarning(&info, "No iOS bundle ID was found in the Flutter project.")
	}
	return info
}
