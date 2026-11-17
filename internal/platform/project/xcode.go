package project

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var bundleSetting = regexp.MustCompile(`PRODUCT_BUNDLE_IDENTIFIER = ([^;]+);`)

func inspectXcode(dir, path string) Info {
	name := strings.TrimSuffix(filepath.Base(path), ".xcodeproj")
	info := Info{
		Kind:     Xcode,
		Dir:      dir,
		Name:     name,
		Label:    LabelFor(Xcode),
		BundleID: pbxBundleID(filepath.Join(path, "project.pbxproj")),
	}
	addWarning(&info, "Xcode projects are not supported on Linux; convert with xcodeproj2xtool.")
	return info
}

func pbxBundleID(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	match := bundleSetting.FindStringSubmatch(string(data))
	if len(match) < 2 {
		return ""
	}
	value := strings.TrimSpace(strings.Trim(match[1], `"`))
	if value == "" || strings.Contains(value, "$(") {
		return ""
	}
	return value
}
