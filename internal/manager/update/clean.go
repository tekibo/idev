package update

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/tekibo/idev/internal/platform/config"
)

// CleanTarget is one removable cache path.
type CleanTarget struct {
	Path  string `json:"path"`
	Bytes int64  `json:"bytes"`
}

// CleanPlan is the cache cleanup plan.
type CleanPlan struct {
	Targets []CleanTarget `json:"targets"`
	Bytes   int64         `json:"bytes"`
}

// CleanResult is the applied cleanup.
type CleanResult struct {
	Targets []CleanTarget `json:"targets"`
	Bytes   int64         `json:"bytes"`
	Removed []string      `json:"removed"`
}

// CleanPlan returns the current removable cache set.
func (s *Service) CleanPlan() CleanPlan {
	plan := CleanPlan{Targets: []CleanTarget{}}
	for _, path := range cleanPaths(s.Config) {
		size := pathSize(path)
		if size == 0 {
			continue
		}
		plan.Targets = append(plan.Targets, CleanTarget{Path: path, Bytes: size})
		plan.Bytes += size
	}
	return plan
}

// CleanApply removes every target in the plan.
func (s *Service) CleanApply(plan CleanPlan) (CleanResult, error) {
	result := CleanResult{Targets: plan.Targets, Bytes: plan.Bytes, Removed: []string{}}
	for _, target := range plan.Targets {
		if err := os.RemoveAll(target.Path); err != nil {
			return result, fmt.Errorf("clean %s: %w", target.Path, err)
		}
		result.Removed = append(result.Removed, target.Path)
	}
	return result, nil
}

func cleanPaths(cfg config.Config) []string {
	return []string{
		cfg.Path("cache"),
		cfg.Path(".cache"),
		cfg.Path(".config", "swiftpm", "cache"),
		cfg.Path(".local", "share", "vite-plus"),
		cfg.Path(".local", "share", "recently-used.xbel"),
		cfg.Path("gems", "cache"),
		cfg.Path("config", "editor.json"),
		filepath.Join(config.UserHome(), ".local", "state", "idev", "logs"),
	}
}

func pathSize(path string) int64 {
	info, err := os.Lstat(path)
	if err != nil {
		return 0
	}
	if !info.IsDir() {
		return info.Size()
	}
	var total int64
	_ = filepath.WalkDir(path, func(_ string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil || !entry.Type().IsRegular() {
			return nil
		}
		if item, err := entry.Info(); err == nil {
			total += item.Size()
		}
		return nil
	})
	return total
}
