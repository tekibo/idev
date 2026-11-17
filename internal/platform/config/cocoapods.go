package config

import (
	"os"
	"path/filepath"
)

// GemRoot returns the private CocoaPods gem tree.
func (c Config) GemRoot() string {
	return c.Path("gems")
}

// PodBin returns the gem binstub for pod.
func (c Config) PodBin() string {
	return filepath.Join(c.GemRoot(), "bin", "pod")
}

// PodLink returns the manager pod wrapper.
func (c Config) PodLink() string {
	return c.Path("bin", "pod")
}

// GemCache returns the private gem cache.
func (c Config) GemCache() string {
	return c.Path("cache", "gem")
}

// GemEnv returns the private gem environment.
func (c Config) GemEnv() []string {
	overrides := map[string]string{
		"GEM_HOME":       c.GemRoot(),
		"GEM_PATH":       c.GemRoot(),
		"GEM_SPEC_CACHE": filepath.Join(c.GemCache(), "specs"),
		"PATH":           filepath.Join(c.GemRoot(), "bin") + string(os.PathListSeparator) + os.Getenv("PATH"),
	}
	for key, value := range c.StateEnv() {
		overrides[key] = value
	}
	return EnvWith(overrides)
}
