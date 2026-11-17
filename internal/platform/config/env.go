package config

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// EnvWith returns the current environment with overrides applied.
func EnvWith(overrides map[string]string) []string {
	env := os.Environ()
	filtered := env[:0]
	for _, entry := range env {
		key, _, _ := strings.Cut(entry, "=")
		if _, replace := overrides[key]; !replace {
			filtered = append(filtered, entry)
		}
	}
	keys := make([]string, 0, len(overrides))
	for key := range overrides {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		filtered = append(filtered, key+"="+overrides[key])
	}
	return filtered
}

// ReactNativeEnv builds the environment for the project Expo and React Native CLIs.
func (c Config) ReactNativeEnv(hook string) []string {
	overrides := map[string]string{
		"IDEV_HOME":        c.Home,
		"IDEV_BIN":         c.Self,
		"EXPO_PROJECT_DIR": c.Project,
	}
	for key, value := range c.StateEnv() {
		overrides[key] = value
	}
	nodeOptions := ""
	if hook != "" {
		nodeOptions = "--require " + hook
	}
	if previous := os.Getenv("NODE_OPTIONS"); previous != "" {
		if nodeOptions != "" {
			nodeOptions += " "
		}
		nodeOptions += previous
	}
	if nodeOptions != "" {
		overrides["NODE_OPTIONS"] = nodeOptions
	}

	pathParts := []string{
		c.Path("react-native", "expo", "bin"),
		c.Path("react-native", "bin"),
		c.Path("apple-tools", "tools"),
		c.Path("bin"),
	}
	if sdk := c.SDKPath(); sdk != "" {
		pathParts = append(pathParts, filepath.Join(sdk, "usr", "bin"))
	}
	swiftBin, swiftLib := c.swiftPaths()
	if swiftBin != "" {
		overrides["IDEV_SWIFT_BIN"] = swiftBin
		pathParts = append([]string{swiftBin}, pathParts...)
	}
	pathParts = append(pathParts, os.Getenv("PATH"))
	overrides["PATH"] = strings.Join(pathParts, string(os.PathListSeparator))

	if sdk := c.SDKPath(); sdk != "" {
		overrides["DEVELOPER_DIR"] = sdk
		overrides["DARWIN_SDK_BUNDLE"] = c.SDKRoot()
	}
	setIfMissing(overrides, "ENGINE", c.Path("react-native", "bin", "build"))
	setIfMissing(overrides, "EXPO_BUILD_BIN", c.Path("bin", "expo-build-bin"))
	setIfMissing(overrides, "SWBBUILDSERVICE_PATH", c.Path("bin", "SWBBuildService"))
	setIfMissing(overrides, "SWIFT_BUILD_PATH", c.Path("src", "swift-build"))
	if swiftLib != "" {
		libs := swiftLib + string(os.PathListSeparator) + c.Path("lib", "compat") + string(os.PathListSeparator) + c.Path("lib")
		if existing := os.Getenv("LD_LIBRARY_PATH"); existing != "" {
			libs += string(os.PathListSeparator) + existing
		}
		overrides["LD_LIBRARY_PATH"] = libs
	}
	overrides["IDEV_APPLE_TOOLS"] = c.Path("apple-tools")
	return EnvWith(overrides)
}

func setIfMissing(overrides map[string]string, key, value string) {
	if _, present := os.LookupEnv(key); present {
		return
	}
	if IsFile(value) {
		overrides[key] = value
	}
}
