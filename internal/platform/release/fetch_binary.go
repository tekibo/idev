package release

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/download"
)

// FetchBinary downloads, verifies, and installs one binary from a tar.gz.
func FetchBinary(ctx context.Context, manifestSource string, component Component, name, target, cache string) error {
	if component.URL == "" || component.SHA256 == "" {
		return errors.New("release: binary archive or checksum is empty")
	}
	archive := target + ".tar.gz"
	if cache != "" {
		archive = filepath.Join(cache, filepath.Base(component.URL))
	}
	url := Resolve(manifestSource, component.URL)
	if err := download.File(ctx, url, component.SHA256, archive, nil); err != nil {
		return fmt.Errorf("release: %w", err)
	}
	if cache == "" {
		defer func() { _ = os.Remove(archive) }()
	}
	parent := filepath.Dir(target)
	if err := os.MkdirAll(parent, 0755); err != nil {
		return err
	}
	temp, err := os.MkdirTemp(parent, ".idev-binary-")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(temp) }()
	file, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	if _, err := extract(file, temp); err != nil {
		return err
	}
	source := ""
	err = filepath.WalkDir(temp, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		if entry.Name() == name {
			source = path
			return fs.SkipAll
		}
		return nil
	})
	if err != nil {
		return err
	}
	if source == "" {
		return fmt.Errorf("release: %s missing in %s", name, component.URL)
	}
	return copyExecutable(source, target)
}

func copyExecutable(source, target string) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() { _ = input.Close() }()
	temp, err := os.CreateTemp(filepath.Dir(target), ".idev-binary-")
	if err != nil {
		return err
	}
	tempName := temp.Name()
	defer func() { _ = os.Remove(tempName) }()
	if _, err := io.Copy(temp, input); err != nil {
		_ = temp.Close()
		return err
	}
	if err := temp.Chmod(0755); err != nil {
		_ = temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	return os.Rename(tempName, target)
}
