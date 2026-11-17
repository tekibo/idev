package config

import (
	"os"
	"path/filepath"
	"strings"
)

// SwiftRoot returns the private Swift toolchain root.
func (c Config) SwiftRoot() string {
	version := c.Sources.Swift.Version
	if version == "" {
		version = "unknown"
	}
	return c.Path("toolchains", "swift-"+version)
}

// SwiftBin returns the private Swift bin directory.
func (c Config) SwiftBin() string {
	return filepath.Join(c.SwiftRoot(), "usr", "bin")
}

// SwiftLib returns the private Swift runtime directory.
func (c Config) SwiftLib() string {
	return filepath.Join(c.SwiftRoot(), "usr", "lib", "swift", "linux")
}

// BaseEnv returns the manager environment without a Darwin SDK.
func (c Config) BaseEnv() []string {
	return EnvWith(c.baseOverrides())
}

func (c Config) baseOverrides() map[string]string {
	overrides := map[string]string{
		"IDEV_HOME":      c.Home,
		"IDEV_SWIFT_BIN": c.SwiftBin(),
		"PATH":           c.SwiftBin() + string(os.PathListSeparator) + c.Path("bin") + string(os.PathListSeparator) + os.Getenv("PATH"),
	}
	for key, value := range c.StateEnv() {
		overrides[key] = value
	}
	libs := c.SwiftLib() + string(os.PathListSeparator) + c.Path("lib", "compat") + string(os.PathListSeparator) + c.Path("lib")
	if existing := os.Getenv("LD_LIBRARY_PATH"); existing != "" {
		libs += string(os.PathListSeparator) + existing
	}
	overrides["LD_LIBRARY_PATH"] = libs
	return overrides
}

// SwiftEnv returns the private Swift environment.
func (c Config) SwiftEnv() []string {
	return c.SwiftEnvWith(nil)
}

// SwiftEnvWith returns the private Swift environment with extra overrides.
func (c Config) SwiftEnvWith(extra map[string]string) []string {
	overrides := c.baseOverrides()
	for key, value := range extra {
		overrides[key] = value
	}
	if sdk := c.SDKPath(); sdk != "" {
		overrides["DEVELOPER_DIR"] = sdk
		overrides["DARWIN_SDK_BUNDLE"] = c.SDKRoot()
	}
	return EnvWith(overrides)
}

// SwiftBuildEnv returns the private Swift environment without DEVELOPER_DIR.
// SwiftPM's SwiftBuild backend registers the SDK platforms itself; a
// DEVELOPER_DIR makes it register the same platforms twice.
func (c Config) SwiftBuildEnv() []string {
	env := c.SwiftEnv()
	filtered := make([]string, 0, len(env))
	for _, entry := range env {
		if strings.HasPrefix(entry, "DEVELOPER_DIR=") {
			continue
		}
		filtered = append(filtered, entry)
	}
	return filtered
}

// XtoolBuildEnv returns the private Swift environment with the xtool version pin.
func (c Config) XtoolBuildEnv() []string {
	return c.SwiftEnvWith(map[string]string{"XTOOL_VERSION": c.Lock.AppleDev.XtoolSHA})
}

func (c Config) swiftPaths() (string, string) {
	if value := os.Getenv("IDEV_SWIFT_BIN"); value != "" {
		return value, filepath.Join(value, "..", "lib", "swift", "linux")
	}
	binary := c.SwiftBin()
	if IsExec(filepath.Join(binary, "swift")) {
		return binary, c.SwiftLib()
	}
	return "", ""
}
