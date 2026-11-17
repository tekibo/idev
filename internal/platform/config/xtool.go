package config

import (
	"path/filepath"
)

// XtoolRoot returns the private xtool install root.
func (c Config) XtoolRoot() string {
	return c.Path("lib", "xtool", c.Lock.AppleDev.XtoolSHA)
}

// XtoolBin returns the private xtool binary.
func (c Config) XtoolBin() string {
	return filepath.Join(c.XtoolRoot(), "xtool")
}

// XtoolSource returns the pinned xtool source cache.
func (c Config) XtoolSource() string {
	return c.Path("cache", "src", "xtool-"+c.Lock.AppleDev.XtoolSHA)
}

// XtoolLink returns the manager xtool link.
func (c Config) XtoolLink() string {
	return c.Path("bin", "xtool")
}
