package config

import "path/filepath"

// DarwinSDKReady reports a registered Darwin SDK bundle.
func (c Config) DarwinSDKReady() bool {
	root := c.SDKRoot()
	return root != "" && IsFile(filepath.Join(root, "info.json"))
}
