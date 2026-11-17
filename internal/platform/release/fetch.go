package release

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/download"
)

// Fetch downloads, verifies, and atomically installs one component archive.
func Fetch(ctx context.Context, manifestSource string, component Component, target, cache string) (int, error) {
	return FetchProgress(ctx, manifestSource, component, target, cache, nil)
}

// FetchProgress is Fetch with download progress reporting.
func FetchProgress(ctx context.Context, manifestSource string, component Component, target, cache string, progress download.Progress) (int, error) {
	if component.URL == "" || component.SHA256 == "" {
		return 0, errors.New("release: component archive or checksum is empty")
	}
	archive := target + ".tar.gz"
	if cache != "" {
		archive = filepath.Join(cache, filepath.Base(component.URL))
	}
	url := Resolve(manifestSource, component.URL)
	if err := download.File(ctx, url, component.SHA256, archive, progress); err != nil {
		return 0, fmt.Errorf("release: %w", err)
	}
	if cache == "" {
		defer func() { _ = os.Remove(archive) }()
	}
	file, err := os.Open(archive)
	if err != nil {
		return 0, err
	}
	defer func() { _ = file.Close() }()
	parent := filepath.Dir(target)
	if err := os.MkdirAll(parent, 0755); err != nil {
		return 0, err
	}
	temp, err := os.MkdirTemp(parent, ".idev-release-")
	if err != nil {
		return 0, err
	}
	defer func() { _ = os.RemoveAll(temp) }()
	count, err := extract(file, temp)
	if err != nil {
		return 0, err
	}
	entries, err := os.ReadDir(temp)
	if err != nil {
		return 0, err
	}
	if len(entries) != 1 || !entries[0].IsDir() {
		return 0, fmt.Errorf("release: archive must have one top-level directory")
	}
	if err := os.RemoveAll(target); err != nil {
		return 0, err
	}
	if err := os.Rename(filepath.Join(temp, entries[0].Name()), target); err != nil {
		return 0, err
	}
	return count, nil
}
