package uninstall

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Apply removes the allowlist paths and writes the receipt.
func Apply(ctx context.Context, report Report) (Receipt, error) {
	return ApplyProgress(ctx, report, nil)
}

// ApplyProgress removes the allowlist paths and reports progress.
func ApplyProgress(ctx context.Context, report Report, progress func(int, int)) (Receipt, error) {
	now := time.Now().UTC()
	receipt := Receipt{
		Time:           now.Format(time.RFC3339),
		Home:           report.Home,
		Kept:           report.Keep,
		KeptComponents: keptComponents(report.Components),
		System:         report.System,
		Configs:        report.Configs,
		ProjectData:    report.ProjectData,
		SystemErr:      report.SystemErr,
	}
	applyErr := removeAll(ctx, &receipt, report, progress)
	if applyErr != nil {
		receipt.Failure = applyErr.Error()
	}
	receipt.Remaining = remaining(report.Home, report.Keep)
	receipt.Path = receiptPath(now)
	writeErr := writeReceipt(receipt)
	if applyErr != nil {
		return receipt, applyErr
	}
	return receipt, writeErr
}

// removeAll removes the allowlist, then empty containers and the empty home.
func removeAll(ctx context.Context, receipt *Receipt, report Report, progress func(int, int)) error {
	for index, path := range report.Remove {
		if err := ctx.Err(); err != nil {
			return err
		}
		if progress != nil {
			progress(index, len(report.Remove))
		}
		if _, err := os.Lstat(path); err != nil {
			continue
		}
		if err := os.RemoveAll(path); err != nil {
			return err
		}
		receipt.Removed = append(receipt.Removed, path)
	}
	if progress != nil {
		progress(len(report.Remove), len(report.Remove))
	}
	for _, dir := range report.Containers {
		_ = os.Remove(dir)
	}
	if entries, err := os.ReadDir(report.Home); err == nil && len(entries) == 0 {
		if err := os.Remove(report.Home); err != nil {
			return err
		}
		receipt.Removed = append(receipt.Removed, report.Home)
	}
	return nil
}

// remaining lists unexpected paths inside the private home.
func remaining(home string, kept []string) []string {
	paths := []string{}
	_ = filepath.WalkDir(home, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || path == home {
			return nil
		}
		if entry.IsDir() {
			if keptPath(path, kept) {
				return filepath.SkipDir
			}
			return nil
		}
		if !keptPath(path, kept) {
			paths = append(paths, path)
		}
		return nil
	})
	sort.Strings(paths)
	return paths
}

// keptPath reports whether path belongs to a kept component.
func keptPath(path string, kept []string) bool {
	for _, root := range kept {
		if isUnder(root, path) {
			return true
		}
	}
	return false
}
