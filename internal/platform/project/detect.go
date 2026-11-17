package project

import (
	"path/filepath"
)

// Detect classifies one project directory.
func Detect(dir string) Info {
	if info, ok := detectNonNode(dir); ok {
		return info
	}
	if info, ok := inspectNode(dir, dir); ok {
		return info
	}
	return unknownInfo(dir)
}

// Inspect returns the detected project fields.
func Inspect(dir string) Info {
	return Detect(dir)
}

func detectNonNode(dir string) (Info, bool) {
	if hasFile(dir, "xtool.yml") {
		return inspectXtool(dir), true
	}
	if hasFile(dir, "pubspec.yaml") {
		return inspectFlutter(dir), true
	}
	if path := hasXcodeProject(dir); path != "" {
		return inspectXcode(dir, path), true
	}
	if hasFile(dir, "Package.swift") {
		return inspectSwiftPackage(dir), true
	}
	return Info{}, false
}

func unknownInfo(dir string) Info {
	return Info{
		Kind:  Unknown,
		Dir:   dir,
		Name:  filepath.Base(dir),
		Label: LabelFor(Unknown),
	}
}
