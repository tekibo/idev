package project

import (
	"os"
	"path/filepath"
	"regexp"
)

var (
	packageName = regexp.MustCompile(`Package\s*\(\s*name:\s*"([^"]+)"`)
	iosMajor    = regexp.MustCompile(`\.iOS\(\.v([0-9]+)\)`)
	iosQuoted   = regexp.MustCompile(`\.iOS\("([^"]+)"\)`)
)

func packageProduct(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	if match := packageName.FindStringSubmatch(string(data)); len(match) == 2 {
		return match[1]
	}
	return ""
}

func packageIOS(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	if match := iosMajor.FindStringSubmatch(string(data)); len(match) == 2 {
		return match[1] + ".0"
	}
	if match := iosQuoted.FindStringSubmatch(string(data)); len(match) == 2 {
		return match[1]
	}
	return ""
}

func inspectSwiftPackage(dir string) Info {
	path := filepath.Join(dir, "Package.swift")
	name := filepath.Base(dir)
	if data, err := os.ReadFile(path); err == nil {
		if match := packageName.FindStringSubmatch(string(data)); len(match) == 2 {
			name = match[1]
		}
	}
	info := Info{
		Kind:        SwiftPackage,
		Dir:         dir,
		Name:        name,
		Label:       LabelFor(SwiftPackage),
		PackageFile: path,
	}
	addWarning(&info, "This is a Swift package, not an iOS app. Add xtool.yml to make it an xtool project.")
	return info
}
