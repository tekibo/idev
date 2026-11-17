// Package config owns manager paths and environment construction.
package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the manager paths for one run.
type Config struct {
	Home        string
	HomeErr     error
	Project     string
	Self        string
	SDKRootPath string
	Sources     Sources
	Lock        Lock
	SourcesErr  error
}

// Load resolves IDEV_HOME, the project directory, and the manager binary.
func Load() Config {
	home := ""
	var homeErr error
	if value := os.Getenv("IDEV_HOME"); value != "" {
		if err := ValidateHome(value); err != nil {
			homeErr = fmt.Errorf("IDEV_HOME: %w", err)
		} else {
			home = filepath.Clean(value)
		}
	}
	if home == "" {
		if userHome, err := os.UserHomeDir(); err == nil {
			home = filepath.Join(userHome, ".local", "share", "idev")
		}
	}
	project, _ := os.Getwd()
	if value := os.Getenv("IDEV_PROJECT"); value != "" {
		project = value
	}
	if abs, err := filepath.Abs(project); err == nil {
		project = abs
	}
	self, _ := os.Executable()
	sources, sourcesErr := LoadSources(home)
	lock, lockErr := LoadLock()
	if sourcesErr == nil {
		sourcesErr = lockErr
	}
	if sourcesErr == nil {
		sourcesErr = homeErr
	}
	return Config{
		Home:       home,
		HomeErr:    homeErr,
		Project:    project,
		Self:       self,
		Sources:    sources,
		Lock:       lock,
		SourcesErr: sourcesErr,
	}
}

// Path joins parts under the private home.
func (c Config) Path(parts ...string) string {
	all := make([]string, 0, len(parts)+1)
	all = append(all, c.Home)
	all = append(all, parts...)
	return filepath.Join(all...)
}

// UserHome returns the real user home.
func UserHome() string {
	home, _ := os.UserHomeDir()
	return home
}
