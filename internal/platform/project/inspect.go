package project

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var skippedDirs = map[string]bool{
	".build": true, ".git": true, ".swiftpm": true,
	"xtool": true, "node_modules": true, "Pods": true,
}

func hasAppEntry(dir string) bool {
	found := false
	sources := filepath.Join(dir, "Sources")
	_ = filepath.WalkDir(sources, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || found {
			return nil
		}
		if entry.IsDir() {
			if skippedDirs[entry.Name()] {
				return fs.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".swift" {
			return nil
		}
		info, err := entry.Info()
		if err != nil || info.Size() > 1<<20 {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		text := string(data)
		if strings.Contains(text, "@main") &&
			(strings.Contains(text, ": App") || strings.Contains(text, ": SwiftUI.App")) {
			found = true
		}
		return nil
	})
	return found
}

func findIcon(dir string) string {
	var found string
	_ = filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || found != "" {
			return nil
		}
		if entry.IsDir() {
			if skippedDirs[entry.Name()] {
				return fs.SkipDir
			}
			if strings.HasSuffix(entry.Name(), ".appiconset") || strings.HasSuffix(entry.Name(), ".icon") {
				found = path
				return fs.SkipDir
			}
		}
		return nil
	})
	return found
}

func findInfoPlist(dir, name string) string {
	candidates := []string{
		filepath.Join(dir, "Info.plist"),
		filepath.Join(dir, "Sources", name, "Info.plist"),
	}
	for _, path := range candidates {
		if fileExists(path) {
			return path
		}
	}
	return ""
}
