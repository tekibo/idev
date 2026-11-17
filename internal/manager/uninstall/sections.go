package uninstall

import "path/filepath"

// userConfigs returns the preserved user data.
func userConfigs(userHome string) []Config {
	return []Config{
		{Path: filepath.Join(userHome, ".local", "share", "xtool"), Contents: "Apple ID auth token"},
		{Path: filepath.Join(userHome, ".config", "xtool"), Contents: "anisette and settings"},
		{Path: filepath.Join(userHome, ".pymobiledevice3"), Contents: "device pairing records"},
		{Path: "/var/lib/lockdown", Contents: "system pairing records"},
		{Path: filepath.Join(userHome, ".appstoreconnect"), Contents: "App Store Connect keys"},
		{Path: filepath.Join(userHome, ".cache", "xtool"), Contents: "SDK and DDI cache"},
	}
}

// projectData returns the user project paths that are never removed.
func projectData() []string {
	return []string{"Project files", "node_modules", "ios/Pods", "build outputs"}
}

// systemPackages maps the app-owned changeset package list to the report rows.
func systemPackages(names []string) []SystemPackage {
	packages := make([]SystemPackage, 0, len(names))
	for _, name := range names {
		packages = append(packages, SystemPackage{Package: name})
	}
	return packages
}
