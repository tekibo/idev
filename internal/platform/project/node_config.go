package project

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type appConfig struct {
	Name     string
	BundleID string
	Source   string
}

func staticAppConfig(dir string) (appConfig, bool) {
	for _, name := range []string{"app.json", "app.config.json"} {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			continue
		}
		var file struct {
			Expo struct {
				Name string `json:"name"`
				IOS  struct {
					BundleIdentifier string `json:"bundleIdentifier"`
				} `json:"ios"`
			} `json:"expo"`
		}
		if json.Unmarshal(data, &file) != nil {
			continue
		}
		return appConfig{Name: file.Expo.Name, BundleID: file.Expo.IOS.BundleIdentifier, Source: name}, true
	}
	return appConfig{}, false
}

func nodeBundleID(dir, projectPath string) (string, string) {
	if config, ok := staticAppConfig(dir); ok && config.BundleID != "" {
		return config.BundleID, config.Source
	}
	if id := infoPlistBundleID(filepath.Join(dir, "ios")); id != "" {
		return id, "Info.plist"
	}
	if id := projectBundleID(projectPath); id != "" {
		return id, "project.pbxproj"
	}
	return "", ""
}

func packageManager(workspaceDir, appDir string) string {
	dirs := []string{appDir}
	if workspaceDir != "" && workspaceDir != appDir {
		dirs = append(dirs, workspaceDir)
	}
	for _, dir := range dirs {
		for _, item := range []struct {
			File string
			Name string
		}{
			{"pnpm-lock.yaml", "pnpm"},
			{"yarn.lock", "yarn"},
			{"package-lock.json", "npm"},
			{"bun.lock", "bun"},
			{"bun.lockb", "bun"},
		} {
			if hasFile(dir, item.File) {
				return item.Name
			}
		}
	}
	return ""
}
