package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

// ValidateHome rejects an IDEV_HOME that points at a system location.
func ValidateHome(home string) error {
	if home == "" {
		return errors.New("home is empty")
	}
	clean := filepath.Clean(home)
	if !filepath.IsAbs(clean) {
		return fmt.Errorf("home %q is not absolute", home)
	}
	if clean == "/" {
		return errors.New("home must not be /")
	}
	for _, root := range []string{
		"/bin", "/boot", "/dev", "/etc", "/lib", "/lib64", "/opt",
		"/proc", "/root", "/run", "/sbin", "/srv", "/sys", "/usr", "/var",
	} {
		if clean == root {
			return fmt.Errorf("home %q is a system root", clean)
		}
	}
	for _, root := range []string{
		"/bin", "/boot", "/dev", "/etc", "/lib", "/lib64",
		"/proc", "/run", "/sbin", "/sys", "/usr", "/var",
	} {
		if strings.HasPrefix(clean, root+"/") {
			return fmt.Errorf("home %q is under the system root %s", clean, root)
		}
	}
	return nil
}
