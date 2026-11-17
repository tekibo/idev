package config

import (
	"os"
	"os/exec"
	"path/filepath"
)

// IsFile reports a regular file.
func IsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// IsExec reports an executable regular file.
func IsExec(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir() && info.Mode()&0111 != 0
}

// FindBinary looks on PATH, then in the private bin roots.
func (c Config) FindBinary(name string) (string, bool) {
	if path, err := exec.LookPath(name); err == nil {
		return path, true
	}
	candidates := []string{c.Path("bin", name)}
	for _, candidate := range candidates {
		if IsExec(candidate) {
			return candidate, true
		}
	}
	return "", false
}

// SDKRoot returns the private Darwin SDK bundle directory.
func (c Config) SDKRoot() string {
	if c.SDKRootPath != "" && IsFile(filepath.Join(c.SDKRootPath, "info.json")) {
		return c.SDKRootPath
	}
	if value := os.Getenv("DARWIN_SDK_BUNDLE"); value != "" && IsFile(filepath.Join(value, "info.json")) {
		return value
	}
	return ""
}

// SDKPath returns the Darwin SDK Developer directory.
func (c Config) SDKPath() string {
	if value := os.Getenv("DEVELOPER_DIR"); value != "" && IsFile(filepath.Join(value, "Info.plist")) {
		return value
	}
	root := c.SDKRoot()
	if root == "" {
		return ""
	}
	return filepath.Join(root, "Developer")
}
