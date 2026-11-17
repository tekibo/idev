package uninstall

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// existing keeps the paths that exist and removes duplicates.
func existing(paths []string) []string {
	seen := map[string]bool{}
	var found []string
	for _, path := range paths {
		if path == "" || seen[path] {
			continue
		}
		if _, err := os.Lstat(path); err != nil {
			continue
		}
		seen[path] = true
		found = append(found, path)
	}
	sort.Strings(found)
	return found
}

// appendUnique adds paths without duplicates.
func appendUnique(paths []string, values ...string) []string {
	seen := map[string]bool{}
	for _, path := range paths {
		seen[path] = true
	}
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		paths = append(paths, value)
	}
	sort.Strings(paths)
	return paths
}

// scanDirs returns child directories with an optional prefix.
func scanDirs(parent, prefix string) []string {
	entries, err := os.ReadDir(parent)
	if err != nil {
		return nil
	}
	var paths []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), prefix) {
			paths = append(paths, filepath.Join(parent, entry.Name()))
		}
	}
	return paths
}

// scanSuffix returns child entries with a suffix.
func scanSuffix(parent, suffix string) []string {
	entries, err := os.ReadDir(parent)
	if err != nil {
		return nil
	}
	var paths []string
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), suffix) {
			paths = append(paths, filepath.Join(parent, entry.Name()))
		}
	}
	return paths
}

// isUnder reports whether path is inside root.
func isUnder(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}
