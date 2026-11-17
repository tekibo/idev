// Package download fetches remote files with retry, resume, and checksums.
package download

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Progress receives the downloaded byte count and the total when known.
type Progress func(done, total int64)

// File downloads url to destination unless the checksum already matches.
// A partial transfer stays at destination+".part" and resumes with Range.
func File(ctx context.Context, url, checksum, destination string, progress Progress) error {
	checksum = strings.ToLower(strings.TrimSpace(checksum))
	if ok, _ := matches(destination, checksum); ok {
		return nil
	}
	temporary := destination + ".part"
	if ok, _ := matches(temporary, checksum); ok {
		return os.Rename(temporary, destination)
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return err
	}
	var last error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		retry, err := fetch(ctx, url, temporary, progress)
		if err == nil {
			sum, err := sum(temporary)
			if err != nil {
				return err
			}
			if sum != checksum {
				_ = os.Remove(temporary)
				return fmt.Errorf("download: checksum mismatch for %s", filepath.Base(destination))
			}
			return os.Rename(temporary, destination)
		}
		last = err
		if !retry {
			return err
		}
		if attempt < maxAttempts {
			wait(ctx, attempt)
		}
	}
	return fmt.Errorf("download: %s: %w", filepath.Base(destination), last)
}

func matches(path, expected string) (bool, error) {
	value, err := sum(path)
	if err != nil {
		return false, err
	}
	return value == expected, nil
}
