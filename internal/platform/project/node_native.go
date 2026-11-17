package project

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var pbxBundleIDPattern = regexp.MustCompile(`PRODUCT_BUNDLE_IDENTIFIER\s*=\s*([^;]+);`)

func nodeNative(dir string) (string, string, string) {
	iosDir := filepath.Join(dir, "ios")
	if info, err := os.Stat(iosDir); err != nil || !info.IsDir() {
		return "", "", ""
	}
	project := firstGlob(filepath.Join(iosDir, "*.xcodeproj"))
	scheme := firstScheme(project)
	if scheme == "" && project != "" {
		scheme = strings.TrimSuffix(filepath.Base(project), ".xcodeproj")
	}
	return iosDir, project, scheme
}

func firstGlob(pattern string) string {
	matches, _ := filepath.Glob(pattern)
	if len(matches) == 0 {
		return ""
	}
	return matches[0]
}

func firstScheme(project string) string {
	if project == "" {
		return ""
	}
	scheme := firstGlob(filepath.Join(project, "xcshareddata", "xcschemes", "*.xcscheme"))
	if scheme == "" {
		return ""
	}
	return strings.TrimSuffix(filepath.Base(scheme), ".xcscheme")
}

func infoPlistBundleID(iosDir string) string {
	candidates := []string{filepath.Join(iosDir, "Info.plist")}
	for _, pattern := range []string{
		filepath.Join(iosDir, "*", "Info.plist"),
		filepath.Join(iosDir, "*", "*", "Info.plist"),
	} {
		matches, _ := filepath.Glob(pattern)
		candidates = append(candidates, matches...)
	}
	for _, path := range candidates {
		if id := plistBundleID(path); id != "" {
			return id
		}
	}
	return ""
}

func plistBundleID(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	text := string(data)
	key := "<key>CFBundleIdentifier</key>"
	index := strings.Index(text, key)
	if index < 0 {
		return ""
	}
	rest := text[index+len(key):]
	start := strings.Index(rest, "<string>")
	if start < 0 {
		return ""
	}
	rest = rest[start+len("<string>"):]
	end := strings.Index(rest, "</string>")
	if end < 0 {
		return ""
	}
	id := strings.TrimSpace(rest[:end])
	if strings.Contains(id, "$") {
		return ""
	}
	return id
}

func projectBundleID(projectPath string) string {
	if projectPath == "" {
		return ""
	}
	data, err := os.ReadFile(filepath.Join(projectPath, "project.pbxproj"))
	if err != nil {
		return ""
	}
	match := pbxBundleIDPattern.FindStringSubmatch(string(data))
	if len(match) < 2 {
		return ""
	}
	id := strings.Trim(strings.TrimSpace(match[1]), `"`)
	if strings.Contains(id, "$") {
		return ""
	}
	return id
}
