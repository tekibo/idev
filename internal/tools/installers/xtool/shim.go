package xtool

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tekibo/idev/internal/platform/config"
)

// registerShim creates one user-facing xtool shim with the manager environment.
func registerShim(cfg config.Config) (string, bool, error) {
	dir := filepath.Join(config.UserHome(), ".local", "bin")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", false, err
	}
	path := filepath.Join(dir, "xtool")
	if data, err := os.ReadFile(path); err == nil && !strings.Contains(string(data), "idev xtool shim") {
		return "", false, fmt.Errorf("%s exists and is not the idev shim", path)
	}
	swiftBin := cfg.SwiftBin()
	script := fmt.Sprintf(`#!/usr/bin/env bash
# idev xtool shim
export IDEV_HOME=%q
export IDEV_SWIFT_BIN=%q
export XDG_CONFIG_HOME=%q
export XDG_DATA_HOME=%q
export XDG_CACHE_HOME=%q
export SDK_CACHE=%q
export PATH=%q:"$PATH"
export LD_LIBRARY_PATH=%q"${LD_LIBRARY_PATH:+:$LD_LIBRARY_PATH}"
exec %q "$@"
`, cfg.Home, swiftBin, cfg.Path(".config"), cfg.Path(".local", "share"),
		cfg.Path(".cache"), cfg.Path(".cache", "xtool"),
		swiftBin+":"+cfg.Path("bin"),
		cfg.SwiftLib()+":"+cfg.Path("lib", "compat")+":"+cfg.Path("lib"),
		cfg.Path("bin", "xtool"))
	if err := os.WriteFile(path, []byte(script), 0755); err != nil {
		return "", false, err
	}
	onPath := false
	for _, entry := range filepath.SplitList(os.Getenv("PATH")) {
		if filepath.Clean(entry) == dir {
			onPath = true
			break
		}
	}
	return path, onPath, nil
}
