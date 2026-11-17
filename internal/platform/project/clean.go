package project

import "path/filepath"

// CleanPaths returns generated build paths for one project.
func CleanPaths(info Info) []string {
	if info.Kind != SwiftUI {
		return []string{
			filepath.Join(info.Dir, "build"),
			filepath.Join(info.Dir, "ios", "build"),
		}
	}
	paths := []string{filepath.Join(info.Dir, ".build")}
	for _, pattern := range []string{"xtool/*.app", "xtool/*.ipa", "xtool/.xtool-tmp"} {
		matches, _ := filepath.Glob(filepath.Join(info.Dir, pattern))
		paths = append(paths, matches...)
	}
	return paths
}
