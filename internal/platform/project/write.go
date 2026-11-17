package project

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var bundleIDPattern = regexp.MustCompile(`^[A-Za-z0-9-]+(\.[A-Za-z0-9-]+)+$`)

// ValidBundleID checks the reverse-DNS form used by Apple.
func ValidBundleID(id string) error {
	if !bundleIDPattern.MatchString(id) {
		return fmt.Errorf("bundle ID %q is not a reverse-DNS identifier", id)
	}
	return nil
}

// WriteBundleID updates the top-level bundleID in xtool.yml.
func WriteBundleID(dir, id string) error {
	if err := ValidBundleID(id); err != nil {
		return err
	}
	path := filepath.Join(dir, "xtool.yml")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read xtool.yml: %w", err)
	}
	lines := strings.Split(string(data), "\n")
	replaced := false
	for index, line := range lines {
		if strings.HasPrefix(line, "bundleID:") {
			lines[index] = "bundleID: " + id
			replaced = true
			break
		}
	}
	if !replaced {
		if len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = append(lines[:len(lines)-1], "bundleID: "+id, "")
		} else {
			lines = append(lines, "bundleID: "+id)
		}
	}
	temp, err := os.CreateTemp(dir, ".xtool-*.yml")
	if err != nil {
		return fmt.Errorf("write xtool.yml: %w", err)
	}
	name := temp.Name()
	defer os.Remove(name)
	if _, err := temp.WriteString(strings.Join(lines, "\n")); err != nil {
		temp.Close()
		return fmt.Errorf("write xtool.yml: %w", err)
	}
	if err := temp.Chmod(0o644); err != nil {
		temp.Close()
		return fmt.Errorf("write xtool.yml: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("write xtool.yml: %w", err)
	}
	if err := os.Rename(name, path); err != nil {
		return fmt.Errorf("write xtool.yml: %w", err)
	}
	return nil
}
