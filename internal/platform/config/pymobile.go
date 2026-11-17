package config

import "path/filepath"

// PymobileRoot returns the private pymobiledevice3 venv.
func (c Config) PymobileRoot() string {
	return c.Path("pymobile3-venv")
}

// PymobileBin returns the venv pymobiledevice3 binary.
func (c Config) PymobileBin() string {
	return filepath.Join(c.PymobileRoot(), "bin", "pymobiledevice3")
}

// PymobileLink returns the manager pymobiledevice3 link.
func (c Config) PymobileLink() string {
	return c.Path("bin", "pymobiledevice3")
}

// PymobileCache returns the private pip cache.
func (c Config) PymobileCache() string {
	return c.Path("cache", "pip")
}
